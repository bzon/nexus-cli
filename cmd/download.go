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
	"os"

	"github.com/bzon/nexus-cli/nexus"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a single or multiple artifact from Nexus.",
	Long: `Downloads a single or multiple artifact from Nexus.

Specify the GAVP [-g, -a, -v, -p] flags. For example:
nexucli download -H http://localhost:8087 --group com.examplegroup --artifact myartifact --version 1.0.0 --packging jar --destination /tmp/

Or give it a source file with correct formatting. For example:
nexucli download -H http://localhost:8087 -f artifacts.txt -d /tmp/

And the example content of 'artifacts.txt' is written in G:A:V:P format:
----------------------------
com.foo.group:bar:LATEST:jar
com.baz.group:foo:1.0.0:jar
----------------------------

`,
	Run: func(cmd *cobra.Command, args []string) {
		artifact.HostURL = NexusHostURL
		artifact.Username = NexusUsername
		artifact.Password = NexusPassword
		err := nexus.DownloadArtifact(artifact)
		if err != nil {
			fmt.Printf("Download Error: %v", err)
			os.Exit(1)
		}
	},
}

var artifact nexus.ArtifactRequest

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().StringVarP(&artifact.RepositoryID, "repository", "r", "", "The Nexus repository id. Example: 'releases' or 'snapshots'")
	downloadCmd.MarkPersistentFlagRequired("repository")
	downloadCmd.PersistentFlags().StringVarP(&artifact.GroupID, "group", "g", "", "The artifact group id.")
	downloadCmd.MarkPersistentFlagRequired("group")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Artifact, "artifact", "a", "", "The artifact id.")
	downloadCmd.MarkPersistentFlagRequired("artifact")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Packaging, "packaging", "p", "", "The artifact packaging. Example: jar, war, zip, or tar, etc.")
	downloadCmd.MarkPersistentFlagRequired("packaging")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Version, "version", "v", "LATEST", "The artifact version.")
	cwd, _ := os.Getwd()
	downloadCmd.PersistentFlags().StringVarP(&artifact.DestinationDir, "destination", "d", cwd, "The directory where to place the file.")

}
