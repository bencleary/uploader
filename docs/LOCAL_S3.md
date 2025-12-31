# Local S3 (MinIO)

This repo includes a `docker-compose.yml` that runs MinIO (an S3-compatible server) for local development and testing.

## Start MinIO + create bucket

1. Create a `.env` file:

```bash
cp .env.example .env
```

2. Start services:

```bash
make s3-up
```

3. Confirm MinIO is running:

- S3 API: `http://localhost:9000`
- Console UI: `http://localhost:9001`

Login uses `MINIO_ROOT_USER` / `MINIO_ROOT_PASSWORD` from `.env`.

## App configuration (for S3 backend)

The S3 storage backend is fully implemented. To use it, set the following environment variables:

**Required:**
- `UPLOADER_STORAGE=s3` - Enable S3 storage backend
- `UPLOADER_S3_ENDPOINT=http://localhost:9000` - S3 endpoint URL
- `UPLOADER_S3_BUCKET=uploader` - S3 bucket name
- `UPLOADER_S3_REGION=us-east-1` - AWS region
- `AWS_ACCESS_KEY_ID=minioadmin` - Access key (MinIO default)
- `AWS_SECRET_ACCESS_KEY=minioadmin` - Secret key (MinIO default)

**Optional:**
- `UPLOADER_S3_PREFIX=uploader` - Prefix for object keys (default: empty)
- `UPLOADER_S3_FORCE_PATH_STYLE=true` - Use path-style addressing (recommended for MinIO, default: true)

**For local storage (default):**
- `UPLOADER_STORAGE=local` or omit the variable
- `UPLOADER_LOCAL_UPLOAD_PATH=temp/` - Upload directory (default: `temp/`)
- `UPLOADER_LOCAL_VAULT_PATH=vault/` - Vault directory (default: `vault/`)

### Example .env file for S3

```bash
UPLOADER_STORAGE=s3
UPLOADER_S3_ENDPOINT=http://localhost:9000
UPLOADER_S3_BUCKET=uploader
UPLOADER_S3_REGION=us-east-1
UPLOADER_S3_PREFIX=uploader
UPLOADER_S3_FORCE_PATH_STYLE=true
AWS_ACCESS_KEY_ID=minioadmin
AWS_SECRET_ACCESS_KEY=minioadmin
```

## Stop MinIO

```bash
make s3-down
```

## Testing with S3

The S3 storage implementation includes comprehensive tests. To run tests that require S3/MinIO:

1. Start MinIO using `make s3-up`
2. Ensure the bucket exists (created automatically by `minio-init` service)
3. Run tests: `go test ./internal/storage/... -v`

Tests will automatically skip S3 integration tests if MinIO is not available, so unit tests will still pass without a running S3 instance.
