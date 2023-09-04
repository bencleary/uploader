package uploader

import (
	"mime/multipart"
	"path/filepath"
	"strings"

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

func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	ext = strings.TrimPrefix(ext, ".")
	return ext
}

func NewAttachment(file *multipart.FileHeader, requestUser int) *Attachment {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil
	}
	return &Attachment{
		UID:       uid,
		OwnerID:   requestUser,
		FileSize:  file.Size,
		FileName:  file.Filename,
		Extension: getFileExtension(file.Filename),
		MimeType:  file.Header.Get("Content-Type"),
	}
}
