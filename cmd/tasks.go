package cmd

import "github.com/spf13/cobra"

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		ENV.CurrentClient.Tasks()
	},
}

func init() {
	RootCmd.AddCommand(tasksCmd)
}
