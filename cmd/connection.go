package cmd

import (
	"fmt"

	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "Freebox ifconfig",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		client, err := getCurrentClient()
		checkErr(err)

		state := new(fbxapi.ConnectionStatus)
		err = client.Query(fbxapi.ConnectionEP).Do(&state)
		checkErr(err)

		fmt.Printf("%s:\t<%s,%s>\n", state.Media, state.Type, state.State)
		fmt.Printf("\tinet %s [%d:%d]\n", state.Ipv4, state.Ipv4PortRange[0], state.Ipv4PortRange[1])
		fmt.Printf("\tinet6 %s\n", state.Ipv6)
		fmt.Printf("\tRX current %d (%d Kb/s)\n", state.RateDown, state.RateDown/1000)
		fmt.Printf("\tRX available %d (%d Kb/s)\n", state.BandwidthDown, state.BandwidthDown/1000)
		fmt.Printf("\tRX bytes %d\n", state.BytesDown)
		fmt.Printf("\tTX current %d (%d Kb/s)\n", state.RateUp, state.RateUp/1000)
		fmt.Printf("\tTX available %d (%d Kb/s)\n", state.BandwidthUp, state.BandwidthUp/1000)
		fmt.Printf("\tTX bytes %d\n", state.BytesUp)
	},
}

func init() {
	RootCmd.AddCommand(connectionCmd)
}
