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

	"github.com/bzon/nexus-cli/nexus"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a single artifact from Nexus.",
	Long: `Downloads a single artifact from Nexus.

Specify the GAVP [-g, -a, -v, -p] flags. For example:
nexus-cli download -H http://localhost:8087 --group com.examplegroup --artifact myartifact --version 1.0.0 --packging jar --destination /tmp/`,
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
	downloadCmd.PersistentFlags().StringVarP(&artifact.GroupID, "group", "g", "", "The artifact group id.")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Artifact, "artifact", "a", "", "The artifact id.")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Packaging, "packaging", "p", "", "The artifact packaging. Example: jar, war, zip, or tar, etc.")
	downloadCmd.PersistentFlags().StringVarP(&artifact.Version, "version", "v", "LATEST", "The artifact version.")
	cwd, _ := os.Getwd()
	downloadCmd.PersistentFlags().StringVarP(&artifact.DestinationDir, "destination", "d", cwd, "The directory where to place the file.")
	downloadCmd.MarkPersistentFlagRequired("group")
	downloadCmd.MarkPersistentFlagRequired("artifact")
	downloadCmd.MarkPersistentFlagRequired("packaging")
}
