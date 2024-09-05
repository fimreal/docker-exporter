/*
Copyright © 2024 fimreal

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/docker/docker/api/types"
	"github.com/fimreal/docker-exporter/dockercli"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Aliases: []string{"l", "ps"}, // 可选：为命令设置别名
	Use:     "list [CONTAINER...]",
	Short:   "List Docker containers",
	Long: `The list command displays all currently running Docker containers along with their configurations. 
Given no arguments, the list command will display all running containers.
Given one or more container names or IDs, the list command will display only those containers.
Use the -a flag to include stopped containers in the output.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll, _ := cmd.Flags().GetBool("all")
		pretty, _ := cmd.Flags().GetBool("pretty")

		var containers []types.Container
		var err error

		if len(args) > 0 {
			// list specific containers
			containers, err = DockerClient.Find(args)
		} else {
			// list all containers
			containers, err = DockerClient.List(showAll)
		}
		if err != nil {
			ezap.Println(err)
			return
		}
		// format and print the list
		dockercli.ListPrint(containers, pretty)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "Include stopped containers in the output")
	listCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the output")
}
