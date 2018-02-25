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
