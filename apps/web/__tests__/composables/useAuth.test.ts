import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useAuth } from '~/composables/useAuth'
import { clearNuxtState } from '#app'

// Mock $fetch
global.$fetch = vi.fn()

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    if (typeof localStorage !== 'undefined') {
      localStorage.clear()
    }
    clearNuxtState()
  })

  it('should initialize with no token when localStorage is empty', () => {
    const auth = useAuth()
    expect(auth.token.value).toBeNull()
    expect(auth.isAuthenticated.value).toBe(false)
  })

  it('should validate and use existing token from localStorage', () => {
    const validToken = 'a'.repeat(32)
    localStorage.setItem('uploader_auth_token', validToken)
    
    const auth = useAuth()
    expect(auth.token.value).toBe(validToken)
    expect(auth.isAuthenticated.value).toBe(true)
  })

  it('should reject invalid token format from localStorage', () => {
    const invalidToken = 'invalid-token'
    localStorage.setItem('uploader_auth_token', invalidToken)
    
    const auth = useAuth()
    expect(auth.token.value).toBeNull()
    expect(auth.isAuthenticated.value).toBe(false)
    expect(localStorage.getItem('uploader_auth_token')).toBeNull()
  })

  it('should reject token that is not 32 characters', () => {
    const shortToken = 'a'.repeat(31)
    localStorage.setItem('uploader_auth_token', shortToken)
    
    const auth = useAuth()
    expect(auth.token.value).toBeNull()
  })

  it('should reject token with non-hex characters', () => {
    const invalidToken = 'g'.repeat(32)
    localStorage.setItem('uploader_auth_token', invalidToken)
    
    const auth = useAuth()
    expect(auth.token.value).toBeNull()
  })

  it('should login and store token', async () => {
    const mockToken = 'a'.repeat(32)
    const mockResponse = { token: mockToken }
    
    ;(global.$fetch as any).mockResolvedValue(mockResponse)
    
    const auth = useAuth()
    await auth.login()
    
    expect(global.$fetch).toHaveBeenCalledWith('/api/auth/login', {
      method: 'POST'
    })
    expect(auth.token.value).toBe(mockToken)
    expect(auth.isAuthenticated.value).toBe(true)
    expect(localStorage.getItem('uploader_auth_token')).toBe(mockToken)
  })

  it('should reuse existing valid token on login', async () => {
    const existingToken = 'b'.repeat(32)
    localStorage.setItem('uploader_auth_token', existingToken)
    
    const auth = useAuth()
    await auth.login()
    
    expect(global.$fetch).not.toHaveBeenCalled()
    expect(auth.token.value).toBe(existingToken)
  })

  it('should handle login errors', async () => {
    const error = new Error('Login failed')
    ;(global.$fetch as any).mockRejectedValue(error)
    
    const auth = useAuth()
    
    await expect(auth.login()).rejects.toThrow('Login failed')
  })

  it('should logout and clear token', () => {
    const token = 'c'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const auth = useAuth()
    auth.logout()
    
    expect(auth.token.value).toBeNull()
    expect(auth.isAuthenticated.value).toBe(false)
    expect(localStorage.getItem('uploader_auth_token')).toBeNull()
  })

  it('should get token', () => {
    const token = 'd'.repeat(32)
    localStorage.setItem('uploader_auth_token', token)
    
    const auth = useAuth()
    expect(auth.getToken()).toBe(token)
  })
})

