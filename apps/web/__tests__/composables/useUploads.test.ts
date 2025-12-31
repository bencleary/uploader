import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useUploads } from '~/composables/useUploads'
import type { UploadMetadata } from '~/types/api'
import { clearNuxtState } from '#app'
import { useIndexedDB } from '~/composables/useIndexedDB'

describe('useUploads', () => {
  beforeEach(async () => {
    vi.clearAllMocks()
    clearNuxtState()
    
    // Clear IndexedDB
    const { remove } = useIndexedDB()
    await remove('uploader_uploads').catch(() => {})
  })

  it('should initialize with empty uploads', async () => {
    const uploads = useUploads()
    await uploads.loadUploads()
    expect(uploads.uploads.value).toEqual([])
    expect(uploads.getUploads()).toEqual([])
  })

  it('should load uploads from IndexedDB', async () => {
    const storedUploads: UploadMetadata[] = [
      {
        uid: 'test-uid-1',
        fileName: 'test1.jpg',
        previewUrl: '/preview/test1.jpg',
        downloadUrl: '/download/test1.jpg',
        uploadedAt: '2025-01-01T00:00:00Z'
      }
    ]
    const { set } = useIndexedDB()
    await set('uploader_uploads', storedUploads)
    
    const uploads = useUploads()
    await uploads.loadUploads()
    expect(uploads.uploads.value).toHaveLength(1)
    expect(uploads.uploads.value[0].uid).toBe('test-uid-1')
  })

  it('should add upload', async () => {
    const uploads = useUploads()
    const upload: UploadMetadata = {
      uid: 'test-uid',
      fileName: 'test.jpg',
      previewUrl: '/preview/test.jpg',
      downloadUrl: '/download/test.jpg',
      uploadedAt: '2025-01-01T00:00:00Z'
    }
    
    await uploads.addUpload(upload)
    
    expect(uploads.uploads.value).toHaveLength(1)
    expect(uploads.uploads.value[0]).toEqual(upload)
  })

  it('should persist uploads to IndexedDB', async () => {
    const uploads = useUploads()
    const upload: UploadMetadata = {
      uid: 'test-uid',
      fileName: 'test.jpg',
      previewUrl: '/preview/test.jpg',
      downloadUrl: '/download/test.jpg',
      uploadedAt: '2025-01-01T00:00:00Z'
    }
    
    await uploads.addUpload(upload)
    
    const { get } = useIndexedDB()
    const stored = await get<UploadMetadata[]>('uploader_uploads')
    expect(stored).toBeTruthy()
    expect(stored).toHaveLength(1)
    expect(stored![0].uid).toBe('test-uid')
  })

  it('should clear uploads', async () => {
    const uploads = useUploads()
    const upload: UploadMetadata = {
      uid: 'test-uid',
      fileName: 'test.jpg',
      previewUrl: '/preview/test.jpg',
      downloadUrl: '/download/test.jpg',
      uploadedAt: '2025-01-01T00:00:00Z'
    }
    
    await uploads.addUpload(upload)
    expect(uploads.uploads.value).toHaveLength(1)
    
    await uploads.clearUploads()
    
    expect(uploads.uploads.value).toEqual([])
    const { get } = useIndexedDB()
    const stored = await get('uploader_uploads')
    expect(stored).toBeNull()
  })

  it('should get all uploads', async () => {
    const uploads = useUploads()
    const upload1: UploadMetadata = {
      uid: 'test-uid-1',
      fileName: 'test1.jpg',
      previewUrl: '/preview/test1.jpg',
      downloadUrl: '/download/test1.jpg',
      uploadedAt: '2025-01-01T00:00:00Z'
    }
    const upload2: UploadMetadata = {
      uid: 'test-uid-2',
      fileName: 'test2.jpg',
      previewUrl: '/preview/test2.jpg',
      downloadUrl: '/download/test2.jpg',
      uploadedAt: '2025-01-02T00:00:00Z'
    }
    
    await uploads.addUpload(upload1)
    await uploads.addUpload(upload2)
    
    const allUploads = uploads.getUploads()
    expect(allUploads).toHaveLength(2)
    expect(allUploads[0].uid).toBe('test-uid-1')
    expect(allUploads[1].uid).toBe('test-uid-2')
  })
})

