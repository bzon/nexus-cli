package nexus

import (
	"fmt"
	"os"
	"testing"
)

var aRequest = ArtifactRequest{"admin", "admin123", "http://localhost:8081/nexus", "releases", "com.example", "LATEST", "artifactA", "jar", "."}

func TestDownload(t *testing.T) {
	f, err := DownloadArtifact(aRequest)
	if err != nil {
		t.Errorf("DownloadArtifact(aRequest) %v", err)
	}
	err = os.Remove(f)
	if err != nil {
		fmt.Printf("Failed deleting file %s\n", f)
		os.Exit(1)
	}
	fmt.Printf("File %s deleted\n", f)
}
