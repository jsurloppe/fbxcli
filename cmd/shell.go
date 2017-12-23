package cmd

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

type tRlShell struct {
	*readline.Instance
}

func (shell *tRlShell) writeString(str string) {
	if !strings.HasSuffix(str, "\n") {
		str = str + "\n"
	}
	shell.Write([]byte(str))
}

func printErr(err error) {
	rlshell.writeString(err.Error())
}

func (shell *tRlShell) refreshPrompt() {
	prompt := fmt.Sprintf("%s:%s > ", ENV.CurrentAlias, ENV.Cwd[ENV.CurrentAlias])
	shell.SetPrompt(prompt)
}

var rlshell tRlShell

var panicHandler = exitOnPanic

var ShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "fbx is a freebox command line utility",
	Long:  "A (non+)interactive cli for managing your currentFreebox",
	Run: func(cmd *cobra.Command, args []string) {
		shell(cmd, args)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		rlshell.Close()
	},
}

func init() {
	RootCmd.AddCommand(ShellCmd)
}

func shell(cmd *cobra.Command, args []string) {
	panicHandler = recoverOnPanic
	defer panicHandler()

	getCurrentClient()
	rlshell.refreshPrompt()

	for {
		line, err := rlshell.Readline()
		if err != nil { // io.EOF
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		args, err := shellwords.Parse(line)
		checkErr(err)
		replCmd, replCmdArgs, err := RootCmd.Find(args)
		if err != nil {
			printErr(err)
			continue
		}

		replCmd.ParseFlags(replCmdArgs)
		// replCmd.ExecuteC()

		// replCmd.Run(replCmd, replCmdArgs)

		/*if cmd == replCmd {
			falias := replCmd.Flag("freebox")
			alias := falias.Value.String()
			if len(alias) > 0 {
				connect(alias)
			} else if len(replCmdArgs) > 0 {
				connect(replCmdArgs[0])
			}

			continue
		}*/
		// cause weird nesting in interactive mode
		if replCmd.PreRun != nil {
			replCmd.PreRun(replCmd, replCmdArgs)
		}

		replCmd.Run(replCmd, replCmdArgs)

		if replCmd.PostRun != nil {
			replCmd.PostRun(replCmd, replCmdArgs)
		}
	}
}
