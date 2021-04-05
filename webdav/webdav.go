package webdav

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/babolivier/scanner/common"
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
func (c *Client) Upload(options *common.ScanOptions, body *bytes.Buffer) (string, error) {
	// Determine the file's name.
	fileNameNoExt := time.Now().Format("2006-01-02_15-04-05")
	if options.FileName != "" {
		fileNameNoExt = options.FileName
	}

	// Add the format extension to the file's name.
	fileName := fmt.Sprintf("%s.%s", fileNameNoExt, options.Format)

	logrus.
		WithField("filename", fileName).
		Info("Uploading file to the WebDAV server")

	// Upload the file.
	status, err := c.requestFile(http.MethodPut, fileName, body)
	if err != nil {
		return "", err
	}

	logrus.WithField("status_code", status).Info("Upload finished")

	// According to RFC4918, the creation of a resource must be indicated by use of a
	// 201 Created response code, so return an error if that's not what we got back.
	if status != http.StatusCreated {
		err = fmt.Errorf("WebDAV server responded with status %d", status)
	}

	// Return the file name so we can pass it on to the user.
	return fileName, err
}

// FileExists checks if a file already exists with the given name.
func (c *Client) FileExists(fileName string) (bool, error) {
	// Send a HEAD request with the file name, if the server responds with a 200 status
	// then a file with this name exists, if the status is 404 then it doesn't.
	status, err := c.requestFile(http.MethodHead, fileName, nil)
	if err != nil {
		return false, err
	}

	logrus.WithField("status_code", status).Info("Checked file existence")

	// Make sure we don't accidentally understand an error as the file not existing.
	if status != http.StatusOK && status != http.StatusNotFound {
		return false, fmt.Errorf("WebDAV server responded with unexpected status %d", status)
	}

	return status == http.StatusOK, nil
}

// requestFile sends a HTTP request to the WebDAV server for the given path with the given
// method and body.
func (c *Client) requestFile(method string, fileName string, body io.Reader) (int, error) {
	// Parse the root URL. Ideally we'd do this in NewClient, but we need to change the
	// path of this URL with the file's name, and we don't want this change to persist
	// on the client.
	u, err := url.Parse(c.cfg.RootURL)
	if err != nil {
		return 0, err
	}

	// Build a path that includes the full path for this file on the WebDAV server.
	u.Path = path.Join(u.Path, c.cfg.UploadPath, fileName)

	// Create the request.
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return 0, err
	}

	// Add basic auth to the request.
	req.SetBasicAuth(c.cfg.User, c.cfg.Password)

	// Send the request.
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
