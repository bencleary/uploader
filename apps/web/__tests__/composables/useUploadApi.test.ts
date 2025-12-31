import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useUploadApi } from '~/composables/useUploadApi'
import { clearNuxtState } from '#app'

// Mock $fetch
global.$fetch = vi.fn()
global.fetch = vi.fn()

describe('useUploadApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    if (typeof localStorage !== 'undefined') {
      localStorage.clear()
    }
    clearNuxtState()
  })

  it('should throw error when not authenticated', async () => {
    const api = useUploadApi()
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    
    await expect(api.uploadFile(file)).rejects.toThrow('Not authenticated')
  })

  it('should upload file when authenticated', async () => {
    const token = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const mockResponse = {
      file_name: 'test.jpg',
      preview_url: '/preview/test.jpg',
      download_url: '/download/test.jpg',
      uploaded_at: '2025-01-01T00:00:00Z'
    }
    
    ;(global.$fetch as any).mockResolvedValue(mockResponse)
    
    const api = useUploadApi()
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    const result = await api.uploadFile(file)
    
    expect(global.$fetch).toHaveBeenCalledWith(
      'http://localhost:1323/file/upload',
      expect.objectContaining({
        method: 'POST',
        headers: {
          key: token
        }
      })
    )
    expect(result).toEqual(mockResponse)
  })

  it('should handle upload errors', async () => {
    const token = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const error = new Error('Upload failed')
    ;(global.$fetch as any).mockRejectedValue(error)
    
    const api = useUploadApi()
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    
    await expect(api.uploadFile(file)).rejects.toThrow('Upload failed')
  })

  it('should get file when authenticated', async () => {
    const token = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const mockBlob = new Blob(['test'], { type: 'image/jpeg' })
    ;(global.fetch as any).mockResolvedValue({
      ok: true,
      blob: vi.fn().mockResolvedValue(mockBlob)
    })
    
    const api = useUploadApi()
    const result = await api.getFile('test-uid')
    
    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:1323/file/test-uid',
      {
        headers: {
          key: token
        }
      }
    )
    expect(result).toBe(mockBlob)
  })

  it('should get preview file when preview is true', async () => {
    const token = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const mockBlob = new Blob(['preview'], { type: 'image/jpeg' })
    ;(global.fetch as any).mockResolvedValue({
      ok: true,
      blob: vi.fn().mockResolvedValue(mockBlob)
    })
    
    const api = useUploadApi()
    await api.getFile('test-uid', true)
    
    expect(global.fetch).toHaveBeenCalledWith(
      'http://localhost:1323/file/test-uid?preview=true',
      {
        headers: {
          key: token
        }
      }
    )
  })

  it('should throw error when get file fails', async () => {
    const token = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    ;(global.fetch as any).mockResolvedValue({
      ok: false,
      statusText: 'Not Found'
    })
    
    const api = useUploadApi()
    
    await expect(api.getFile('test-uid')).rejects.toThrow('Failed to fetch file: Not Found')
  })

  it('should throw error when not authenticated for getFile', async () => {
    const api = useUploadApi()
    
    await expect(api.getFile('test-uid')).rejects.toThrow('Not authenticated')
  })
})

