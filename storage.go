package uploader

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type StorageService interface {
	Initialise(ctx context.Context) error
	Hold(ctx context.Context, attachment *multipart.FileHeader) (*Attachment, error)
	Save(ctx context.Context, attachment *Attachment, key string) error
	Load(ctx context.Context, fileUID uuid.UUID) (*Attachment, error)
	Delete(ctx context.Context, fileUID uuid.UUID) error
}
