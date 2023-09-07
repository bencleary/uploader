package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/bencleary/uploader"

	"github.com/google/uuid"
)

var _ uploader.StorageService = (*LocalStorage)(nil)

const (
	DIRECTORY_PERMISSIONS = 0755
)

type LocalStorage struct {
	directory  string
	vault      string
	encryption uploader.EncryptionService
}

// NewLocalStorageStrategy creates a new instance of the LocalStorage struct with the provided directory and encryption service.
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

func (l *LocalStorage) newFileName(uid uuid.UUID) string {
	return fmt.Sprintf("%s.enc", uid)
}

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
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (l *LocalStorage) Hold(ctx context.Context, file *multipart.FileHeader) (*uploader.Attachment, error) {
	vaultID, err := uuid.NewUUID()
	if err != nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	vaultDir := filepath.Join("vault", vaultID.String())

	_, err = os.Stat(vaultDir)

	if os.IsNotExist(err) {
		err = os.Mkdir(vaultDir, DIRECTORY_PERMISSIONS)
		if err != nil {
			return nil, err
		}
		fmt.Println("Directory created:", vaultDir)
	} else {
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

	// TODO extract directory creation to function

	attachment := uploader.NewAttachment(file, 1)

	attachment.LocalPath = vaultPath

	return attachment, nil
}

func (l *LocalStorage) Load(ctx context.Context, fileUID uuid.UUID) (*uploader.Attachment, error) {
	return nil, nil
}

func (l *LocalStorage) Delete(ctx context.Context, fileUID uuid.UUID) error {
	return nil
}

// Save - Copies files from the vault path, encrypts them and then stores them as per the implementation
func (l *LocalStorage) Save(ctx context.Context, attachment *uploader.Attachment, key string) error {
	finalPath := filepath.Join(l.directory, attachment.UID.String())

	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		if err := os.Mkdir(finalPath, DIRECTORY_PERMISSIONS); err != nil {
			return err
		}
	} else if err != nil {
		fmt.Println(err)
		return err
	}

	source, err := os.Open(attachment.LocalPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer source.Close()

	encrypted, err := l.encryption.EncryptStream(ctx, source, key)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fileName := l.newFileName(attachment.UID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	uploadPath := filepath.Join(finalPath, fileName)
	dst, err := os.Create(uploadPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, encrypted); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
