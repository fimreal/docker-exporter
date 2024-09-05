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
	"github.com/fimreal/docker-exporter/dockercli"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <CONTAINER...>",
	Short: "Inspect a Docker container",
	Long:  `The inspect command displays the complete configuration of a specified Docker container in JSON format.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 查询容器
		containers, err := DockerClient.Find(args)
		if err != nil {
			ezap.Error(err)
			return
		}
		containersJSON, err := DockerClient.Inspect(containers)
		if err != nil {
			ezap.Error(err)
			return
		}
		ezap.Println(dockercli.Containers2JSON(containersJSON))
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
