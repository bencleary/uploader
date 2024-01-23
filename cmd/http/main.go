package main

import (
	"context"

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
	storageService := storage.NewLocalStorage("temp/", "vault/", encryptionService)
	err := storageService.Initialise(context.Background())

	if err != nil {
		panic(err)
	}

	sqlite, err := db.NewSQLiteDatabase("filer.sqlite")

	if err != nil {
		panic(err)
	}

	sqlite.CreateTable()

	filingService := db.NewSqliteFilerService(sqlite)

	supportedImageScalerMimeTypes := []string{"image/png", "image/gif", "image/jpeg"}
	drawScaler := scaler.NewDrawImageScaler(supportedImageScalerMimeTypes)

	imagePreviewGenerator := preview.NewImagePreviewGenerator(drawScaler)

	previewService := uploader.NewPreviewService()

	previewService.Register("image/png", imagePreviewGenerator)
	previewService.Register("image/gif", imagePreviewGenerator)
	previewService.Register("image/jpeg", imagePreviewGenerator)

	server := http.NewServer(filingService, storageService, drawScaler, previewService)

	server.Start()

}
