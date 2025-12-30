package storage

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/bencleary/uploader"
	"github.com/bencleary/uploader/internal/encryption"
	"github.com/bencleary/uploader/internal/keystore"
)

// createS3Storage creates an S3Storage instance for testing
func createS3Storage(t *testing.T) *S3Storage {
	t.Helper()
	keyStore := keystore.NewInMemoryKeyStore()
	aes := encryption.NewAESService(keyStore)
	prefix := t.TempDir()

	storage := NewS3Storage(&S3Options{
		Endpoint:       "http://localhost:9000",
		Bucket:         "uploader",
		Region:         "us-east-1",
		Prefix:         prefix,
		ForcePathStyle: true,
		AccessKeyID:    "minioadmin",
		SecretKey:      "minioadmin",
	}, aes)

	if storage == nil {
		t.Fatal("expected S3 storage to be created")
	}

	return storage
}

// skipIfS3Unavailable skips the test if S3/MinIO is not available
func skipIfS3Unavailable(t *testing.T) {
	t.Helper()
	// Try to create a storage instance - if it fails, S3 is likely not available
	keyStore := keystore.NewInMemoryKeyStore()
	aes := encryption.NewAESService(keyStore)
	storage := NewS3Storage(&S3Options{
		Endpoint:       "http://localhost:9000",
		Bucket:         "uploader",
		Region:         "us-east-1",
		Prefix:         t.TempDir(),
		ForcePathStyle: true,
		AccessKeyID:    "minioadmin",
		SecretKey:      "minioadmin",
	}, aes)

	if storage == nil {
		t.Skip("S3/MinIO not available - skipping integration test")
	}

	// Try to initialize - if this fails, S3 is not running
	if err := storage.Initialise(context.Background()); err != nil {
		t.Skipf("S3/MinIO not available - skipping integration test: %v", err)
	}
}

func TestNewS3Storage(t *testing.T) {
	keyStore := keystore.NewInMemoryKeyStore()
	aes := encryption.NewAESService(keyStore)
	prefix := t.TempDir()

	storage := NewS3Storage(&S3Options{
		Endpoint:       "http://localhost:9000",
		Bucket:         "uploader",
		Region:         "us-east-1",
		Prefix:         prefix,
		ForcePathStyle: true,
		AccessKeyID:    "minioadmin",
		SecretKey:      "minioadmin",
	}, aes)

	if storage == nil {
		t.Fatal("expected storage to be created")
	}
}

func TestNewS3Storage_NilOptions(t *testing.T) {
	keyStore := keystore.NewInMemoryKeyStore()
	aes := encryption.NewAESService(keyStore)

	storage := NewS3Storage(nil, aes)
	if storage != nil {
		t.Fatal("expected storage to be nil when options is nil")
	}
}

func TestNewS3Storage_NilEncryption(t *testing.T) {
	storage := NewS3Storage(&S3Options{
		Endpoint:       "http://localhost:9000",
		Bucket:         "uploader",
		Region:         "us-east-1",
		Prefix:         t.TempDir(),
		ForcePathStyle: true,
		AccessKeyID:    "minioadmin",
		SecretKey:      "minioadmin",
	}, nil)

	if storage != nil {
		t.Fatal("expected storage to be nil when encryption is nil")
	}
}

