package cmd

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of fbxcli registration for this host",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()
		client, err := getCurrentClient()
		checkErr(err)

		resp, err := client.TrackLogin(client.Freebox.TrackID)
		checkErr(err)
		rlshell.writeString(resp.Status)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
