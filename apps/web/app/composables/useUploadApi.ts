import type { UploadResponse } from '~/types/api'

const TOKEN_KEY = 'uploader_auth_token'

export const useUploadApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl || 'http://localhost:1323'

  const getAuthToken = (): string | null => {
    if (import.meta.client) {
      return localStorage.getItem(TOKEN_KEY)
    }
    return null
  }

  const uploadFile = async (file: File): Promise<UploadResponse> => {
    const token = getAuthToken()
    if (!token) {
      throw new Error('Not authenticated. Please login first.')
    }

    const formData = new FormData()
    formData.append('file', file)

    try {
      const response = await $fetch<UploadResponse>(`${baseURL}/file/upload`, {
        method: 'POST',
        headers: {
          key: token
        },
        body: formData
      })

      return response
    } catch (error) {
      console.error('Upload failed:', error)
      throw error
    }
  }

  const getFile = async (uid: string, preview = false): Promise<Blob> => {
    const token = getAuthToken()
    if (!token) {
      throw new Error('Not authenticated. Please login first.')
    }

    const url = `${baseURL}/file/${uid}${preview ? '?preview=true' : ''}`

    try {
      const response = await fetch(url, {
        headers: {
          key: token
        }
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch file: ${response.statusText}`)
      }

      return await response.blob()
    } catch (error) {
      console.error('Get file failed:', error)
      throw error
    }
  }

  return {
    uploadFile,
    getFile
  }
}

