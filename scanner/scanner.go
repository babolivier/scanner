package scanner

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tjgq/sane"

	"github.com/babolivier/scanner/config"
	"github.com/babolivier/scanner/pdf"
	"github.com/babolivier/scanner/webdav"
)

const (
	inchToMMRatio = 25.4
)

var (
	// ErrUnsupportedFormat is the error returned by ScanAndUpload if the format isn't
	// among the supported ones.
	ErrUnsupportedFormat = errors.New("Unsupported format")
)

// Scanner interacts with SANE to control the scanner.
type Scanner struct {
	cfg             *config.ScannerConfig
	conn            *sane.Conn
	webDAVClient    *webdav.Client
	defaultScanArea *defaultScanArea
}

// NewScanner returns a new Scanner. It also opens the SANE connection to the scanning
// device, and sets the mode.
func NewScanner(cfg *config.ScannerConfig, client *webdav.Client) (s *Scanner, err error) {
	s = &Scanner{
		cfg:          cfg,
		webDAVClient: client,
	}

	// Try to open a connection with the device.
	if err = s.openConn(); err != nil {
		// If that didn't work, we'll try again when trying to get an image.
		logrus.
			WithField("name", s.cfg.DeviceName).
			WithError(err).
			Warn("Failed to connect to device")
	}

	return s, nil
}

// openConn opens a SANE connection to the scanning device and sets the mode.
func (s *Scanner) openConn() (err error) {
	if s.conn, err = sane.Open(s.cfg.DeviceName); err != nil {
		return err
	}

	if _, err = s.conn.SetOption("mode", s.cfg.Mode); err != nil {
		return err
	}

	logrus.WithField("name", s.cfg.DeviceName).Info("Connected to device")

	return nil
}

// Preview triggers a low-resolution scan on the scanning device and returns the
// resulting image.
func (s *Scanner) Preview() (*sane.Image, error) {
	logrus.Info("Getting preview")

	options := &ScanOptions{
		resolution: s.cfg.PreviewRes,
	}
	return s.getImage(options)
}

// ScanAndUpload triggers a high-resolution scan on the scanning device and uploads the
// resulting image to the WebDAV server.
func (s *Scanner) ScanAndUpload(options *ScanOptions) (fileName string, err error) {
	entry := logrus.WithField("format", options.Format)
	if options.WithRect {
		entry = entry.WithFields(logrus.Fields{
			"x":      options.X,
			"y":      options.Y,
			"width":  options.Width,
			"height": options.Height,
		})
	}
	entry.Info("Triggering scan")

	// Select the encoding function to run the resulting image through, and at the same
	// time make sure the format is a supported one. We do this early because the scan
	// can take some time to complete, and we don't want to wait that long to tell the
	// requester the requested format isn't supported.
	var encode func(w io.Writer, img image.Image, o *jpeg.Options) error
	switch options.Format {
	case "jpeg":
		encode = jpeg.Encode
	case "pdf":
		encode = pdf.Encode
	default:
		return "", ErrUnsupportedFormat
	}

	// Trigger the scan and get the resulting image.
	options.resolution = s.cfg.ScanRes
	img, err := s.getImage(options)
	if err != nil {
		return
	}

	// Encode the resulting image.
	buf := new(bytes.Buffer)
	if err = encode(buf, img, nil); err != nil {
		return "", err
	}

	// Upload the encoded bytes to the WebDAV server.
	return s.webDAVClient.Upload(buf, options.Format)
}

// getImage triggers a scan with the provided resolution on the scanning device.
func (s *Scanner) getImage(options *ScanOptions) (*sane.Image, error) {
	logrus.WithFields(logrus.Fields{
		"resolution": options.resolution,
		"with_rect":  options.WithRect,
	}).Info("Reading image")

	// If the SANE connection hasn't already been established, try to do it now.
	if s.conn == nil {
		if err := s.openConn(); err != nil {
			// If that didn't work, return the error, and try again next time.
			return nil, err
		}
	}

	// Set the scan resolution.
	if _, err := s.conn.SetOption("resolution", options.resolution); err != nil {
		return nil, err
	}

	if options.WithRect {
		// If we're scanning a rectangle within the scanning area (and not the whole
		// area), then set the parameters on the scanner.
		if _, err := s.conn.SetOption("tl-x", s.pxToMM(options.X)); err != nil {
			return nil, err
		}

		if _, err := s.conn.SetOption("tl-y", s.pxToMM(options.Y)); err != nil {
			return nil, err
		}

		if _, err := s.conn.SetOption("br-x", s.pxToMM(options.Width)); err != nil {
			return nil, err
		}

		if _, err := s.conn.SetOption("br-y", s.pxToMM(options.Height)); err != nil {
			return nil, err
		}
	} else {
		// Otherwise reset the options on the SANE connection to include the whole
		// available surface.
		if err := s.resetScanArea(); err != nil {
			return nil, err
		}
	}

	return s.conn.ReadImage()
}

// pxToMM calculates a value in millimeters from the given value in pixels. The preview
// resolution (in DPI) is used because the value in pixels is expected to be coming from
// a rectangle drawn on a preview.
func (s *Scanner) pxToMM(pxValue int) float64 {
	return float64(pxValue) * inchToMMRatio / float64(s.cfg.PreviewRes)
}
