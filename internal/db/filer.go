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

func (s *SqliteFiler) Fetch(fileUID uuid.UUID) (*uploader.Attachment, error) {
	row := s.db.db.QueryRow(`
		SELECT owner_id, file_name, file_size, extension, mime_type
		FROM uploads
		WHERE uuid = ?
	`, fileUID.String())

	attachment := &uploader.Attachment{UID: fileUID}

	err := row.Scan(&attachment.OwnerID, &attachment.FileName, &attachment.FileSize, &attachment.Extension, &attachment.MimeType)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *SqliteFiler) Delete(fileUID uuid.UUID) error {
	_, err := s.db.db.Exec(`
		DELETE FROM uploads
		WHERE uuid = ?
	`, fileUID.String())
	if err != nil {
		return err
	}
	return nil
}
