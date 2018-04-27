package nexus3

// Nexus contains the fields requied for accessing a nexus server
type Client struct {
	Repository, HostURL, Username, Password string
}
