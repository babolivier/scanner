package http

import (
	"fmt"
	"image/jpeg"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"github.com/tjgq/sane"

	"github.com/babolivier/scanner/common"
	"github.com/babolivier/scanner/config"
	"github.com/babolivier/scanner/scanner"
	"github.com/babolivier/scanner/webdav"
)

const (
	internalErrorMsg = "Something happened"
)

// handlers define the HTTP handlers to serve on top of the static files.
type handlers struct {
	scanner *scanner.Scanner
	webdav  *webdav.Client
}

// handlePanics recovers from a panic that occurred when processing a request, sends a
// 500 response to the request and logs the panic reason as well as the stack trace.
func handlePanics(w http.ResponseWriter) {
	if err := recover(); err != nil {
		logrus.WithField("err", err).Error("Recovering from panic")
		debug.PrintStack()
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
	}
}

// handlePreview generates a JPEG preview of what's currently on the scanner's plate.
func (h *handlers) handlePreview(w http.ResponseWriter, req *http.Request) {
	defer handlePanics(w)

	// Given the endpoint looks like a static image, browsers might try to cache it, but
	// we don't want that.
	w.Header().Add("Cache-Control", "no-cache")

	// Generate the preview.
	img, err := h.scanner.Preview()
	if err != nil {
		logrus.WithError(err).Error("Failed to get preview from scanner")
		if err == sane.ErrBusy {
			http.Error(w, "Device busy", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}

	// Encode the image as JPEG and send it in the response.
	w.Header().Add("Content-Type", "image/jpeg")
	if err = jpeg.Encode(w, img, nil); err != nil {
		logrus.WithError(err).Error("Failed to encode into JPEG")
		http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		return
	}
}

// handleScan generates a scan of what's currently on the scanner's plate and uploads it
// to the WebDAV server.
func (h *handlers) handleScan(w http.ResponseWriter, req *http.Request) {
	defer handlePanics(w)

	// Tell browsers not to cache this endpoint.
	w.Header().Add("Cache-Control", "no-cache")

	// Try to parse the URL query parameters.
	options, err := common.NewOptionsFromQuery(req.URL.Query())
	if err == common.ErrMissingFormat {
		http.Error(w, "Missing format", http.StatusBadRequest)
		return
	} else if err == common.ErrMalformedRect {
		http.Error(w, "Missing or malformed rect arguments", http.StatusBadRequest)
		return
	} else if err != nil {
		logrus.WithError(err).Error("Failed to parse URL query")
	}

	// If a file name has been provided, check that it's not already used by another file.
	if options.FileName != "" {
		exists, err := h.webdav.FileExists(options.FileName)
		if err != nil {
			http.Error(w, internalErrorMsg, http.StatusInternalServerError)
			return
		}

		if exists {
			http.Error(w, "File name already in use", http.StatusConflict)
			return
		}
	}

	// Scan the file and upload it, and get the name of the file that's been uploaded to
	// the WebDAV server.
	fileName, err := h.scanner.ScanAndUpload(options)
	if err != nil {
		logrus.
			WithError(err).
			Error("Failed to scan or to upload to the WebDAV server")

		if err == scanner.ErrUnsupportedFormat {
			http.Error(w, "Unsupported format", http.StatusBadRequest)
		} else {
			http.Error(w, internalErrorMsg, http.StatusInternalServerError)
		}
		return
	}

	// Send the file name back to the client.
	w.WriteHeader(200)
	if _, err = w.Write([]byte(fileName)); err != nil {
		logrus.WithError(err).Error("Failed to respond to /scan request")
	}
}

// ListenAndServe registers the HTTP handlers and starts the HTTP server.
func ListenAndServe(cfg *config.HTTPConfig, s *scanner.Scanner, c *webdav.Client) error {
	h := &handlers{
		scanner: s,
		webdav:  c,
	}

	// Register a file server to serve the front end.
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	// Register the handlers to preview and scan documents.
	http.HandleFunc("/preview.jpg", h.handlePreview)
	http.HandleFunc("/scan", h.handleScan)

	// Figure out which address to listen on, and whether to enable TLS.
	addr := fmt.Sprintf("%s:%s", cfg.Address, cfg.Port)
	useTLS := cfg.TLSCert != "" && cfg.TLSKey != ""

	logrus.WithFields(logrus.Fields{
		"address": addr,
		"use_tls": useTLS,
	}).Info("Started HTTP(S) server")

	// If TLS credentials have been provided, start a HTTPS server, otherwise start a
	// plain text HTTP server.
	if useTLS {
		return http.ListenAndServeTLS(addr, cfg.TLSCert, cfg.TLSKey, nil)
	} else {
		return http.ListenAndServe(addr, nil)
	}
}
