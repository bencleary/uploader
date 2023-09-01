package main

import (
	"context"
	"fmt"

	"github.com/bencleary/uploader/internal/db"
	"github.com/bencleary/uploader/internal/encryption"
	"github.com/bencleary/uploader/internal/keystore"
	"github.com/bencleary/uploader/internal/storage"
)

func main() {
	keyService := keystore.NewInMemoryKeyStore()
	encryptionService := encryption.NewAESService(keyService)
	storageService := storage.NewLocalStorage("", "", encryptionService)
	err := storageService.Initialise(context.Background())

	if err != nil {
		panic(err)
	}

	sqlite, err := db.NewSQLiteDatabase("")

	if err != nil {
		panic(err)
	}

	filingService := db.NewSqliteFilerService(sqlite)
	fmt.Println(filingService)
}
