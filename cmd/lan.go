package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lanCmd = &cobra.Command{
	Use:   "lan",
	Short: "List currently connected devices",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		iface, err := cmd.Flags().GetString("iface")
		checkErr(err)

		client, err := getCurrentClient()
		checkErr(err)
		devices, err := client.Interface(iface)
		checkErr(err)

		for _, device := range devices {
			if device.Active {
				line := fmt.Sprintf("[%s] %s %s %s\n", device.HostType, device.PrimaryName, device.L2Ident.Type, device.GetIPv4s())
				_, err := rlshell.Write([]byte(line))
				checkErr(err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(lanCmd)
	lanCmd.Flags().String("iface", "pub", "The freebox interface to scan (default: pub)")
}
