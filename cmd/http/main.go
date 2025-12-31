package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	// Load storage configuration from environment variables
	storageType := getEnv("UPLOADER_STORAGE", "local")
	var storageService uploader.StorageService

	if storageType == "s3" {
		s3Options := &storage.S3Options{
			Endpoint:       getEnv("UPLOADER_S3_ENDPOINT", ""),
			Bucket:         getEnv("UPLOADER_S3_BUCKET", "uploader"),
			Region:         getEnv("UPLOADER_S3_REGION", "us-east-1"),
			Prefix:         getEnv("UPLOADER_S3_PREFIX", ""),
			ForcePathStyle: getEnvBool("UPLOADER_S3_FORCE_PATH_STYLE", true),
			AccessKeyID:    getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretKey:      getEnv("AWS_SECRET_ACCESS_KEY", ""),
		}

		s3Storage := storage.NewS3Storage(s3Options, encryptionService)
		if s3Storage == nil {
			panic("failed to initialize S3 storage: invalid configuration")
		}
		storageService = s3Storage
	} else {
		// Default to local storage
		uploadPath := getEnv("UPLOADER_LOCAL_UPLOAD_PATH", "temp/")
		vaultPath := getEnv("UPLOADER_LOCAL_VAULT_PATH", "vault/")
		storageService = storage.NewLocalStorage(uploadPath, vaultPath, encryptionService)
	}

	err := storageService.Initialise(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to initialize storage: %v", err))
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

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool retrieves an environment variable as a boolean or returns a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		value = strings.ToLower(strings.TrimSpace(value))
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
