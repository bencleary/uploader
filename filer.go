package uploader

import "github.com/google/uuid"

type Upload struct{}

type FilerService interface {
	Record(attachment *Attachment) error
	Fetch(fileUID uuid.UUID) (*Upload, error)
	Delete(fileUID uuid.UUID) error
}
