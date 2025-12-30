# File Upload Service (Go)

A small service that handles image uploads for an imaginary chat application: it stores encrypted files, creates a resized original + a preview image, and records upload metadata for later retrieval.

Inspired by [Code Aesthetic’s](https://www.youtube.com/watch?v=J1f5b4vcxCQ) “file upload service” walkthrough (the video uses TypeScript; this project is the Go implementation).

## Features

- Upload images over HTTP
- Resize originals (max width) + generate preview images
- Encrypt stored files (AES-GCM)
- Record metadata in SQLite for later downloads

## Quickstart

### Requirements

- Go 1.21+
- A C toolchain (required by `github.com/mattn/go-sqlite3`)

### Run the server

```bash
make server
```

The server listens on `http://localhost:1323` and writes:

- Metadata: `filer.sqlite`
- Encrypted files: `temp/<upload-uuid>/...`
- Working files: `vault/<vault-uuid>/...` (temporary)

### Upload a file

The API requires an `key` header. It must be **32 characters** (AES-256 key material) and pass validation.

```bash
KEY='0123456789abcdef0123456789abcdef'
curl -sS \
  -H "key: ${KEY}" \
  -F "file=@./path/to/image.png" \
  http://localhost:1323/file/upload
```

Response shape:

```json
{
  "file_name": "example.png",
  "preview_url": "http://localhost:1323/file/<uid>?preview=true",
  "download_url": "http://localhost:1323/file/<uid>",
  "uploaded_at": "2025-01-01T00:00:00Z"
}
```

### Download an original or preview

```bash
UID='<uid-from-upload-response>'
curl -sS -H "key: ${KEY}" "http://localhost:1323/file/${UID}" -o original.bin
curl -sS -H "key: ${KEY}" "http://localhost:1323/file/${UID}?preview=true" -o preview.bin
```

## API

- `POST /file/upload` (multipart form field: `file`)
- `GET /file/:uid` (query: `preview=true|false`)

More details: `docs/API.md`.
Local S3 setup (MinIO): `docs/LOCAL_S3.md`.

## Project layout

- Root package: interfaces and core types (`uploader.*`)
- `internal/http`: Echo server + handlers
- `internal/storage`: local storage backend
- `internal/encryption`: AES-GCM encryption provider
- `internal/scaler` + `internal/preview`: image scaling + preview generation
- `internal/db`: SQLite-backed filer (metadata store)

Architecture notes: `docs/ARCHITECTURE.md`.

## Developer commands

```bash
make test
make fmt
make ci
```

## Notes / assumptions

- This service assumes authentication/authorization is handled elsewhere; the encryption key is provided per-request.
- Only common image types are supported today (`image/png`, `image/jpeg`, `image/gif`).

## Roadmap

- [ ] Finish CLI (`cmd/cli`)
- [ ] Add S3 storage backend
- [ ] Improve error responses + consistent JSON errors
- [ ] Streaming encryption (avoid buffering whole files in memory)
