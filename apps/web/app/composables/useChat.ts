import type { ChatMessage, UploadResponse } from '~/types/api'

const STORAGE_KEY = 'uploader_chat_messages'

export const useChat = () => {
  const messages = useState<ChatMessage[]>('chat_messages', () => [])
  const { get, set, remove } = useIndexedDB()
  const isLoaded = useState<boolean>('chat_messages_loaded', () => false)

  const loadMessages = async (): Promise<void> => {
    if (!import.meta.client || isLoaded.value) {
      return
    }

    try {
      const stored = await get<ChatMessage[]>(STORAGE_KEY)
      if (stored) {
        messages.value = stored.map((msg: any) => ({
          ...msg,
          timestamp: new Date(msg.timestamp)
        }))
      }
      isLoaded.value = true
    } catch (error) {
      console.error('Failed to load messages from IndexedDB:', error)
      isLoaded.value = true
    }
  }

  // Auto-load messages on client-side
  if (import.meta.client && !isLoaded.value) {
    loadMessages()
  }

  const sendMessage = async (text: string, attachments?: UploadResponse[]): Promise<void> => {
    const message: ChatMessage = {
      id: crypto.randomUUID(),
      text,
      timestamp: new Date(),
      attachments,
      isFromUser: true // User messages are from the current user
    }

    // Optimistic update
    messages.value.push(message)

    // Persist to IndexedDB
    if (import.meta.client) {
      try {
        await set(STORAGE_KEY, messages.value)
      } catch (error) {
        console.error('Failed to save message to IndexedDB:', error)
        // Rollback on error
        messages.value = messages.value.filter(m => m.id !== message.id)
        throw error
      }
    }
  }

  const getMessages = (): ChatMessage[] => {
    return messages.value
  }

  const clearMessages = async (): Promise<void> => {
    // Optimistic update
    messages.value = []

    // Remove from IndexedDB
    if (import.meta.client) {
      try {
        await remove(STORAGE_KEY)
      } catch (error) {
        console.error('Failed to clear messages from IndexedDB:', error)
        throw error
      }
    }
  }

  return {
    messages: readonly(messages),
    sendMessage,
    getMessages,
    clearMessages,
    loadMessages
  }
}

