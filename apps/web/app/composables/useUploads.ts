import type { UploadMetadata } from '~/types/api'

const STORAGE_KEY = 'uploader_uploads'

export const useUploads = () => {
  const uploads = useState<UploadMetadata[]>('uploads', () => {
    if (import.meta.client) {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        try {
          return JSON.parse(stored)
        } catch {
          return []
        }
      }
    }
    return []
  })

  const addUpload = (upload: UploadMetadata): void => {
    uploads.value.push(upload)

    if (import.meta.client) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(uploads.value))
    }
  }

  const getUploads = (): UploadMetadata[] => {
    return uploads.value
  }

  const clearUploads = (): void => {
    uploads.value = []
    if (import.meta.client) {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  return {
    uploads: readonly(uploads),
    addUpload,
    getUploads,
    clearUploads
  }
}

