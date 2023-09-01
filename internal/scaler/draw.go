package scaler

import (
	"image"
	"os"

	"github.com/bencleary/uploader"
	"golang.org/x/exp/slices"
	"golang.org/x/image/draw"
)

const (
	PNG  = "image/png"
	GIF  = "image/gif"
	JPEG = "image/jpeg"
)

type DrawImageScaler struct {
	supported []string
}

func NewDrawImageScaler(supportedMimeTypes []string) *DrawImageScaler {
	return &DrawImageScaler{
		supported: supportedMimeTypes,
	}
}

func (d *DrawImageScaler) Supported(mimeType string) bool {
	return slices.Contains(d.supported, mimeType)
}

func getImageFormat(mimeType string) (uploader.ImageFormat, error) {
	// Determine the appropriate image format based on the provided MIME type.
	// Returns an instance of the resize.ImageFormat interface and an error.

	switch mimeType {
	case PNG:
		return &uploader.PNGFormat{}, nil
	case GIF:
		return &uploader.GIFFormat{}, nil
	case JPEG:
		return &uploader.JPEGFormat{}, nil
	}

	// If the MIME type does not match any of the cases, return nil for the format and error.
	return nil, nil
}

func (d *DrawImageScaler) Scale(inputPath string, targetWidth int) error {
	format, err := getImageFormat("")
	if err != nil {
		return err
	}
	// Open input image file
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	// Decode the input image
	src, _, err := image.Decode(input)
	if err != nil {
		return err
	}

	// only resize if over the target width
	if src.Bounds().Dx() <= targetWidth {
		return nil
	}

	// Calculate the target height to maintain aspect ratio
	aspectRatio := float64(src.Bounds().Dx()) / float64(src.Bounds().Dy())
	targetHeight := int(float64(targetWidth) / aspectRatio)

	// Create output image file
	output, err := os.Create(inputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	// Create the destination image with the calculated size
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Resize the image using the BiLinear algorithm
	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode the resized image and write to output file
	return format.Encode(output, dst)
}
