package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/chzyer/readline"
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

		if len(alias) > 0 {
			err := connect(alias)
			checkErr(err)
		} else {
			for alias := range ENV.Freeboxs {
				err := connect(alias)
				checkErr(err)
				break
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if !ENV.KeepSession {
			PoolLogout()
		}
		rlshell.Close()
	},
}

func connect(alias string) (err error) {
	cfgEntry, ok := ENV.Freeboxs[alias]
	if !ok {
		return errors.New("Cant load config")
	}

	ENV.CurrentClient, err = NewClientFromPool(alias)
	checkErr(err)
	if len(ENV.CurrentClient.SessionToken) == 0 {
		err = ENV.CurrentClient.OpenSession(APPID, cfgEntry.AppToken)
		checkErr(err)
	}
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
	ENV.Freeboxs = make(map[string]ConfigEntry)

	usr, err := user.Current()
	checkErr(err)

	RootCmd.PersistentFlags().StringVar(&ENV.CfgFile, "config", usr.HomeDir+"/.fbxcli.json", "config file (default is $HOME/.fbxcli.json)")
	RootCmd.PersistentFlags().StringP("freebox", "f", "", "Local configureed name of Freebox to query")
	RootCmd.PersistentFlags().BoolVarP(&ENV.KeepSession, "keep", "k", false, "Keep session open")

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
		err = decoder.Decode(&ENV.Freeboxs)
		checkErr(err)
	}
}
