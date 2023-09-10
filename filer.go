package uploader

import "github.com/google/uuid"

type FilerService interface {
	Record(attachment *Attachment) error
	Fetch(fileUID uuid.UUID) (*Attachment, error)
	Delete(fileUID uuid.UUID) error
}
