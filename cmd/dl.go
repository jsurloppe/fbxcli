package cmd

import (
	"io"
	"mime"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl",
	Short: "Download a file",
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		path := args[0]
		path = makePath(ENV.Cwd, path)
		client, err := getCurrentClient()
		checkErr(err)
		resp, err := client.Dl(path)
		checkErr(err)

		bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
		bar.Start()

		_, params, err := mime.ParseMediaType(resp.Header["Content-Disposition"][0])
		checkErr(err)

		file, err := os.Create(params["filename"])
		checkErr(err)

		reader := bar.NewProxyReader(resp.Body)

		_, err = io.Copy(file, reader)
		checkErr(err)

		bar.Finish()
	},
}

func init() {
	RootCmd.AddCommand(dlCmd)
}
