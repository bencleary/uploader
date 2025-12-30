# Architecture

This repo is structured to keep interfaces and core types in the root `uploader` package, with concrete implementations living under `internal/`.

## Key flows

### Upload (`POST /file/upload`)

1. Validate `key` header (`internal/middleware/validators.go`).
2. Read multipart file header from the request (`echo.Context.FormFile`).
3. `StorageService.Hold`: persist the raw upload to a working area (`internal/storage/local.go`).
4. Create a preview copy next to the working file (`Attachment.CopyFileToPath`).
5. `ScalerService.Scale`: resize the original to a max width, and the preview to a smaller width.
6. `StorageService.Upload`: encrypt and store files under `temp/<uid>/...`.
7. `FilerService.Record`: store upload metadata in SQLite (`internal/db`).

### Download (`GET /file/:uid`)

1. `FilerService.Fetch`: retrieve metadata for the UID.
2. `StorageService.Download`: open the encrypted blob (original or preview) and decrypt it.
3. Stream decrypted bytes back to the client.

## Interfaces

- `uploader.StorageService`: storage backend (local today; can be extended to S3).
- `uploader.EncryptionService`: encrypt/decrypt files.
- `uploader.FilerService`: metadata store (SQLite today).
- `uploader.ScalerService` + `uploader.PreviewGeneratorService`: image processing pipeline.

## Implementation notes

- Local storage writes working files to `vault/` and encrypted output to `temp/` by default (see `cmd/http/main.go`).
- The current encryption implementation buffers files in memory before encrypting/decrypting; converting this to true streaming is a good next enhancement.

