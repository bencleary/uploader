package encryption

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestIsValidKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "valid key with variation",
			key:      "12345678901234567890123456789012",
			expected: true,
		},
		{
			name:     "weak key all zeros",
			key:      "00000000000000000000000000000000",
			expected: false,
		},
		{
			name:     "weak key all ones",
			key:      "11111111111111111111111111111111",
			expected: false,
		},
		{
			name:     "invalid length too short",
			key:      "short",
			expected: false,
		},
		{
			name:     "invalid length too long",
			key:      "123456789012345678901234567890123",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsValidKey(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestNewAESService(t *testing.T) {
	aes := NewAESService(nil)
	if aes == nil {
		t.Fatal("expected aes service")
	}

}

func TestAESEncrypt(t *testing.T) {
	key := "some-secret-key"
	aes := NewAESService(nil)
	_, err := aes.encrypt([]byte("some-data"), key)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAESDecrypt(t *testing.T) {
	key := "some-secret-key1"
	aes := NewAESService(nil)
	enc, err := aes.encrypt([]byte("some-data"), key)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := aes.decrypt(enc, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(dec) != "some-data" {
		t.Fatal("decrypted data does not match")
	}
}

func TestAESServiceEncryptStream(t *testing.T) {
	key := "some-secret-key"
	aes := NewAESService(nil)
	data := bytes.NewBufferString("some-data")
	_, err := aes.EncryptStream(context.Background(), data, key)
	if err == nil {
		t.Fatal(err)
	}
}

func TestAESServiceDecryptStream(t *testing.T) {
	key := "some-secret-key1"
	aes := NewAESService(nil)
	data := bytes.NewBufferString("some-data")
	enc, err := aes.EncryptStream(context.Background(), data, key)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := aes.DecryptStream(context.Background(), enc, key)
	if err != nil {
		t.Fatal(err)
	}
	content, err := io.ReadAll(dec)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "some-data" {
		t.Fatal("decrypted data does not match")
	}
}
