package storage

import (
	"context"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"os"
	"testing"

	"github.com/bencleary/uploader/internal/encryption"
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

func TestLocalStorage(t *testing.T) {
	aes := encryption.NewAESService(nil)
	storage := NewLocalStorage("upload", "vault", aes)
	if storage == nil {
		t.Fatal("expected local storage")
	}
}

func TestLocalStorageInitialise(t *testing.T) {
	storage := NewLocalStorage("upload", "vault", nil)
	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestLocalStorageHold(t *testing.T) {

	defer cleanUpFiles("temp/test.png")

	aes := encryption.NewAESService(nil)
	storage := NewLocalStorage("temp", "vault", aes)

	generateEmptyImage("temp/test", 1, 1)

	opened, err := os.Open("temp/test.png")

	if err != nil {
		t.Fatal(err)
	}
	defer opened.Close()

	fileHeader := createFileHeader(opened)

	if _, err := storage.Hold(context.Background(), fileHeader); err != nil {
		t.Fatal(err)
	}

}

func createFileHeader(file *os.File) *multipart.FileHeader {
	fileInfo, _ := file.Stat()
	fileHeader := &multipart.FileHeader{
		Filename: fileInfo.Name(),
		Size:     fileInfo.Size(),
		Header:   make(map[string][]string),
	}
	return fileHeader
}
