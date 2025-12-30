import type { AuthResponse } from '~/types/api'

const TOKEN_KEY = 'uploader_auth_token'

export const useAuth = () => {
  const token = useState<string | null>('auth_token', () => {
    if (import.meta.client) {
      return localStorage.getItem(TOKEN_KEY)
    }
    return null
  })

  const isAuthenticated = computed(() => !!token.value)

  const login = async (): Promise<void> => {
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

