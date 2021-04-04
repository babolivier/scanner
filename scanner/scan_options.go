package scanner

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
	WithRect   bool
	X          int
	Y          int
	Width      int
	Height     int
	resolution int
}

// NewOptionsFromQuery instantiates a new ScanOptions and fills it with the provided
// URL query parameters.
// Returns ErrMissingFormat if the format is missing from the query parameters, or
// ErrMalformedRect if a rectangle is defined in the query parameters but one of its
// parameters is missing or malformed.
func NewOptionsFromQuery(query url.Values) (*ScanOptions, error) {
	options := new(ScanOptions)

	// Make sure a format has been provided, and return an error if not.
	options.Format = query.Get("format")
	if options.Format == "" {
		return nil, ErrMissingFormat
	}

	x := query.Get("x")
	y := query.Get("y")
	width := query.Get("width")
	height := query.Get("height")

	// Don't do anything more if no rectangle was provided.
	if x == "" && y == "" && width == "" && height == "" {
		return options, nil
	}

	options.WithRect = true

	// Check if any of the rectangle parameters is missing.
	if x == "" || y == "" || width == "" || height == "" {
		return nil, ErrMalformedRect
	}

	// Parse the parameters into integers. This is probably not required, and SANE can
	// probably deal with them just as well if they're strings, but this at least
	// provides an extra layer of input validation
	var err error
	if options.X, err = strconv.Atoi(x); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse x value for rectangle")

		return nil, ErrMalformedRect
	}

	if options.Y, err = strconv.Atoi(y); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse y value for rectangle")

		return nil, ErrMalformedRect
	}

	if options.Width, err = strconv.Atoi(width); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse width value for rectangle")

		return nil, ErrMalformedRect
	}

	if options.Height, err = strconv.Atoi(height); err != nil {
		logrus.
			WithError(err).
			Error("Failed to parse height value for rectangle")

		return nil, ErrMalformedRect
	}

	return options, nil
}
