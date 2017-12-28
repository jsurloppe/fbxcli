package cmd

import (
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Upload a file",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		path := args[0]
		dest := args[1]

		client, err := getCurrentClient()
		checkErr(err)
		err = client.Upload(path, dest)
		checkErr(err)
	},
}

func init() {
	RootCmd.AddCommand(putCmd)
}
