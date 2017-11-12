package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jsurloppe/fbxapi"
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
	prompt := fmt.Sprintf("%s:%s > ", ENV.CurrentAlias, ENV.Cwd)
	shell.SetPrompt(prompt)
}

var rlshell tRlShell

var panicHandler = exitOnPanic

var RootCmd = &cobra.Command{
	Use:   "fbx",
	Short: "fbx is a freebox command line utility",
	Long:  "A (non+)interactive cli for managing your currentFreebox",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		alias, _ := cmd.Flags().GetString("freebox")
		if alias != "" {
			err := connect(alias)
			checkErr(err)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		shell(cmd, args)
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		PoolLogout()
		updateConfig()
		rlshell.Close()
	},
}

func connect(alias string) (err error) {
	_, ok := ENV.FreeboxsList[alias]
	if !ok {
		return errors.New("Freebox not found, you'll need to enroll it first")
	}

	_, err = NewClientFromPool(alias)
	checkErr(err)
	ENV.CurrentAlias = alias

	ENV.Cwd = "/"
	rlshell.refreshPrompt()
	return nil
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.Execute()
}

func init() {
	ENV.FreeboxsList = make(map[string]*fbxapi.Freebox)

	usr, err := user.Current()
	checkErr(err)

	RootCmd.PersistentFlags().StringVar(&ENV.CfgFile, "config", usr.HomeDir+"/.fbxcli.json", "config file (default is $HOME/.fbxcli.json)")
	RootCmd.PersistentFlags().StringP("freebox", "f", "", "Local configureed name of Freebox to query")

	cobra.OnInitialize(initConfig)
	rl, err := readline.New("> ")
	checkErr(err)
	rlshell.Instance = rl
}

func initConfig() {
	file, err := os.Open(ENV.CfgFile)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&ENV.FreeboxsList)
		checkErr(err)
	}
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
		replCmd, replCmdArgs, err := cmd.Find(args)
		if err != nil {
			printErr(err)
			continue
		}

		replCmd.ParseFlags(replCmdArgs)
		replCmd.Run(replCmd, replCmdArgs)

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
		/*if replCmd.PreRun != nil {
			replCmd.PreRun(replCmd, replCmdArgs)
		}

		replCmd.Run(replCmd, replCmdArgs)
		if replCmd.PostRun != nil {
			replCmd.PostRun(replCmd, replCmdArgs)
		}*/
	}
}
