import type { AuthResponse } from '~/types/api'

function generateEncryptionKey(): string {
  // Generate a random 32-character hex string (16 bytes = 32 hex chars)
  const bytes = new Uint8Array(16)
  crypto.getRandomValues(bytes)
  return Array.from(bytes)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('')
}

export default defineEventHandler(async (event): Promise<AuthResponse> => {
  const key = generateEncryptionKey()

  return {
    token: key,
    message: 'Demo authentication - this is a development environment only!'
  }
})

