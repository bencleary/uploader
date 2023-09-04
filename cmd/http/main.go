package main

import (
	"context"
	"fmt"

	"github.com/bencleary/uploader"
	"github.com/bencleary/uploader/internal/db"
	"github.com/bencleary/uploader/internal/encryption"
	"github.com/bencleary/uploader/internal/http"
	"github.com/bencleary/uploader/internal/keystore"
	"github.com/bencleary/uploader/internal/preview"
	"github.com/bencleary/uploader/internal/scaler"
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

	supportedMimeTypes := []string{"image/png", "image/gif", "image/jpeg"}
	drawScaler := scaler.NewDrawImageScaler(supportedMimeTypes)

	imagePreviewGenerator := preview.NewImagePreviewGenerator(drawScaler)

	previewService := uploader.NewPreviewService()

	previewService.Register("image/png", imagePreviewGenerator)
	previewService.Register("image/gif", imagePreviewGenerator)
	previewService.Register("image/jpeg", imagePreviewGenerator)

	server := http.NewServer(filingService, storageService, previewService)

	server.Open()

}