package nexus3

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// SiteComponent contains the fields that will be passed as a parameter for NexusUpload
type SiteComponent struct {
	File, Filename, Directory string
}

// SiteFileUpload uploads a file to Nexus returns the uploaded file url
func (n *Client) SiteFileUpload(c SiteComponent) (string, error) {
	file, err := os.Open(c.File)
	if err != nil {
		return "", err
	}
	defer file.Close()
	uri := n.GetRepoURL() + "/" + c.Directory + "/" + c.Filename
	req, err := http.NewRequest("PUT", uri, file)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(n.Username, n.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("%s: %s", resp.Status, string(b))
	}
	return uri, nil
}

func (n *Client) GetRepoURL() string {
	return n.HostURL + "/repository/" + n.Repository
}
