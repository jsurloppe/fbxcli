package cmd

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"

	"github.com/chzyer/readline"
	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

const APPID = "com.github.jsurloppe.fbxcli"
const APPNAME = "fbxcli"
const APPVERSION = "0"

var App = &fbxapi.App{
	AppID:      APPID,
	AppVersion: APPVERSION,
}

var RootCmd = &cobra.Command{
	Use:   "fbx",
	Short: "fbx is a freebox command line utility",
	Long:  "A (non+)interactive cli for managing your currentFreebox",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		alias, _ := cmd.Flags().GetString("freebox")
		if alias == "" {
			alias = getDefaultFreebox()
		}
		if alias != "" {
			err := connect(alias)
			checkErr(err)
		}
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		PoolLogout()
		updateConfig()
	},
}

func connect(alias string) (err error) {
	_, ok := ENV.FreeboxsList[alias]
	if !ok {
		return errors.New("Freebox not found, you'll need to enroll it first")
	}

	_, err = NewClient(alias)
	checkErr(err)
	ENV.CurrentAlias = alias

	ENV.Cwd[alias] = "/"
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
