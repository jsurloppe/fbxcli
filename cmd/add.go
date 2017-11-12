package cmd

import (
	"log"

	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

const defaultPort = 80
const defaultPortSSL = 443

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add local or remote freebox with address and port",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		alias, err := flags.GetString("alias")
		checkErr(err)
		if len(alias) == 0 {
			log.Panic("alias flag required")
		}
		host, err := flags.GetString("host")
		checkErr(err)
		if len(host) == 0 {
			log.Panic("host flag required")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		port, err := flags.GetInt("port")
		checkErr(err)
		alias, err := flags.GetString("alias")
		checkErr(err)
		host, err := flags.GetString("host")
		checkErr(err)
		freebox, err := fbxapi.HttpDiscover(host, port)
		checkErr(err)
		Register(alias, freebox)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().String("alias", "", "Alias of the freebox")
	addCmd.Flags().String("host", "", "Host of the freebox")
	addCmd.Flags().Int("port", defaultPort, "Port of the freebox")
	addCmd.Flags().Bool("ssl", false, "Require https")

}
