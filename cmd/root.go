// Copyright © 2018 bryansazon@hotmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nexus-cli",
	Short: "Calls Nexus REST API from the commandline.",
	Long:  `Search and Download artifacts from Nexus.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// NexusHostURL is the nexus host url
var NexusHostURL string

// NexusUsername may be required for nexus authentication
var NexusUsername string

// NexusPassword may be required for nexus authentication
var NexusPassword string

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&NexusHostURL, "hostURL", "H", "", "The Nexus host url including the protocol. Defaults to Env $NEXUS_HOST.")
	RootCmd.PersistentFlags().StringVarP(&NexusUsername, "username", "U", "", "The Nexus host url including the protocol. Defaults to Env $NEXUS_USERNAME.")
	RootCmd.PersistentFlags().StringVarP(&NexusPassword, "password", "P", "", "The Nexus host url including the protocol. Defaults to Env $NEXUS_PASSWORD")

	// Lookup for the Environment variables
	if h := os.Getenv("NEXUS_HOST"); h == "" {
		RootCmd.MarkPersistentFlagRequired("hostURL")
	} else {
		NexusHostURL = h
	}
	if u := os.Getenv("NEXUS_USERNAME"); u == "" {
		RootCmd.MarkPersistentFlagRequired("password")
	} else {
		NexusUsername = u
	}
	if p := os.Getenv("NEXUS_PASSWORD"); p == "" {
		RootCmd.MarkPersistentFlagRequired("username")
	} else {
		NexusPassword = p
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nexuscli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".nexuscli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nexuscli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
