package cmd

import (
	"strconv"

	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of fbxcli registration for this host",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()
		client, err := getCurrentClient()
		checkErr(err)

		params := map[string]string{
			"track_id": strconv.Itoa(client.Freebox.TrackID),
		}

		state := new(fbxapi.AuthorizationState)
		err = client.Query(fbxapi.TrackAuthorizeEP).As(params).Do(&state)
		checkErr(err)

		rlshell.writeString(state.Status)
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
