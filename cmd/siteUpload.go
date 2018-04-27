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

	"github.com/bzon/nexus-cli/nexus3"
	"github.com/spf13/cobra"
)

var site nexus3.Client

// siteUploadCmd represents the siteUpload command
var siteUploadCmd = &cobra.Command{
	Use:   "site-upload",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(site.HostURL)
	},
}

func init() {
	RootCmd.AddCommand(siteUploadCmd)
	siteUploadCmd.PersistentFlags().StringVarP(&site.Repository, "repo", "r", "", "nexus 3 site raw repository")
	siteUploadCmd.MarkPersistentFlagRequired("repo")
	fmt.Println(NexusHostURL)
	site.HostURL = NexusHostURL
	site.Username = NexusUsername
	site.Password = NexusPassword
}
