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

Recommended env vars for your upcoming `S3Storage` implementation:

- `AWS_ACCESS_KEY_ID=minioadmin`
- `AWS_SECRET_ACCESS_KEY=minioadmin`
- `AWS_REGION=us-east-1`
- `UPLOADER_STORAGE=s3`
- `UPLOADER_S3_ENDPOINT=http://localhost:9000`
- `UPLOADER_S3_BUCKET=uploader`
- `UPLOADER_S3_REGION=us-east-1`
- `UPLOADER_S3_PREFIX=uploader` (optional)
- `UPLOADER_S3_FORCE_PATH_STYLE=true` (recommended for MinIO)

## Stop MinIO

```bash
make s3-down
```
