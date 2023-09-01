package uploader

import (
	"context"
	"io"
)

type EncryptionService interface {
	EncryptStream(ctx context.Context, src io.Reader, key string) (io.ReadCloser, error)
	DecryptStream(ctx context.Context, src io.Reader, ket string) (io.ReadCloser, error)
}
