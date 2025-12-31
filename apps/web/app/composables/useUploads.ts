import type { UploadMetadata } from '~/types/api'

const STORAGE_KEY = 'uploader_uploads'

export const useUploads = () => {
  const uploads = useState<UploadMetadata[]>('uploads', () => [])
  const { get, set, remove } = useIndexedDB()
  const isLoaded = useState<boolean>('uploads_loaded', () => false)

  const loadUploads = async (): Promise<void> => {
    if (!import.meta.client || isLoaded.value) {
      return
    }

    try {
      const stored = await get<UploadMetadata[]>(STORAGE_KEY)
      if (stored) {
        uploads.value = stored
      }
      isLoaded.value = true
    } catch (error) {
      console.error('Failed to load uploads from IndexedDB:', error)
      isLoaded.value = true
    }
  }

  // Auto-load uploads on client-side
  if (import.meta.client && !isLoaded.value) {
    loadUploads()
  }

  const addUpload = async (upload: UploadMetadata): Promise<void> => {
    // Optimistic update
    uploads.value.push(upload)

    // Persist to IndexedDB
    if (import.meta.client) {
      try {
        await set(STORAGE_KEY, uploads.value)
      } catch (error) {
        console.error('Failed to save upload to IndexedDB:', error)
        // Rollback on error
        const index = uploads.value.findIndex(u => u.uid === upload.uid)
        if (index !== -1) {
          uploads.value.splice(index, 1)
        }
        throw error
      }
    }
  }

  const getUploads = (): UploadMetadata[] => {
    return uploads.value
  }

  const clearUploads = async (): Promise<void> => {
    // Optimistic update
    uploads.value = []

    // Remove from IndexedDB
    if (import.meta.client) {
      try {
        await remove(STORAGE_KEY)
      } catch (error) {
        console.error('Failed to clear uploads from IndexedDB:', error)
        throw error
      }
    }
  }

  return {
    uploads: readonly(uploads),
    addUpload,
    getUploads,
    clearUploads,
    loadUploads
  }
}

