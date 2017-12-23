package cmd

import (
	"fmt"

	"github.com/jsurloppe/fbxapi"
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

		params := map[string]string{
			"iface": iface,
		}

		var hosts []fbxapi.LanHost
		err = client.Query(fbxapi.InterfaceEP).As(params).Do(&hosts)
		checkErr(err)

		for _, host := range hosts {
			if host.Active {
				line := fmt.Sprintf("[%s] %s %s %s\n", host.HostType, host.PrimaryName, host.L2Ident.ID, host.GetIPv4s())
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
