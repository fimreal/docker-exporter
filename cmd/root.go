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

	"github.com/fimreal/docker-exporter/dockercli"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	VERSION = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docker-exporter",
	Short: "A CLI tool to export Docker container configurations",
	Long: `docker-exporter is a command-line interface (CLI) tool designed to 
export and display the configuration and parameters of Docker containers.
`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	Version: VERSION,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initCli)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docker-exporter.yaml)")

	rootCmd.PersistentFlags().StringP("docker-host", "H", "unix:///var/run/docker.sock", "Docker daemon api address  (e.g., tcp://localhost:2375 or unix:///var/run/docker.sock)")
	rootCmd.PersistentFlags().StringP("client-version", "V", "1.39", "Docker client version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	ezap.SetLogTime("")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".docker-exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docker-exporter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		ezap.Errorf("Error reading config file[%s]: %v", viper.ConfigFileUsed(), err)
	}

	viper.BindPFlags(rootCmd.Flags())
}

var DockerClient *dockercli.DockerClient

func initCli() {
	var err error
	host := viper.GetString("docker-host")
	clientVersion := viper.GetString("client-version")
	DockerClient, err = dockercli.NewCli(host, clientVersion)
	if err != nil {
		ezap.Fatalf("Error creating Docker client: %v, docker-host: %s, client-version: %s", err, host, clientVersion)
	}
}
