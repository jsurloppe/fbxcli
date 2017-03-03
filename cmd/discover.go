package cmd

import (
	"fmt"
	"strings"

	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

func isKnownFreebox(host string) bool {
	for _, knHost := range getKnownHosts() {
		if host == knHost {
			return true
		}
	}
	return false
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover freebox(s) on the lan",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		boxs := fbxapi.MdnsDiscover()
		nboxs := len(boxs)

		for _, freebox := range boxs {
			if isKnownFreebox(freebox.Host) {
				nboxs--
				continue
			}
			line := fmt.Sprintf("Found %s at %s:%d\nEnter a name for saving it or leave blank for ignore:\n",
				freebox.DeviceName, freebox.Host, freebox.Port)
			_, err := rlshell.Write([]byte(line))
			checkErr(err)
			alias, err := rlshell.Readline()
			checkErr(err)
			alias = strings.TrimSpace(alias)
			if len(alias) > 0 {
				client, track_id, err := Register(alias, freebox, false)
				checkErr(err)
				line = fmt.Sprintf("Added %s as %s\n(touch the right arrow on the freebox, then press enter)\n", freebox.DeviceName, alias)
				_, err = rlshell.Write([]byte(line))
				rlshell.Readline()
				resp, err := client.TrackLogin(track_id)
				rlshell.writeString(resp.Status)
				checkErr(err)
			}
		}

		if nboxs == 0 {
			rlshell.Write([]byte("No new freebox found\n"))
		}
	},
}

func init() {
	// logrus.SetOutput(os.Stdout)
	RootCmd.AddCommand(discoverCmd)
}
