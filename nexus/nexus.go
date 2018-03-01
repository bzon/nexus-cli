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

package nexus

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/fatih/color"
)

const (
	// MavenRedirectPath is used to download an artifact
	MavenRedirectPath = "/service/local/artifact/maven/redirect"
	// MavenResolvePath is used to get artifact's metadata
	MavenResolvePath = "/service/local/artifact/maven/resolve"
)

// ArtifactResolution contains all the data when querying
// http://localhost:8081/nexus/nexus-restlet1x-plugin/default/docs/path__artifact_maven_resolve.htm
// When media type is `application/json`
type ArtifactResolution struct {
	Data struct {
		PresentLocally      bool   `json:"presentLocally"`
		GroupID             string `json:"groupId"`
		ArtifactID          string `json:"artifactId"`
		Version             string `json:"version"`
		Extension           string `json:"extension"`
		Snapshot            bool   `json:"snapshot"`
		SnapshotBuildNumber int    `json:"snapshotBuildNumber"`
		SnapshotTimeStamp   int    `json:"snapshotTimeStamp"`
		Sha1                string `json:"sha1"`
		RepositoryPath      string `json:"repositoryPath"`
	} `json:"data"`
}

// ArtifactRequest holds the required fields for performing NewNexusQuery
type ArtifactRequest struct {
	Username, Password, HostURL, RepositoryID, GroupID, Version, Artifact, Packaging, DestinationDir string
}

func handleError(err error) {
	color.Set(color.FgRed)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	color.Unset()
}

func setRepository(aRequest *ArtifactRequest) {
	if matched, _ := regexp.MatchString(".+-SNAPSHOT", aRequest.Version); matched {
		aRequest.RepositoryID = "snapshots"
	} else {
		aRequest.RepositoryID = "releases"
	}
}

// NewNexusQuery adds the required request Body parameters to the Query and then executes it
func NewNexusQuery(req *http.Request, n ArtifactRequest) *http.Response {
	req.SetBasicAuth(n.Username, n.Password)
	query := req.URL.Query()
	if n.RepositoryID == "" {
		setRepository(&n)
	}
	query.Add("r", n.RepositoryID)
	query.Add("g", n.GroupID)
	query.Add("v", n.Version)
	query.Add("a", n.Artifact)
	query.Add("p", n.Packaging)
	req.URL.RawQuery = query.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	handleError(err)
	if resp.StatusCode != http.StatusOK {
		handleError(fmt.Errorf("Got %s while querying %s", resp.Status, req.URL.String()))
	} else {
		color.Set(color.FgGreen)
		fmt.Println("/"+resp.Request.Method, resp.Status, resp.Request.URL)
		color.Unset()
	}
	return resp
}

// DownloadArtifact downloads artifacts from Nexus and validates it
func DownloadArtifact(aRequest ArtifactRequest) (string, error) {
	// Resolve and validate the artifact to download
	aResolution, err := GetArtifactResolution(aRequest)
	handleError(err)
	data := aResolution.Data

	// Declare the file path where to place the downloaded bytes
	filePath := aRequest.DestinationDir + "/" + data.ArtifactID + "-" + data.Version + "." + data.Extension

	// Download the resolved artifact
	fmt.Printf("Downloading file %s:%s:%s:%s\n", data.GroupID, data.ArtifactID, data.Version, data.Extension)
	req, err := http.NewRequest("GET", aRequest.HostURL+MavenRedirectPath, nil)
	req.Header.Add("Accept", "application/xml")
	handleError(err)
	resp := NewNexusQuery(req, aRequest)
	defer resp.Body.Close()

	// Create the file
	downloadedBytes, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	err = ioutil.WriteFile(filePath, downloadedBytes, 0644)
	handleError(err)

	// Get Remote file metadata SHA1
	remoteSHA1 := data.Sha1
	fmt.Println("Got remote sha1:", remoteSHA1)

	// Get Local downloaded file SHA1
	b, err := ioutil.ReadFile(filePath)
	hash := sha1.New()
	_, err = hash.Write(b)
	hashInBytes := hash.Sum(nil)
	localSHA1 := hex.EncodeToString(hashInBytes)
	fmt.Printf("Got downloaded sha1: %s\n", localSHA1)

	// Compare SHA1s and return and error if it didn't match
	if remoteSHA1 != localSHA1 {
		return "", fmt.Errorf("Download error. There is a mismatch in sha1sum")
	}

	// Print a successful message!
	color.Set(color.FgGreen)
	fmt.Printf("Successfully downloaded the file %s\n", filePath)
	color.Unset()
	return filePath, nil
}

// GetArtifactResolution resolves the ArtifactRequest and return the data needed for ArtifactResolution
func GetArtifactResolution(aRequest ArtifactRequest) (*ArtifactResolution, error) {
	fmt.Println("Resolving the artifact to download.")
	req, err := http.NewRequest("GET", aRequest.HostURL+MavenResolvePath, nil)
	req.Header.Add("Accept", "application/json")
	handleError(err)
	resp := NewNexusQuery(req, aRequest)
	defer resp.Body.Close()
	aResolution := new(ArtifactResolution)
	decodedJSON := json.NewDecoder(resp.Body)
	err = decodedJSON.Decode(aResolution)
	handleError(err)
	return aResolution, nil
}
