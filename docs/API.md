# API

Base URL: `http://localhost:1323`

## Authentication / encryption key

Requests must include a `key` header. The key is used to encrypt/decrypt stored files.

- Header: `key: <32 characters>`
- Validation: see `internal/encryption/aes.go` (`encryption.IsValidKey`)

## `POST /file/upload`

Uploads an image, creates a preview image, encrypts both files, and records upload metadata.

### Request

- Content-Type: `multipart/form-data`
- Form field: `file` (required)
- Header: `key` (required)

Example:

```bash
KEY='0123456789abcdef0123456789abcdef'
curl -sS \
  -H "key: ${KEY}" \
  -F "file=@./path/to/image.png" \
  http://localhost:1323/file/upload
```

### Response (200)

```json
{
  "file_name": "image.png",
  "preview_url": "http://localhost:1323/file/<uid>?preview=true",
  "download_url": "http://localhost:1323/file/<uid>",
  "uploaded_at": "2025-01-01T00:00:00Z"
}
```

## `GET /file/:uid`

Downloads and decrypts a previously uploaded file.

### Request

- Path param: `uid` (required, UUID)
- Header: `key` (required)
- Query param: `preview` (optional, `true|false`, defaults to `false`)

Examples:

```bash
KEY='0123456789abcdef0123456789abcdef'
UID='<uuid>'

curl -sS -H "key: ${KEY}" "http://localhost:1323/file/${UID}" -o original.bin
curl -sS -H "key: ${KEY}" "http://localhost:1323/file/${UID}?preview=true" -o preview.bin
```

## Error behavior

Errors are currently a mix of Echo HTTP errors and internal typed errors. A cleanup to return consistent JSON error bodies is on the roadmap (see `README.md`).

