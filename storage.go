package uploader

import (
	"context"
	"io"
	"mime/multipart"
)

type StorageService interface {
	Initialise(ctx context.Context) error
	Hold(ctx context.Context, attachment *multipart.FileHeader) (*Attachment, error)
	Upload(ctx context.Context, attachment *Attachment, key string) error
	Download(ctx context.Context, attachment *Attachment, preview bool, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, attachmentUID string) error
}
