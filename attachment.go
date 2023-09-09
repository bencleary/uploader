package uploader

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Attachment struct {
	UID              uuid.UUID
	OwnerID          int
	FileName         string
	FileSize         int64
	Extension        string
	MimeType         string
	LocalPath        string
	PreviewLocalPath string
}

// CreatePreviewLocalPath adds .preview into the localpath befroe the file extenion
func (a *Attachment) CreatePreviewLocalPath() string {
	// Find the last dot (.) in the FileName
	lastDotIndex := strings.LastIndex(a.FileName, ".")

	// If there's no dot or the dot is at the beginning of the file name, return an error or handle it as needed
	if lastDotIndex <= 0 {
		return "" // or return an error
	}

	// Create the preview file name by inserting ".preview" before the last dot
	previewFileName := a.FileName[:lastDotIndex] + ".preview" + a.FileName[lastDotIndex:]

	// Replace the old extension with ".preview"
	previewLocalPath := strings.Replace(a.LocalPath, a.FileName, previewFileName, 1)

	// Update the PreviewLocalPath field
	a.PreviewLocalPath = previewLocalPath

	// Return the new PreviewLocalPath
	return previewLocalPath
}

func (a *Attachment) CopyFileToPath(path string) error {
	if a.LocalPath == "" {
		return Errorf(INVALID, "LocalPath is empty")
	}
	file, err := os.Open(a.LocalPath)
	if err != nil {
		return err
	}
	defer file.Close()
	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err = io.Copy(dest, file); err != nil {
		return err
	}

	return nil

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
