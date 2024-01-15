package encryption_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/bencleary/uploader/internal/encryption"
)

func TestNewAESService(t *testing.T) {
	aes := encryption.NewAESService(nil)
	if aes == nil {
		t.Fatal("expected aes service")
	}

}

func TestAESServiceEncryptStream(t *testing.T) {
	key := "some-secret-key"
	aes := encryption.NewAESService(nil)
	data := bytes.NewBufferString("some-data")
	_, err := aes.EncryptStream(context.Background(), data, key)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAESServiceDecryptStream(t *testing.T) {
	key := "some-secret-key1"
	aes := encryption.NewAESService(nil)
	data := bytes.NewBufferString("some-data")
	enc, err := aes.EncryptStream(context.Background(), data, key)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := aes.DecryptStream(context.Background(), enc, key)
	if err != nil {
		t.Fatal("expected error")
	}
	content, err := io.ReadAll(dec)
	if err != nil {
		t.Fatal("expected error")
	}

	if string(content) != "some-data" {
		t.Fatal("expected data")
	}
}
