<template>
  <div class="flex flex-col h-full w-full">
    <!-- Demo Warning Banner -->
    <UAlert
      v-if="isAuthenticated"
      color="warning"
      variant="soft"
      icon="i-lucide-alert-triangle"
      title="Demo Environment"
      description="This is a demo authentication service. The encryption key is generated client-side and should not be used in production."
      :close-button="{ icon: 'i-lucide-x', color: 'amber', variant: 'link', 'aria-label': 'Close' }"
    />

    <!-- Login Section -->
    <div
      v-if="!isAuthenticated"
      class="flex-1 w-full flex items-center justify-center p-6 bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900"
    >
      <UCard class="w-full max-w-md bg-zinc-900 dark:bg-zinc-950 shadow-xl border-zinc-950 dark:border-zinc-900">
        <template #header>
          <div class="text-center">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-emerald-500/10 mb-4">
              <UIcon name="i-lucide-upload" class="w-8 h-8 text-emerald-600 dark:text-emerald-400" />
            </div>
            <h2 class="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mb-2">Welcome</h2>
            <p class="text-zinc-600 dark:text-zinc-400">Get started with image uploads</p>
          </div>
        </template>

        <div class="space-y-6 py-2">
          <div class="bg-amber-50 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-900 rounded-lg p-4">
            <p class="text-sm text-amber-900 dark:text-amber-200">
              <UIcon name="i-lucide-info" class="w-4 h-4 inline mr-1" />
              This is a demo environment. Click below to generate a demo authentication token.
            </p>
          </div>

          <UButton
            @click="handleLogin"
            :loading="isLoggingIn"
            block
            size="lg"
            color="primary"
            class="font-semibold"
          >
            <UIcon name="i-lucide-sparkles" class="w-5 h-5 mr-2" />
            Start Demo Session
          </UButton>
        </div>
      </UCard>
    </div>

    <!-- Chat Interface -->
    <div
      v-else
      class="flex-1 flex flex-col h-full w-full overflow-hidden"
    >
      <!-- Header -->
      <div class="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-4 shrink-0">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center">
              <UIcon name="i-lucide-message-circle" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" />
            </div>
            <div>
              <h1 class="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mb-2">Chat</h1>
              <p class="text-sm text-zinc-500 dark:text-zinc-400">
                Share images and messages
              </p>
            </div>
          </div>
          <div class="flex gap-2">
            <UButton
              to="/uploads"
              icon="i-lucide-images"
              variant="ghost"
              color="zinc"
              size="lg"
            >
              Gallery
            </UButton>
            <UButton
              @click="handleLogout"
              icon="i-lucide-log-out"
              variant="ghost"
              color="zinc"
              size="lg"
            >
              Logout
            </UButton>
          </div>
        </div>
      </div>

      <!-- Messages -->
      <div
        ref="messagesContainer"
        class="flex-1 overflow-y-auto px-6 py-6 space-y-6 bg-zinc-50 dark:bg-zinc-950"
      >
        <div
          v-if="messages.length === 0"
          class="flex flex-col items-center justify-center h-full text-center"
        >
          <div class="w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center mb-4">
            <UIcon name="i-lucide-message-square-plus" class="w-10 h-10 text-emerald-600 dark:text-emerald-400" />
          </div>
          <h3 class="text-lg font-semibold text-zinc-900 dark:text-zinc-50 mb-2">No messages yet</h3>
          <p class="text-sm text-zinc-500 dark:text-zinc-400 max-w-sm">
            Start a conversation by sending a message or uploading an image
          </p>
        </div>

        <ChatMessage
          v-for="message in Array.from(messages)"
          :key="message.id"
          :message="{ ...message, attachments: message.attachments ? [...message.attachments] : undefined }"
          :get-preview-url="getPreviewUrl"
          @image-click="handleImageClick"
        />
      </div>

      <!-- Input Area -->
      <ChatInput
        v-model:message-text="messageText"
        :selected-file="selectedFile"
        :is-uploading="isUploading"
        @file-selected="handleFileSelect"
        @clear-file="selectedFile = null"
        @submit="handleSendMessage"
      />
    </div>

    <!-- Image Modal for chat images -->
    <ImageModal
      v-model="isImageModalOpen"
      :image="selectedImage"
      :on-download="handleImageDownload"
    />
  </div>
