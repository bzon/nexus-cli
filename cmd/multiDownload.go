// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bzon/nexus-cli/nexus"

	"github.com/spf13/cobra"
)

// multiDownloadCmd represents the multiDownload command
var multiDownloadCmd = &cobra.Command{
	Use:   "multi-download",
	Short: "Downloads a single artifact from Nexus using a configuration file.",
	Long: `Downloads a single artifact from Nexus using a configuration file.

Or give it a configuration '.txt' file with correct formatting. For example:
nexus-cli multi-download -H http://localhost:8087 -f artifacts.txt -d /tmp/

And the example content of 'artifacts.txt' is written in G:A:V:P format:
----------------------------
com.foo.group:bar:LATEST:jar
com.baz.group:foo:1.0.0:jar
----------------------------`,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
		artifacts := strings.Split(string(b), "\n")
		var aRequest nexus.ArtifactRequest
		aRequest.HostURL = NexusHostURL
		aRequest.Username = NexusUsername
		aRequest.Password = NexusPassword
		aRequest.DestinationDir = destinationDir
		// for txt files
		if strings.HasSuffix(configFile, ".txt") {
			for i, a := range artifacts {
				fmt.Printf("============== [%d] - Found %s in %s ==============\n", i, a, configFile)
				aRequest.GroupID = strings.Split(a, ":")[0]
				aRequest.Artifact = strings.Split(a, ":")[1]
				aRequest.Version = strings.Split(a, ":")[2]
				aRequest.Packaging = strings.Split(a, ":")[3]
				nexus.DownloadArtifact(aRequest)
			}
		}
	},
}

var configFile, destinationDir string

func init() {
	RootCmd.AddCommand(multiDownloadCmd)
	multiDownloadCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "", "The artifacts file.")
	multiDownloadCmd.MarkPersistentFlagRequired("file")
	cwd, _ := os.Getwd()
	multiDownloadCmd.PersistentFlags().StringVarP(&destinationDir, "destination", "d", cwd, "The directory where to place the file.")
}
