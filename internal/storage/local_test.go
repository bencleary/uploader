package storage

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bencleary/uploader/internal/encryption"
)

func createMultipartFileHeader(t *testing.T, fieldName, filename string, contents []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(part, bytes.NewReader(contents)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "http://example.com/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(32 << 20); err != nil {
		t.Fatal(err)
	}

	files := req.MultipartForm.File[fieldName]
	if len(files) != 1 {
		t.Fatalf("expected 1 file for %q, got %d", fieldName, len(files))
	}
	return files[0]
}

func TestLocalStorage(t *testing.T) {
	aes := encryption.NewAESService(nil)
	uploadDir := t.TempDir()
	vaultDir := t.TempDir()
	storage := NewLocalStorage(uploadDir, vaultDir, aes)
	if storage == nil {
		t.Fatal("expected local storage")
	}
}

func TestLocalStorageInitialise(t *testing.T) {
	aes := encryption.NewAESService(nil)
	base := t.TempDir()
	uploadDir := filepath.Join(base, "upload")
	vaultDir := filepath.Join(base, "vault")
	storage := NewLocalStorage(uploadDir, vaultDir, aes)
	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(uploadDir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(vaultDir); err != nil {
		t.Fatal(err)
	}
}

func TestLocalStorageHold(t *testing.T) {
	aes := encryption.NewAESService(nil)
	uploadDir := t.TempDir()
	vaultDir := t.TempDir()
	storage := NewLocalStorage(uploadDir, vaultDir, aes)

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var pngBuffer bytes.Buffer
	if err := png.Encode(&pngBuffer, img); err != nil {
		t.Fatal(err)
	}
	fileHeader := createMultipartFileHeader(t, "file", "test.png", pngBuffer.Bytes())

	attachment, err := storage.Hold(context.Background(), fileHeader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(attachment.LocalPath); err != nil {
		t.Fatal(err)
	}

}
