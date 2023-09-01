package uploader

import (
	"context"

	"github.com/google/uuid"
)

type StorageService interface {
	Initialise(ctx context.Context) error
	Hold(ctx context.Context, attachment *Attachment) (string, error)
	Save(ctx context.Context, vaultPath string, attachment *Attachment, key string) error
	Load(ctx context.Context, fileUID uuid.UUID) (*Attachment, error)
	Delete(ctx context.Context, fileUID uuid.UUID) error
}
