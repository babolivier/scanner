package scanner

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"

	"github.com/tjgq/sane"

	"github.com/babolivier/scanner/config"
	"github.com/babolivier/scanner/pdf"
	"github.com/babolivier/scanner/webdav"
)

var (
	// ErrUnsupportedFormat is the error returned by ScanAndUpload if the format isn't
	// among the supported ones.
	ErrUnsupportedFormat = errors.New("Unsupported format")
)

// Scanner interacts with SANE to control the scanner.
type Scanner struct {
	cfg          *config.ScannerConfig
	conn         *sane.Conn
	webDAVClient *webdav.Client
}

// NewScanner returns a new Scanner. It also opens the SANE connection to the scanning
// device, and sets the mode.
func NewScanner(cfg *config.ScannerConfig, client *webdav.Client) (s *Scanner, err error) {
	s = &Scanner{
		cfg:          cfg,
		webDAVClient: client,
	}

	if s.conn, err = sane.Open(cfg.DeviceName); err != nil {
		return nil, err
	}

	if _, err = s.conn.SetOption("mode", cfg.Mode); err != nil {
		return nil, err
	}

	return
}

// Preview triggers a low-resolution scan on the scanning device and returns the
// resulting image.
func (s *Scanner) Preview() (*sane.Image, error) {
	return s.getImage(s.cfg.PreviewRes)
}

// ScanAndUpload triggers a high-resolution scan on the scanning device and uploads the
// resulting image to the WebDAV server.
func (s *Scanner) ScanAndUpload(format string) (fileName string, err error) {
	// Select the encoding function to run the resulting image through, and at the same
	// time make sure the format is a supported one. We do this early because the scan
	// can take some time to complete, and we don't want to wait that long to tell the
	// requester the requested format isn't supported.
	var encode func(w io.Writer, img image.Image, o *jpeg.Options) error
	switch format {
	case "jpeg":
		encode = jpeg.Encode
	case "pdf":
		encode = pdf.Encode
	default:
		return "", ErrUnsupportedFormat
	}

	// Trigger the scan and get the resulting image.
	img, err := s.getImage(s.cfg.ScanRes)
	if err != nil {
		return
	}

	// Encode the resulting image.
	buf := new(bytes.Buffer)
	if err = encode(buf, img, nil); err != nil {
		return "", err
	}

	// Upload the encoded bytes to the WebDAV server.
	return s.webDAVClient.Upload(buf, format)
}

// getImage triggers a scan with the provided resolution on the scanning device.
func (s *Scanner) getImage(resolution int) (*sane.Image, error) {
	if _, err := s.conn.SetOption("resolution", resolution); err != nil {
		return nil, err
	}

	return s.conn.ReadImage()
}