func TestS3StorageInitialise(t *testing.T) {
	storage := createS3Storage(t)

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestS3StorageHold(t *testing.T) {
	storage := createS3Storage(t)

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	var pngBuffer bytes.Buffer
	if err := png.Encode(&pngBuffer, img); err != nil {
		t.Fatal(err)
	}

	// Create multipart file header
	fileHeader := createMultipartFileHeader(t, "file", "test.png", pngBuffer.Bytes())

	attachment, err := storage.Hold(context.Background(), fileHeader)
	if err != nil {
		t.Fatal(err)
	}

	if attachment == nil {
		t.Fatal("expected attachment to be created")
	}

	if attachment.UID.String() == "" {
		t.Fatal("expected attachment to have a UID")
	}

	if attachment.LocalPath == "" {
		t.Fatal("expected attachment to have a LocalPath")
	}

	// Verify file exists
	if _, err := os.Stat(attachment.LocalPath); err != nil {
		t.Fatalf("expected file to exist at %s: %v", attachment.LocalPath, err)
	}
}

func TestS3StorageUpload(t *testing.T) {
	skipIfS3Unavailable(t)

	storage := createS3Storage(t)
	keyStore := keystore.NewInMemoryKeyStore()
	encryptionKey := "thisisa32byteencryptionkey123456" // 32 bytes

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Create a test file
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create attachment
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

	// Store encryption key
	if err := keyStore.StoreKey(attachment.UID.String(), []byte(encryptionKey)); err != nil {
		t.Fatal(err)
	}

	// Upload to S3
	if err := storage.Upload(context.Background(), attachment, encryptionKey); err != nil {
		t.Fatalf("failed to upload: %v", err)
	}
}

func TestS3StorageUpload_NilAttachment(t *testing.T) {
	storage := createS3Storage(t)

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	err := storage.Upload(context.Background(), nil, "key")
	if err == nil {
		t.Fatal("expected error when attachment is nil")
	}

	uploaderErr, ok := err.(*uploader.Error)
	if !ok || uploaderErr.Code != uploader.INVALID {
		t.Fatalf("expected INVALID error, got: %v", err)
	}
}

func TestS3StorageDownload(t *testing.T) {
	skipIfS3Unavailable(t)

	storage := createS3Storage(t)
	keyStore := keystore.NewInMemoryKeyStore()
	encryptionKey := "thisisa32byteencryptionkey123456" // 32 bytes

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Create and upload a file first
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

	if err := keyStore.StoreKey(attachment.UID.String(), []byte(encryptionKey)); err != nil {
		t.Fatal(err)
	}

	// Upload to S3
	if err := storage.Upload(context.Background(), attachment, encryptionKey); err != nil {
		t.Fatalf("failed to upload: %v", err)
	}

	// Download from S3
	reader, err := storage.Download(context.Background(), attachment, false, encryptionKey)
	if err != nil {
		t.Fatalf("failed to download: %v", err)
	}
	defer reader.Close()

	// Read the content
	content, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("failed to read downloaded content: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("expected downloaded content to not be empty")
	}
}

func TestS3StorageDownload_NilAttachment(t *testing.T) {
	storage := createS3Storage(t)

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	_, err := storage.Download(context.Background(), nil, false, "key")
	if err == nil {
		t.Fatal("expected error when attachment is nil")
	}

	uploaderErr, ok := err.(*uploader.Error)
	if !ok || uploaderErr.Code != uploader.INVALID {
		t.Fatalf("expected INVALID error, got: %v", err)
	}
}

func TestS3StorageDelete(t *testing.T) {
	skipIfS3Unavailable(t)

	storage := createS3Storage(t)
	keyStore := keystore.NewInMemoryKeyStore()
	encryptionKey := "thisisa32byteencryptionkey123456" // 32 bytes

	if err := storage.Initialise(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Create and upload a file first
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

	if err := keyStore.StoreKey(attachment.UID.String(), []byte(encryptionKey)); err != nil {
		t.Fatal(err)
	}

	// Upload to S3
	if err := storage.Upload(context.Background(), attachment, encryptionKey); err != nil {
		t.Fatalf("failed to upload: %v", err)
	}

	// Delete from S3
	if err := storage.Delete(context.Background(), attachment.UID.String()); err != nil {
		t.Fatalf("failed to delete: %v", err)
	}
}

func TestS3StorageObjectKey(t *testing.T) {
	storage := createS3Storage(t)

	tests := []struct {
		name      string
		uid       string
		isPreview bool
		prefix    string
		expected  string
	}{
		{
			name:      "main file without prefix",
			uid:       "test-uid",
			isPreview: false,
			prefix:    "",
			expected:  "test-uid.enc",
		},
		{
			name:      "preview file without prefix",
			uid:       "test-uid",
			isPreview: true,
			prefix:    "",
			expected:  "test-uid.preview.enc",
		},
		{
			name:      "main file with prefix",
			uid:       "test-uid",
			isPreview: false,
			prefix:    "uploader",
			expected:  "uploader/test-uid.enc",
		},
		{
			name:      "preview file with prefix",
			uid:       "test-uid",
			isPreview: true,
			prefix:    "uploader",
			expected:  "uploader/test-uid.preview.enc",
		},
		{
			name:      "main file with prefix ending in slash",
			uid:       "test-uid",
			isPreview: false,
			prefix:    "uploader/",
			expected:  "uploader/test-uid.enc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage.options.Prefix = tt.prefix
			result := storage.objectKey(tt.uid, tt.isPreview)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
