package db

import (
	"github.com/bencleary/uploader"
	"github.com/google/uuid"
)

var _ uploader.FilerService = (*SqliteFiler)(nil)

type SqliteFiler struct {
	db *DB
}

func NewSqliteFilerService(db *DB) *SqliteFiler {
	return &SqliteFiler{
		db: db,
	}
}

func (s *SqliteFiler) Record(attachment *uploader.Attachment) error {
	_, err := s.db.db.Exec(`
		INSERT INTO uploads (uuid, owner_id, file_name, file_size, extension, mime_type)
		VALUES (?, ?, ?, ?, ?, ?)
	`, attachment.UID.String(), attachment.OwnerID, attachment.FileName, attachment.FileSize, attachment.Extension, attachment.MimeType)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqliteFiler) Fetch(fileUID uuid.UUID) (*uploader.Upload, error) {
	return nil, nil
}

func (s *SqliteFiler) Delete(fileUID uuid.UUID) error {
	return nil
}
