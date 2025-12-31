# Uploader Demo UI (Nuxt)

A Nuxt UI “chat app” demo that exercises the Go upload service in the repo root.

## Requirements

- Node.js 20+
- A running backend API (default: `http://localhost:1323`)

## Configuration

- `NUXT_PUBLIC_API_BASE_URL` (optional): backend base URL (default `http://localhost:1323`)

Example:

```bash
NUXT_PUBLIC_API_BASE_URL=http://localhost:1323
```

## Run locally

From `apps/web`:

```bash
npm ci
npm run dev
```

Open `http://localhost:3000` and click **Start Demo Session**. This generates a demo encryption key client-side and stores it in `localStorage`. Uploading an image sends that key as the `key` header to the Go API.

## Scripts

```bash
npm test
```

## Notes

- This frontend is intentionally a demo environment. The `/api/auth/login` route just generates a random 32‑char key for local use.
- File downloads/images go through a Nuxt server route (`/api/file/:uid`) so the browser can load images while still passing the key to the backend.

Optional (extra checks):

```bash
npm run typecheck
npm run lint
npm run lint -- --fix
```