</template>

<script setup lang="ts">
import type { UploadMetadata, UploadResponse } from '~/types/api'

const { isAuthenticated, login, logout, getToken } = useAuth()
const { sendMessage, messages } = useChat()
const { uploadFile } = useUploadApi()
const { addUpload, getUploads } = useUploads()

const messageText = ref('')
const selectedFile = ref<File | null>(null)
const isLoggingIn = ref(false)
const isUploading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const isImageModalOpen = ref(false)
const selectedImage = ref<UploadMetadata | null>(null)

const handleLogin = async () => {
  isLoggingIn.value = true
  try {
    await login()
  } catch (error) {
    console.error('Login failed:', error)
  } finally {
    isLoggingIn.value = false
  }
}

const handleLogout = () => {
  logout()
  messageText.value = ''
  selectedFile.value = null
}

const handleFileSelect = (file: File) => {
  selectedFile.value = file
}

const handleImageClick = async (attachment: UploadResponse) => {
  // Extract UID from attachment URL
  const uid = extractUidFromUrl(attachment.preview_url || attachment.download_url)
  if (!uid) return

  // Find the upload in our uploads list or create metadata from attachment
  const { getUploads } = useUploads()
  const uploads = getUploads()
  let upload = uploads.find(u => u.uid === uid)

  // If not found in uploads, create metadata from attachment
  if (!upload) {
    const token = getToken()
    upload = {
      uid,
      fileName: attachment.file_name,
      previewUrl: `/api/file/${uid}?preview=true${token ? `&key=${encodeURIComponent(token)}` : ''}`,
      downloadUrl: `/api/file/${uid}${token ? `?key=${encodeURIComponent(token)}` : ''}`,
      uploadedAt: attachment.uploaded_at
    }
  }

  selectedImage.value = upload
  isImageModalOpen.value = true
}

const handleImageDownload = async (image: UploadMetadata) => {
  const { getFile } = useUploadApi()
  try {
    const blob = await getFile(image.uid, false)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = image.fileName
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('Download failed:', error)
    alert('Failed to download file. Please try again.')
    throw error
  }
}

const handleSendMessage = async () => {
  if (!messageText.value.trim() && !selectedFile.value) {
    return
  }

  let uploadResponse = null

  // Upload file if selected
  if (selectedFile.value) {
    isUploading.value = true
    try {
      uploadResponse = await uploadFile(selectedFile.value)
      
      // Store upload metadata
      const uid = extractUidFromUrl(uploadResponse.download_url)
      if (uid) {
        const token = getToken()
        // Convert backend URLs to use our proxy endpoint
        const previewUrl = `/api/file/${uid}?preview=true${token ? `&key=${encodeURIComponent(token)}` : ''}`
        const downloadUrl = `/api/file/${uid}${token ? `?key=${encodeURIComponent(token)}` : ''}`
        
        addUpload({
          uid,
          fileName: uploadResponse.file_name,
          previewUrl,
          downloadUrl,
          uploadedAt: uploadResponse.uploaded_at
        })
      }

      sendMessage(
        messageText.value || `Uploaded ${uploadResponse.file_name}`,
        [uploadResponse]
      )
    } catch (error) {
      console.error('Upload failed:', error)
      alert('Failed to upload file. Please try again.')
      isUploading.value = false
      return
    } finally {
      isUploading.value = false
      selectedFile.value = null
    }
  } else {
    sendMessage(messageText.value)
  }

  messageText.value = ''
  nextTick(() => {
    scrollToBottom()
  })
}

const extractUidFromUrl = (url: string | undefined | null): string | null => {
  if (!url) return null
  const match = url.match(/\/file\/([^/?]+)/)
  return match ? (match[1] || null) : null
}

const getPreviewUrl = (attachment: UploadResponse): string => {
  const uid = extractUidFromUrl(attachment.preview_url || attachment.download_url)
  if (uid) {
    const token = getToken()
    if (token) {
      return `/api/file/${uid}?preview=true&key=${encodeURIComponent(token)}`
    }
  }
  return attachment.preview_url || ''
}

const formatTime = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  }).format(date)
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

// Auto-scroll when messages change
watch(messages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })
</script>
