package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// cdCmd represents the ls command
var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Change directory",
	Run: func(cmd *cobra.Command, args []string) {
		path := "/"
		if len(args) > 0 {
			path = strings.TrimSpace(args[0])
			if len(path) == 0 {
				path = "/"
			}
		}

		if path == "." {
			return
		}

		if path == ".." {
			i := strings.LastIndex(ENV.Cwd, "/")
			if i > 0 {
				path = ENV.Cwd[:i]
			} else {
				path = "/"
			}
		}

		if !strings.HasPrefix(path, "/") {
			if ENV.Cwd != "/" {
				path = fmt.Sprintf("%s/%s", ENV.Cwd, path)
			} else {
				path = fmt.Sprintf("/%s", path)
			}
		}
		resp, err := ENV.CurrentClient.Info(path)
		checkErr(err)
		if resp.Type == "dir" {
			ENV.Cwd = path
			rlshell.refreshPrompt()
		}
	},
}

func init() {
	RootCmd.AddCommand(cdCmd)
}
