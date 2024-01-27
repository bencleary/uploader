package db

import (
	"testing"

	"github.com/bencleary/uploader"
	"github.com/google/uuid"
)

func TestNewSqliteFiler(t *testing.T) {

	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	filer := NewSqliteFilerService(db)

	if filer == nil {
		t.Fatal("expected filer")
	}

}

func TestSqliteFilerRecord(t *testing.T) {
	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	filer := NewSqliteFilerService(db)

	if filer == nil {
		t.Fatal("expected filer")
	}

	attachment := &uploader.Attachment{
		UID:       uuid.New(),
		OwnerID:   1,
		FileName:  "test",
		FileSize:  0,
		Extension: "",
		MimeType:  "",
	}

	err = filer.Record(attachment)
	if err != nil {
		t.Fatal(err)
	}

}

func TestFilerFetch(t *testing.T) {
	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	filer := NewSqliteFilerService(db)

	if filer == nil {
		t.Fatal("expected filer")
	}

	attachment := &uploader.Attachment{
		UID:       uuid.New(),
		OwnerID:   1,
		FileName:  "test",
		FileSize:  0,
		Extension: "",
		MimeType:  "",
	}

	err = filer.Record(attachment)
	if err != nil {
		t.Fatal(err)
	}

	row, err := filer.Fetch(attachment.UID)
	if err != nil {
		t.Fatal(err)
	}

	if row == nil {
		t.Fatal("expected row")
	}

	// if row.OwnerID == attachment.OwnerID {
	// 	t.Fatal("owner id does not match")
	// }

}

func TestFilerDelete(t *testing.T) {
	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTable()
	if err != nil {
		t.Fatal(err)
	}

	filer := NewSqliteFilerService(db)

	if filer == nil {
		t.Fatal("expected filer")
	}

	attachment := &uploader.Attachment{
		UID:       uuid.New(),
		OwnerID:   1,
		FileName:  "test",
		FileSize:  0,
		Extension: "",
		MimeType:  "",
	}

	err = filer.Record(attachment)
	if err != nil {
		t.Fatal(err)
	}

	err = filer.Delete(attachment.UID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = filer.Fetch(attachment.UID)
	if err == nil {
		t.Fatal("expected not to find device")
	}

}
