package cmd

import (
	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of fbxcli registration for this host",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()
		alias, _ := cmd.Flags().GetString("freebox")

		freebox := ENV.Freeboxs[alias]
		client, err := fbxapi.NewClientFromFreebox(freebox.Freebox, freebox.UseSSL)
		checkErr(err)
		resp, err := client.TrackLogin(freebox.TrackID)
		checkErr(err)
		rlshell.writeString(resp.Status)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
