package webdav

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/babolivier/scanner/config"
)

// Client is a WebDAV client that can upload file contents to a WebDAV server.
type Client struct {
	client *http.Client
	cfg    *config.WebDAVConfig
}

// NewClient returns a new Client.
func NewClient(cfg *config.WebDAVConfig) *Client {
	return &Client{
		client: new(http.Client),
		cfg:    cfg,
	}
}

// Upload creates a file with the given content on the WebDAV server. It names this file
// using the current time and the given file type, and returns the generated name.
func (c *Client) Upload(body *bytes.Buffer, fileType string) (string, error) {
	// Determine the file's name.
	fileName := fmt.Sprintf(
		"%s.%s",
		time.Now().Format("2006-01-02_15-04-05"),
		fileType,
	)

	logrus.
		WithField("filename", fileName).
		Info("Uploading file to the WebDAV server")

	// Parse the root URL. Ideally we'd do this in NewClient, but we need to change the
	// path of this URL with the file's name, and we don't want this change to persist
	// on the client.
	u, err := url.Parse(c.cfg.RootURL)
	if err != nil {
		return "", err
	}

	// Build a path that includes the full path for this file on the WebDAV server.
	u.Path = path.Join(u.Path, c.cfg.UploadPath, fileName)

	// Create the request.
	req, err := http.NewRequest("PUT", u.String(), body)
	if err != nil {
		return "", err
	}

	// Add basic auth to the request.
	req.SetBasicAuth(c.cfg.User, c.cfg.Password)

	// Send the request.
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	logrus.WithField("status_code", resp.StatusCode).Info("Upload finished")

	// According to RFC4918, the creation of a resource must be indicated by use of a
	// 201 Created response code, so return an error if that's not what we got back.
	if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("WebDAV server responded with status %d", resp.StatusCode)
	}

	// Return the file name so we can pass it on to the user.
	return fileName, err
}
