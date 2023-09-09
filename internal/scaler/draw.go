package scaler

import (
	"context"
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

var (
	_ uploader.ScalerService = (*DrawImageScaler)(nil)
)

type DrawImageScaler struct {
	supported []string
}

// NewDrawImageScaler creates a new instance of DrawImageScaler.
// It takes a list of supported MIME types and returns a pointer to the new instance.
func NewDrawImageScaler(supportedMimeTypes []string) *DrawImageScaler {
	return &DrawImageScaler{
		supported: supportedMimeTypes,
	}
}

// Supported checks if the given MIME type is supported by the DrawImageScaler.
func (d *DrawImageScaler) Supported(mimeType string) bool {
	// Use the Contains function from the slices package to check if the mimeType is in the supported slice.
	return slices.Contains(d.supported, mimeType)
}

// getImageFormat determines the appropriate image format based on the provided MIME type.
// It returns an instance of the uploader.ImageFormat interface and an error.
func getImageFormat(mimeType string) (uploader.ImageFormat, error) {
	switch mimeType {
	case PNG:
		return &uploader.PNGFormat{}, nil
	case GIF:
		return &uploader.GIFFormat{}, nil
	case JPEG:
		return &uploader.JPEGFormat{}, nil
	}

	return nil, uploader.Errorf(uploader.INVALID, "Unsupported image format: %s", mimeType)
}

// Scale resizes an image to a target width while maintaining the aspect ratio.
// It takes a context, an attachment, and the target width as input parameters.
// It returns an error if there was any issue during the scaling process.
func (d *DrawImageScaler) Scale(ctx context.Context, filePath string, targetWidth int, mimeType string) error {
	// Get the image format from the attachment's MIME type
	format, err := getImageFormat(mimeType)
	if err != nil {
		return err
	}

	// Open the input image file
	input, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer input.Close()

	// Decode the input image
	src, _, err := image.Decode(input)
	if err != nil {
		return err
	}

	// Only resize if the width of the source image is larger than the target width
	if src.Bounds().Dx() <= targetWidth {
		return nil
	}

	// Calculate the target height to maintain the aspect ratio
	aspectRatio := float64(src.Bounds().Dx()) / float64(src.Bounds().Dy())
	targetHeight := int(float64(targetWidth) / aspectRatio)

	// Create the output image file
	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()

	// Create the destination image with the calculated size
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Resize the image using the BiLinear algorithm
	draw.BiLinear.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode the resized image and write it to the output file
	return format.Encode(output, dst)
}
