import type { AuthResponse } from '~/types/api'

const TOKEN_KEY = 'uploader_auth_token'

const validateToken = (token: string | null): boolean => {
  return token !== null && token.length === 32 && /^[0-9a-f]+$/i.test(token)
}

export const useAuth = () => {
  const token = useState<string | null>('auth_token', () => {
    if (import.meta.client) {
      const stored = localStorage.getItem(TOKEN_KEY)
      // Validate token format (should be 32 hex characters)
      if (validateToken(stored)) {
        return stored
      }
      // If invalid, remove it
      if (stored) {
        localStorage.removeItem(TOKEN_KEY)
      }
    }
    return null
  })

  const isAuthenticated = computed(() => !!token.value)

  const login = async (): Promise<void> => {
    // Check if we already have a valid token in localStorage
    if (import.meta.client) {
      const existingToken = localStorage.getItem(TOKEN_KEY)
      if (validateToken(existingToken)) {
        // Reuse existing token - this ensures persistence across refreshes
        token.value = existingToken
        return
      }
    }

    // Only generate a new token if we don't have one
    try {
      const response = await $fetch<AuthResponse>('/api/auth/login', {
        method: 'POST'
      })

      if (response.token) {
        token.value = response.token
        if (import.meta.client) {
          localStorage.setItem(TOKEN_KEY, response.token)
        }
      }
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const logout = (): void => {
    token.value = null
    if (import.meta.client) {
      localStorage.removeItem(TOKEN_KEY)
    }
  }

  const getToken = (): string | null => {
    return token.value
  }

  return {
    token: readonly(token),
    isAuthenticated,
    login,
    logout,
    getToken
  }
}

