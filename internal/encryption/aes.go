package encryption

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/bencleary/uploader"
)

// Encryption
// var encryptionKey = []byte("thisisa32byteencryptionkey123456")

// Ensure service implements interface.
var _ uploader.EncryptionService = (*AES)(nil)

type AES struct {
	keystore uploader.KeyStoreService
}

func NewAESService(keystore uploader.KeyStoreService) *AES {
	return &AES{
		keystore: keystore,
	}
}

func (a *AES) encrypt(data []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (a *AES) decrypt(ciphertext []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (a *AES) EncryptStream(ctx context.Context, src io.Reader, key string) (io.ReadCloser, error) {
	encryptedContent, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	encryptedData, err := a.encrypt(encryptedContent, key)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(encryptedData)), nil
}
func (a *AES) DecryptStream(ctx context.Context, src io.Reader, key string) (io.ReadCloser, error) {
	decryptedContent, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	decryptedData, err := a.decrypt(decryptedContent, key)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(decryptedData)), nil
}
