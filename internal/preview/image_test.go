package preview

import (
	"context"
	"testing"

	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/bencleary/uploader"
	"github.com/bencleary/uploader/internal/scaler"
)

const (
	testImageWidth  = 600
	testImageHeight = 600
	previewWidth    = 300
)

func generateEmptyImage(name string, width, height int) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with a transparent color
	clearColor := color.RGBA{0, 0, 0, 0}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, clearColor)
		}
	}

	// Create a file to save the image
	file, err := os.Create("" + name + ".png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image to PNG and write it to the file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func cleanUpFiles(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		os.Remove("test.png")
	}
}

func TestImagePreviewService(t *testing.T) {
	scaler := scaler.NewDrawImageScaler([]string{"image/jpeg", "image/png"})
	preview := NewImagePreviewGenerator(scaler)
	if preview == nil {
		t.Fatal("expected image preview generator")
	}
}

// TestImagePreviewGeneration tests the image preview generation functionality.
//
// This test ensures that an image preview can be generated, by creating an empty
// image file, the test resizes and then we confirm it has been resized.
func TestImagePreviewGeneration(t *testing.T) {

	defer cleanUpFiles("test.png")

	imageScaler := scaler.NewDrawImageScaler([]string{"image/png"})
	preview := NewImagePreviewGenerator(imageScaler)
	ctx := context.Background()

	generateEmptyImage("test", testImageWidth, testImageHeight)

	attachment := &uploader.Attachment{}
	attachment.MimeType = "image/png"
	attachment.PreviewLocalPath = "./test.png"

	if err := preview.Generate(ctx, attachment, previewWidth); err != nil {
		t.Fatal(err)
	}

	opened, err := os.Open("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer opened.Close()

	decoded, _, err := image.Decode(opened)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds().Dx() != 300 {
		t.Fatal("expected 300 width")
	}
	if decoded.Bounds().Dy() != 300 {
		t.Fatal("expected 300 height")
	}

}
