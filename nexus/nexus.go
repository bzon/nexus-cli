package nexus

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

// NewNexusQuery adds the required request Body parameters to the Query and then executes it
func NewNexusQuery(req *http.Request, n ArtifactRequest) *http.Response {
	req.SetBasicAuth(n.Username, n.Password)
	q := req.URL.Query()
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
func DownloadArtifact(n ArtifactRequest) error {
	filePath := n.DestinationDir + "/" + n.Artifact + "-" + n.Version + "." + n.Packaging
	fmt.Printf("Downloading file %s:%s:%s:%s\n", n.GroupID, n.Artifact, n.Version, n.Packaging)
	req, err := http.NewRequest("GET", n.HostURL+MavenRedirectPath, nil)
	req.Header.Add("Accept", "application/xml")
	checkErr(err)
	resp := NewNexusQuery(req, n)
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	checkErr(err)
	io.Copy(out, resp.Body)

	// Compare remote and local file sha1
	remoteSHA1, _ := GetRemoteSHA1(n)
	fmt.Println("Got remote sha1:", remoteSHA1)
	f, err := os.Open(filePath)
	checkErr(err)
	defer f.Close()
	h := sha1.New()
	b, err := io.Copy(h, f)
	checkErr(err)
	fmt.Printf("%d bytes downloaded\n", b)
	hashInBytes := h.Sum(nil)
	localSHA1 := hex.EncodeToString(hashInBytes)
	fmt.Printf("Got downloaded sha1: %s\n", localSHA1)
	if remoteSHA1 != localSHA1 {
		return fmt.Errorf("Download error. There is a mismatch in sha1sum")
	}
	color.Set(color.FgGreen)
	fmt.Println("Successfully validated the file! Download complete!")
	color.Unset()
	return nil
}

// GetRemoteSHA1 gets the remote sha1 of the file from Nexus
func GetRemoteSHA1(n ArtifactRequest) (string, error) {
	req, err := http.NewRequest("GET", n.HostURL+MavenResolvePath, nil)
	req.Header.Add("Accept", "application/json")
	checkErr(err)
	resp := NewNexusQuery(req, n)
	defer resp.Body.Close()
	ar := new(ArtifactResolution)
	decodedJSON := json.NewDecoder(resp.Body)
	err = decodedJSON.Decode(ar)
	checkErr(err)
	return ar.Data.Sha1, nil
}
