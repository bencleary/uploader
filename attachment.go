package uploader

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type Attachment struct {
	UID              uuid.UUID
	Contents         multipart.File
	OwnerID          int
	FileName         string
	FileSize         int64
	Extension        string
	MimeType         string
	LocalPath        string
	PreviewLocalPath string
}

type AttachmentService interface{}
