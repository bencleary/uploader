package db

import (
	"testing"
)

func TestNewSQLiteDatabase(t *testing.T) {
	_, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSQLiteDatabaseClose(t *testing.T) {
	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSQLiteDatabaseCreateTable(t *testing.T) {
	db, err := NewSQLiteDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTable()
	if err != nil {
		t.Fatal(err)
	}
}
