import { describe, it, expect, beforeEach, vi } from 'vitest'
import { useChat } from '~/composables/useChat'
import type { ChatMessage, UploadResponse } from '~/types/api'
import { clearNuxtState } from '#app'
import { useIndexedDB } from '~/composables/useIndexedDB'

// Mock crypto.randomUUID to return consistent IDs for testing
const mockUUID = vi.fn(() => 'test-uuid-123')
Object.defineProperty(global, 'crypto', {
  value: {
    randomUUID: mockUUID
  },
  writable: true
})

describe('useChat', () => {
  beforeEach(async () => {
    vi.clearAllMocks()
    clearNuxtState()
    mockUUID.mockReturnValue(`test-uuid-${Date.now()}-${Math.random()}`)
    
    // Clear IndexedDB
    const { remove } = useIndexedDB()
    await remove('uploader_chat_messages').catch(() => {})
  })

  it('should initialize with empty messages', async () => {
    const chat = useChat()
    await chat.loadMessages()
    expect(chat.messages.value).toEqual([])
    expect(chat.getMessages()).toEqual([])
  })

  it('should load messages from IndexedDB', async () => {
    const storedMessages: ChatMessage[] = [
      {
        id: '1',
        text: 'Hello',
        timestamp: new Date('2025-01-01'),
        attachments: []
      }
    ]
    const { set } = useIndexedDB()
    await set('uploader_chat_messages', storedMessages)
    
    const chat = useChat()
    await chat.loadMessages()
    expect(chat.messages.value).toHaveLength(1)
    expect(chat.messages.value[0].text).toBe('Hello')
  })

  it('should send a message', async () => {
    const chat = useChat()
    await chat.sendMessage('Test message')
    
    expect(chat.messages.value).toHaveLength(1)
    expect(chat.messages.value[0].text).toBe('Test message')
    expect(chat.messages.value[0].id).toBeTruthy()
    expect(typeof chat.messages.value[0].id).toBe('string')
    expect(chat.messages.value[0].timestamp).toBeInstanceOf(Date)
  })

  it('should send a message with attachments', async () => {
    const chat = useChat()
    const attachments: UploadResponse[] = [
      {
        file_name: 'test.jpg',
        preview_url: '/preview/test.jpg',
        download_url: '/download/test.jpg',
        uploaded_at: '2025-01-01T00:00:00Z'
      }
    ]
    
    await chat.sendMessage('Message with attachment', attachments)
    
    expect(chat.messages.value).toHaveLength(1)
    expect(chat.messages.value[0].attachments).toEqual(attachments)
  })

  it('should persist messages to IndexedDB', async () => {
    const chat = useChat()
    await chat.sendMessage('Persisted message')
    
    const { get } = useIndexedDB()
    const stored = await get<ChatMessage[]>('uploader_chat_messages')
    expect(stored).toBeTruthy()
    expect(stored).toHaveLength(1)
    expect(stored![0].text).toBe('Persisted message')
  })

  it('should clear messages', async () => {
    const chat = useChat()
    await chat.sendMessage('Message 1')
    await chat.sendMessage('Message 2')
    
    expect(chat.messages.value).toHaveLength(2)
    
    await chat.clearMessages()
    
    expect(chat.messages.value).toEqual([])
    const { get } = useIndexedDB()
    const stored = await get('uploader_chat_messages')
    expect(stored).toBeNull()
  })

  it('should get all messages', async () => {
    const chat = useChat()
    await chat.sendMessage('Message 1')
    await chat.sendMessage('Message 2')
    
    const messages = chat.getMessages()
    expect(messages).toHaveLength(2)
    expect(messages[0].text).toBe('Message 1')
    expect(messages[1].text).toBe('Message 2')
  })
})

