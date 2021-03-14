package pdf

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
)

const (
	// The DPI used by gopdf seems to be around 128, so we need to multiply the size of
	// the page in pt by 1.778 to get it in pixels.
	// See https://github.com/signintech/gopdf/issues/168
	ptToPx = 1.778
)

var (
	// For now let's hardcode A4 as the page size to use. It might be nice th be able to
	// set it via the configuration file in the future.
	pageSize = gopdf.PageSizeA4
)

// Encode encodes an image into a PDF document, sizing it so it takes the whole page.
func Encode(w io.Writer, img image.Image, jpegEncodeOptions *jpeg.Options) error {
	// Make sure the image isn't too large for the document.
	img = resizeImageForPage(img, pageSize)

	// Start defining the PDF document.
	pdf := new(gopdf.GoPdf)
	pdf.Start(gopdf.Config{PageSize: *pageSize})
	pdf.AddPage()

	// Encode the (maybe) resized image as JPEG.
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, jpegEncodeOptions); err != nil {
		return err
	}

	// Add the image to the document.
	ih, err := gopdf.ImageHolderByReader(buf)
	if err != nil {
		return err
	}
	if err = pdf.ImageByHolder(ih, 0, 0, nil); err != nil {
		return err
	}

	// Compile the PDF document.
	return pdf.Write(w)
}

// resizeImageForPage resizes the given image if this image is too large for the given
// page size, while preserving its aspect ratio.
func resizeImageForPage(img image.Image, pageSize *gopdf.Rect) image.Image {
	// Get the dimensions of the current image.
	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())

	// Calculate the dimensions of the page in pixels.
	pageWidth := pageSize.W * ptToPx
	pageHeight := pageSize.H * ptToPx

	// Check if the image is too large for the page.
	var tooWide, tooHigh bool
	if width > pageWidth {
		tooWide = true
	}
	if height > pageHeight {
		tooHigh = true
	}

	// If the image is both too wide and too high for the page, pick the largest
	// dimension and use it as the reference when resizing, so the image takes as much
	// space as possible in the document.
	if tooWide && tooHigh {
		tooWide = width > height
		tooHigh = height > width
	}

	// Resize the image if needed.
	if tooWide {
		return resize.Resize(uint(pageWidth), 0, img, resize.Lanczos3)
	} else if tooHigh {
		return resize.Resize(0, uint(pageHeight), img, resize.Lanczos3)
	}

	// If no resize was needed, just return the original image.
	return img
}
