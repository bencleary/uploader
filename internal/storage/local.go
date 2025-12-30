package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/bencleary/uploader"
	"github.com/google/uuid"
)

const (
	// DIRECTORY_PERMISSIONS represents the directory permission mode.
	DIRECTORY_PERMISSIONS = 0755
)

// LocalStorage implements the uploader.StorageService interface for storing files locally.
type LocalStorage struct {
	directory  string
	vault      string
	encryption uploader.EncryptionService
}

// NewLocalStorage creates a new instance of LocalStorage.
func NewLocalStorage(uploadPath, vaultPath string, encryptionService uploader.EncryptionService) *LocalStorage {
	if encryptionService == nil {
		return nil
	}

	return &LocalStorage{
		directory:  uploadPath,
		vault:      vaultPath,
		encryption: encryptionService,
	}
}

// Initialise creates the necessary directories if they don't exist.
func (l *LocalStorage) Initialise(ctx context.Context) error {
	directories := []string{l.directory, l.vault}
	for _, dir := range directories {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dir, DIRECTORY_PERMISSIONS)
			if err != nil {
				return err
			}
			fmt.Println("Directory created:", dir)
			continue
		} else if err != nil {
			return err
		}
	}
	return nil
}

// Hold stores the uploaded file in the vault directory.
func (l *LocalStorage) Hold(ctx context.Context, file *multipart.FileHeader) (*uploader.Attachment, error) {
	vaultID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	vaultDir := filepath.Join(l.vault, vaultID.String())

	_, err = os.Stat(vaultDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(vaultDir, DIRECTORY_PERMISSIONS)
		if err != nil {
			return nil, err
		}
		fmt.Println("Directory created:", vaultDir)
	} else if err != nil {
		return nil, err
	}

	vaultPath := filepath.Join(vaultDir, file.Filename)
	dst, err := os.Create(vaultPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	contents, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer contents.Close()

	if _, err = io.Copy(dst, contents); err != nil {
		return nil, err
	}

	attachment := uploader.NewAttachment(file, 1)
	attachment.LocalPath = vaultPath

	return attachment, nil
}

// Download retrieves an encrypted attachment from the local storage, decrypts it using the specified key,
// and returns a ReadCloser for the decrypted content.
//
// Parameters:
//   - ctx: The context.Context to control the execution flow and cancellation.
//   - attachment: The uploader.Attachment representing the attachment to be downloaded.
//   - preview: A boolean indicating whether to download the preview version of the attachment.
//   - key: The encryption key used to decrypt the attachment.
//
// Returns:
//   - io.ReadCloser: A ReadCloser providing access to the decrypted content of the attachment.
//   - error: An error, if any, encountered during the download or decryption process.
//
// Example:
//
//	localStore := NewLocalStorage("/path/to/storage", encryptionProvider)
//	attachment := &uploader.Attachment{UID: attachmentUID}
//	key := "encryption_key"
//	preview := false
//	reader, err := localStore.Download(context.Background(), attachment, preview, key)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//	Use the reader to access the decrypted content.
func (l *LocalStorage) Download(ctx context.Context, attachment *uploader.Attachment, preview bool, key string) (io.ReadCloser, error) {
	storagePath := filepath.Join(l.directory, attachment.UID.String())

	var filePath string
	if preview {
		filePath = filepath.Join(storagePath, attachment.UID.String()) + ".preview.enc"
	} else {
		filePath = filepath.Join(storagePath, attachment.UID.String()) + ".enc"
	}

	source, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer source.Close()

	decrypted, err := l.encryption.DecryptStream(ctx, source, key)
	if err != nil {
		return nil, err
	}
	defer decrypted.Close()

	return decrypted, nil
}

// Delete removes a folder and its contents based on its unique identifier.
func (l *LocalStorage) Delete(ctx context.Context, attachmentUID string) error {
	return os.RemoveAll(filepath.Join(l.directory, attachmentUID))
}

// uploadFile encrypts and uploads a file.
func (l *LocalStorage) uploadFile(ctx context.Context, destinationPath string, uid, filePath, key string) error {
	source, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	encrypted, err := l.encryption.EncryptStream(ctx, source, key)
	if err != nil {
		return err
	}

	var fileName string
	if strings.Contains(filePath, "preview") {
		fileName = fmt.Sprintf("%s.preview.enc", uid)
	} else {
		fileName = fmt.Sprintf("%s.enc", uid)
	}

	uploadPath := filepath.Join(destinationPath, fileName)
	dst, err := os.Create(uploadPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, encrypted); err != nil {
		return err
	}
	return nil
}

// Upload encrypts and stores files in the specified directory.
func (l *LocalStorage) Upload(ctx context.Context, attachment *uploader.Attachment, key string) error {
	finalPath := filepath.Join(l.directory, attachment.UID.String())
	err := getOrCreateDirectory(finalPath)
	if err != nil {
		return err
	}

	for _, filePath := range attachment.GetFilePaths() {
		err := l.uploadFile(ctx, finalPath, attachment.UID.String(), filePath, key)
		if err != nil {
			return err
		}
	}

	return nil
}

// getOrCreateDirectory creates the directory if it doesn't exist.
func getOrCreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, DIRECTORY_PERMISSIONS); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
