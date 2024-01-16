package scaler

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"testing"
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

func TestNewDrawImageScaler(t *testing.T) {
	draw_scaler := NewDrawImageScaler([]string{"image/png"})
	if draw_scaler == nil {
		t.Fatal("expected draw scaler")
	}
}

func TestDrawImageScalerSupported(t *testing.T) {
	draw_scaler := NewDrawImageScaler([]string{"image/png"})
	if !draw_scaler.Supported("image/png") {
		t.Fatal("expected image/png")
	}
}

func TestDrawImageScalerNotSupported(t *testing.T) {
	draw_scaler := NewDrawImageScaler([]string{"image/png"})
	if draw_scaler.Supported("image/jpeg") {
		t.Fatal("expected image/png")
	}
}

func TestGetImageFormat(t *testing.T) {

	defer cleanUpFiles("test.png")

	generateEmptyImage("test", 1, 1)

	buffer := make([]byte, 512)
	opened, err := os.Open("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer opened.Close()

	_, err = opened.Read(buffer)
	if err != nil {
		t.Fatal(err)
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/png" {
		t.Fatal("expected image/png")
	}

}

func TestDrawImageScalerScale(t *testing.T) {
	defer cleanUpFiles("test.png")
	draw_scaler := NewDrawImageScaler([]string{"image/png"})
	generateEmptyImage("test", 1000, 1000)
	err := draw_scaler.Scale(context.Background(), "test.png", 100, "image/png")
	if err != nil {
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
	if decoded.Bounds().Dx() != 100 {
		t.Fatal("expected 100 width")
	}
	if decoded.Bounds().Dy() != 100 {
		t.Fatal("expected 100 height")
	}
}
