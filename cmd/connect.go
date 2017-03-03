package cmd

import (
	"strings"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

func printErr(err error) {
	rlshell.writeString(err.Error())
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "An interactive shell",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		panicHandler = recoverOnPanic
		RootCmd.PersistentPreRun(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

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
			replCmd.SetArgs(replCmdArgs)

			if cmd == replCmd {
				falias := replCmd.Flag("freebox")
				alias := falias.Value.String()
				if len(alias) > 0 {
					connect(alias)
				} else if len(replCmdArgs) > 0 {
					connect(replCmdArgs[0])
				}

				continue
			}
			// cause weird nesting in interactive mode
			if replCmd.PreRun != nil {
				replCmd.PreRun(replCmd, replCmdArgs)
			}
			replCmd.Run(replCmd, replCmdArgs)
			if replCmd.PostRun != nil {
				replCmd.PostRun(replCmd, replCmdArgs)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)

}
