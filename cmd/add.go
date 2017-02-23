// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/jsurloppe/fbxapi"
	"github.com/spf13/cobra"
)

const defaultPort = 80
const defaultPortSSL = 443

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add local or remote freebox with address and port",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		alias, err := flags.GetString("alias")
		checkErr(err)
		if len(alias) == 0 {
			log.Panic("alias flag required")
		}
		host, err := flags.GetString("host")
		checkErr(err)
		if len(host) == 0 {
			log.Panic("host flag required")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		port, err := flags.GetInt("port")
		checkErr(err)
		ssl, err := flags.GetBool("ssl")
		checkErr(err)
		if ssl && port == defaultPort {
			port = defaultPortSSL
		}
		alias, err := flags.GetString("alias")
		checkErr(err)
		host, err := flags.GetString("host")
		checkErr(err)
		freebox, err := fbxapi.HttpDiscover(host, port, ssl)
		checkErr(err)
		Register(alias, freebox, ssl)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	addCmd.Flags().String("alias", "", "Alias of the freebox")
	addCmd.Flags().String("host", "", "Host of the freebox")
	addCmd.Flags().Int("port", defaultPort, "Port of the freebox")
	addCmd.Flags().Bool("ssl", false, "Require https")

}
