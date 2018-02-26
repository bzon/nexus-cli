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
	"io"
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

func checkErr(err error) {
	color.Set(color.FgRed)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
	color.Unset()
}

func setRepository(n *ArtifactRequest) {
	if matched, _ := regexp.MatchString(".+-SNAPSHOT", n.Version); matched {
		n.RepositoryID = "snapshots"
	} else {
		n.RepositoryID = "releases"
	}
}

// NewNexusQuery adds the required request Body parameters to the Query and then executes it
func NewNexusQuery(req *http.Request, n ArtifactRequest) *http.Response {
	req.SetBasicAuth(n.Username, n.Password)
	q := req.URL.Query()
	if n.RepositoryID == "" {
		setRepository(&n)
	}
	q.Add("r", n.RepositoryID)
	q.Add("g", n.GroupID)
	q.Add("v", n.Version)
	q.Add("a", n.Artifact)
	q.Add("p", n.Packaging)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	checkErr(err)
	if resp.StatusCode != http.StatusOK {
		checkErr(fmt.Errorf("Got %s while querying %s", resp.Status, req.URL.String()))
	} else {
		color.Set(color.FgGreen)
		fmt.Println("/"+resp.Request.Method, resp.Status, resp.Request.URL)
		color.Unset()
	}
	return resp
}

// DownloadArtifact downloads artifacts from Nexus and validates it
func DownloadArtifact(aRequest ArtifactRequest) error {
	// Resolve and validate the artifact to download
	aResolution, err := GetArtifactResolution(aRequest)
	checkErr(err)
	data := aResolution.Data

	filePath := aRequest.DestinationDir + "/" + data.ArtifactID + "-" + data.Version + "." + data.Extension

	// Download the resolved artifact
	fmt.Printf("Downloading file %s:%s:%s:%s\n", data.GroupID, data.ArtifactID, data.Version, data.Extension)
	req, err := http.NewRequest("GET", aRequest.HostURL+MavenRedirectPath, nil)
	req.Header.Add("Accept", "application/xml")
	checkErr(err)
	resp := NewNexusQuery(req, aRequest)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	checkErr(err)
	io.Copy(out, resp.Body)

	// Compare remote and local file sha1
	remoteSHA1 := data.Sha1
	fmt.Println("Got remote sha1:", remoteSHA1)
	f, err := os.Open(filePath)
	checkErr(err)
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	checkErr(err)
	hashInBytes := h.Sum(nil)
	localSHA1 := hex.EncodeToString(hashInBytes)
	fmt.Printf("Got downloaded sha1: %s\n", localSHA1)
	if remoteSHA1 != localSHA1 {
		return fmt.Errorf("Download error. There is a mismatch in sha1sum")
	}
	color.Set(color.FgGreen)
	fmt.Printf("Successfully downloaded the file %s\n", filePath)
	color.Unset()
	return nil
}

// GetArtifactResolution resolves the ArtifactRequest and return the data needed for ArtifactResolution
func GetArtifactResolution(n ArtifactRequest) (*ArtifactResolution, error) {
	fmt.Println("Resolving the artifact to download.")
	req, err := http.NewRequest("GET", n.HostURL+MavenResolvePath, nil)
	req.Header.Add("Accept", "application/json")
	checkErr(err)
	resp := NewNexusQuery(req, n)
	defer resp.Body.Close()
	ar := new(ArtifactResolution)
	decodedJSON := json.NewDecoder(resp.Body)
	err = decodedJSON.Decode(ar)
	checkErr(err)
	return ar, nil
}

func (aResolve *ArtifactResolution) String() string {
	b, _ := json.Marshal(aResolve)
	return string(b)
}
