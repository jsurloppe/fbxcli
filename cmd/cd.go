package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

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

		alias := ENV.CurrentAlias

		cwd := getCwd(alias)

		if path == "." {
			return
		}

		if path == ".." {
			i := strings.LastIndex(cwd, "/")
			if i > 0 {
				path = cwd[:i]
			} else {
				path = "/"
			}
		}

		if !strings.HasPrefix(path, "/") {
			if cwd != "/" {
				path = fmt.Sprintf("%s/%s", cwd, path)
			} else {
				path = fmt.Sprintf("/%s", path)
			}
		}
		client, err := getCurrentClient()
		checkErr(err)
		resp, err := client.Info(path)
		checkErr(err)
		if resp.Type == "dir" {
			ENV.Cwd[alias] = path
			rlshell.refreshPrompt()
		}
	},
}

func init() {
	RootCmd.AddCommand(cdCmd)
}
