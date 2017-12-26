package cmd

import (
	"fmt"
	"time"

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

		if showLogs, _ := cmd.Flags().GetBool("logs"); showLogs {
			var logs []fbxapi.ConnectionLog
			err = client.Query(fbxapi.ConnectionLogEP).Do(&logs)
			for _, log := range logs {
				// FIXME: wrong timezone after DST
				tm := time.Unix(int64(log.Date), 0)
				name := log.Conn
				if name == "" {
					name = log.Link
				}
				fmt.Printf("%s\t%s\t%-10s\t%s\t%10d\t%10d\n", tm, log.Type, name, log.State, log.BandwidthDown, log.BandwidthUp)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(connectionCmd)
	connectionCmd.Flags().Bool("logs", false, "show connection logs")
}
