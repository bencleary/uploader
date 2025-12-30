import type { ChatMessage, UploadResponse } from '~/types/api'

const STORAGE_KEY = 'uploader_chat_messages'

export const useChat = () => {
  const messages = useState<ChatMessage[]>('chat_messages', () => {
    if (import.meta.client) {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) {
        try {
          const parsed = JSON.parse(stored)
          return parsed.map((msg: any) => ({
            ...msg,
            timestamp: new Date(msg.timestamp)
          }))
        } catch {
          return []
        }
      }
    }
    return []
  })

  const sendMessage = (text: string, attachments?: UploadResponse[]): void => {
    const message: ChatMessage = {
      id: crypto.randomUUID(),
      text,
      timestamp: new Date(),
      attachments
    }

    messages.value.push(message)

    if (import.meta.client) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(messages.value))
    }
  }

  const getMessages = (): ChatMessage[] => {
    return messages.value
  }

  const clearMessages = (): void => {
    messages.value = []
    if (import.meta.client) {
      localStorage.removeItem(STORAGE_KEY)
    }
  }

  return {
    messages: readonly(messages),
    sendMessage,
    getMessages,
    clearMessages
  }
}

