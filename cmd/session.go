package cmd

import "github.com/spf13/cobra"

var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Give me a session key",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()
		ENV.KeepSession = true

		rlshell.writeString(ENV.CurrentClient.SessionToken)
	},
}

func init() {
	RootCmd.AddCommand(sessionCmd)
}
