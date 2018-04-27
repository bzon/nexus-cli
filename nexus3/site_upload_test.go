package nexus3

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var nexus = Client{
	HostURL:    "http://localhost:8081",
	Username:   "admin",
	Password:   "admin123",
	Repository: "site",
}

func TestSiteFileUpload(t *testing.T) {
	ioutil.WriteFile("file.txt", []byte("foo"), 0644)
	defer os.Remove("file.txt")
	var testComponent = SiteComponent{
		File:      "file.txt",
		Filename:  "file.txt",
		Directory: "go_upload_test",
	}
	uri, err := nexus.SiteFileUpload(testComponent)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("nexus url:", uri)
}
