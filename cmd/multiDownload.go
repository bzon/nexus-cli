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
	"io/ioutil"
	"os"
	"strings"

	"github.com/bzon/nexus-cli/nexus2"

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
		var aRequest nexus2.ArtifactRequest
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
				nexus2.DownloadArtifact(aRequest)
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
