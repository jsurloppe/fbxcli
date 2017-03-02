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
	"strings"

	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

func printErr(err error) {
	rlshell.writeString(err.Error())
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "An interactive shell",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		panicHandler = recoverOnPanic
		RootCmd.PersistentPreRun(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer panicHandler()

		for {
			line, err := rlshell.Readline()
			if err != nil { // io.EOF
				break
			}

			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			args, err := shellwords.Parse(line)
			checkErr(err)

			replCmd, replCmdArgs, err := RootCmd.Find(args)
			if err != nil {
				printErr(err)
				continue
			}
			replCmd.SetArgs(replCmdArgs)

			if cmd == replCmd {
				falias := replCmd.Flag("freebox")
				alias := falias.Value.String()
				if len(alias) > 0 {
					connect(alias)
				} else if len(replCmdArgs) > 0 {
					connect(replCmdArgs[0])
				}

				continue
			}
			// cause weird nesting in interactive mode
			if replCmd.PreRun != nil {
				replCmd.PreRun(replCmd, replCmdArgs)
			}
			replCmd.Run(replCmd, replCmdArgs)
			if replCmd.PostRun != nil {
				replCmd.PostRun(replCmd, replCmdArgs)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
