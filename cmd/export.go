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
	"os"
	"path"
	"time"

	"github.com/fimreal/docker-exporter/dockercli"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the configuration of a specified Docker container",
	Long: `The export command retrieves and displays the full configuration 
of a specified Docker container. It provides an easy way to see the parameters 
and settings used when the container was created, which can be useful for 
replicating setups or troubleshooting issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll, _ := cmd.Flags().GetBool("all")
		format, _ := cmd.Flags().GetString("format")
		cjson, err := DockerClient.ExportContainersJSON(args, showAll)
		if err != nil {
			ezap.Error(err)
			return
		}

		output, _ := cmd.Flags().GetString("output-dir")
		pretty, _ := cmd.Flags().GetBool("pretty")
		dump := dockercli.ParseContainers(cjson, format, pretty)

		switch t := dump.(type) {
		case string:
			if output != "" {
				os.MkdirAll(output, 0755)
				filename := path.Join(output, "docker_dump-"+time.Now().Format("2006_01_02"))
				if format == "json" {
					filename += ".json"
				} else {
					filename += ".sh"
				}
				ezap.Infof("Writing to %s\n", filename)
				err := os.WriteFile(filename, []byte(dump.(string)), 0644)
				if err != nil {
					ezap.Fatal(err)
				}
			}
			ezap.Println(dump.(string))
		case map[string]string:
			var out string
			for serviceName, content := range t {
				out += "# " + serviceName + "\n"
				out += content + "\n---\n"
				if output != "" {
					os.MkdirAll(output, 0755)
					filename := path.Join(output, serviceName+".yml")
					ezap.Infof("Writing to %s\n", filename)
					err := os.WriteFile(filename, []byte(content), 0644)
					if err != nil {
						ezap.Fatal(err)
					}
				}
			}
			ezap.Println(out)
		default:
			ezap.Fatal("Unknown error")
		}

	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Example flags for the export command (optional, customize as needed)
	exportCmd.Flags().BoolP("pretty", "p", false, "Pretty-print the output")
	exportCmd.Flags().StringP("output-dir", "o", "", "Set output directory for the generated files, if not set, output to stdout")
	exportCmd.Flags().StringP("format", "f", "command", "Set output format (eg. command (shell), compose (yaml))")
}
