package cmd

import (
	"github.com/spf13/cobra"
)

var showHidden bool

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List files",
	Run: func(cmd *cobra.Command, args []string) {
		path := ""
		if len(args) > 0 {
			path = args[0]
		}

		alias := ENV.CurrentAlias
		cwd := ENV.Cwd[alias]

		path = makePath(cwd, path)
		client, err := getCurrentClient()
		checkErr(err)
		resp, err := client.Ls(path, false, false, !showHidden)
		checkErr(err)

		for _, f := range resp {
			rlshell.writeString(f.Name)
		}
	},
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolVarP(&showHidden, "hidden", "a", false, "show hidden files")
}
