package common

import (
	"errors"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	ErrMissingFormat = errors.New("missing format")
	ErrMalformedRect = errors.New("malformed rect")
)

// ScanOptions stores the parameters to use when scanning an image and processing the
// result.
type ScanOptions struct {
	Format     string
	ScanArea   *ScanArea
	FileName   string
	Resolution int
}

// NewOptionsFromQuery instantiates a new ScanOptions and fills it with the provided
// URL query parameters.
// Returns ErrMissingFormat if the format is missing from the query parameters, or
// ErrMalformedRect if a rectangle is defined in the query parameters but one of its
// parameters is missing or malformed.
func NewOptionsFromQuery(query url.Values) (*ScanOptions, error) {
	options := &ScanOptions{
		Format:   query.Get("format"),
		FileName: query.Get("name"),
	}

	// Make sure a format has been provided, and return an error if not.
	if options.Format == "" {
		return nil, ErrMissingFormat
	}

	x := query.Get("x")
	y := query.Get("y")
	rawWidth := query.Get("width")
	rawHeight := query.Get("height")

	// Don't do anything more if no rectangle was provided.
	if x == "" && y == "" && rawWidth == "" && rawHeight == "" {
		return options, nil
	}

	options.ScanArea = new(ScanArea)

	// Check if any of the rectangle parameters is missing.
	if x == "" || y == "" || rawWidth == "" || rawHeight == "" {
		return nil, ErrMalformedRect
	}

	// Parse the parameters into integers. This is probably not required, and SANE can
	// probably deal with them just as well if they're strings, but this at least
	// provides an extra layer of input validation
	var err error
	if options.ScanArea.TLX, err = strconv.Atoi(x); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse x value for rectangle")

		return nil, ErrMalformedRect
	}

	if options.ScanArea.TLY, err = strconv.Atoi(y); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse y value for rectangle")

		return nil, ErrMalformedRect
	}

	var width, height int
	if width, err = strconv.Atoi(rawWidth); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse width value for rectangle")

		return nil, ErrMalformedRect
	}

	if height, err = strconv.Atoi(rawHeight); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse height value for rectangle")

		return nil, ErrMalformedRect
	}

	options.ScanArea.BRX = options.ScanArea.TLX.(int) + width
	options.ScanArea.BRY = options.ScanArea.TLY.(int) + height

	return options, nil
}
