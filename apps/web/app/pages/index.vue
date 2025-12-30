<template>
  <div class="flex flex-col h-screen bg-zinc-50 dark:bg-zinc-950">
    <!-- Demo Warning Banner -->
    <UAlert
      v-if="isAuthenticated"
      color="amber"
      variant="soft"
      icon="i-lucide-alert-triangle"
      title="Demo Environment"
      description="This is a demo authentication service. The encryption key is generated client-side and should not be used in production."
      class="m-4 mb-0"
      :close-button="{ icon: 'i-lucide-x', color: 'amber', variant: 'link', 'aria-label': 'Close' }"
    />

    <!-- Login Section -->
    <div
      v-if="!isAuthenticated"
      class="flex-1 flex items-center justify-center p-6 bg-gradient-to-br from-zinc-50 to-zinc-100 dark:from-zinc-950 dark:to-zinc-900"
    >
      <UCard class="max-w-md w-full shadow-xl border-zinc-200 dark:border-zinc-800">
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
            color="emerald"
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
      class="flex-1 flex flex-col h-full max-w-4xl mx-auto w-full"
    >
      <!-- Header -->
      <div class="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 px-6 py-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-full bg-emerald-500/10 flex items-center justify-center">
            <UIcon name="i-lucide-message-circle" class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
          </div>
          <div>
            <h1 class="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Chat</h1>
            <p class="text-xs text-zinc-500 dark:text-zinc-400">Share images and messages</p>
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

        <div
          v-for="message in messages"
          :key="message.id"
          class="flex gap-4 group"
        >
          <UAvatar
            :alt="'User'"
            size="md"
            class="ring-2 ring-emerald-500/20"
          />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-2">
              <span class="font-semibold text-sm text-zinc-900 dark:text-zinc-50">You</span>
              <span class="text-xs text-zinc-500 dark:text-zinc-400">
                {{ formatTime(message.timestamp) }}
              </span>
            </div>
            <div
              v-if="message.text"
              class="inline-block px-4 py-2.5 rounded-2xl rounded-tl-sm bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 shadow-sm mb-2"
            >
              <p class="text-sm text-zinc-900 dark:text-zinc-50 whitespace-pre-wrap">{{ message.text }}</p>
            </div>
            <div
              v-if="message.attachments && message.attachments.length > 0"
              class="space-y-3"
            >
              <div
                v-for="(attachment, idx) in message.attachments"
                :key="idx"
                class="relative group/image"
              >
                <div class="rounded-xl overflow-hidden border border-zinc-200 dark:border-zinc-700 shadow-md hover:shadow-lg transition-shadow">
                  <img
                    :src="getPreviewUrl(attachment)"
                    :alt="attachment.file_name"
                    class="max-w-sm max-h-80 object-cover w-full"
                  />
                </div>
                <p class="text-xs text-zinc-500 dark:text-zinc-400 mt-2 px-1">{{ attachment.file_name }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Input Area -->
      <div class="border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 px-6 py-4">
        <form
          @submit.prevent="handleSendMessage"
          class="flex gap-3 items-end"
        >
          <div class="flex-1 relative">
            <UInput
              v-model="messageText"
              placeholder="Type a message..."
              :disabled="isUploading"
              size="lg"
              color="zinc"
              :ui="{ rounded: 'rounded-xl' }"
            />
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              class="hidden"
              @change="handleFileSelect"
            />
          </div>
          <UButton
            type="button"
            icon="i-lucide-paperclip"
            variant="ghost"
            color="zinc"
            size="lg"
            :disabled="isUploading"
            @click="fileInput?.click()"
            :ui="{ rounded: 'rounded-xl' }"
          />
          <UButton
            type="submit"
            :loading="isUploading"
            :disabled="!messageText.trim() && !selectedFile"
            color="emerald"
            size="lg"
            :ui="{ rounded: 'rounded-xl' }"
          >
            <UIcon name="i-lucide-send" class="w-4 h-4" />
          </UButton>
        </form>
        <div
          v-if="selectedFile"
          class="mt-3 flex items-center gap-2 px-3 py-2 bg-emerald-50 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-900 rounded-lg"
        >
          <UIcon name="i-lucide-file-image" class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
          <span class="text-sm text-emerald-900 dark:text-emerald-200 flex-1 truncate">{{ selectedFile.name }}</span>
          <UButton
            icon="i-lucide-x"
            size="xs"
            variant="ghost"
            color="emerald"
            @click="selectedFile = null"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated, login, logout, getToken } = useAuth()
const { sendMessage, messages } = useChat()
const { uploadFile } = useUploadApi()
const { addUpload } = useUploads()

const messageText = ref('')
const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)
const isLoggingIn = ref(false)
const isUploading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

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

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    // Validate image type
    if (file.type.startsWith('image/')) {
      selectedFile.value = file
    } else {
      alert('Please select an image file')
    }
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
      if (fileInput.value) {
        fileInput.value.value = ''
      }
    }
  } else {
    sendMessage(messageText.value)
  }

  messageText.value = ''
  nextTick(() => {
    scrollToBottom()
  })
}

const extractUidFromUrl = (url: string): string | null => {
  const match = url.match(/\/file\/([^/?]+)/)
  return match ? match[1] : null
}

const getPreviewUrl = (attachment: any): string => {
  const uid = extractUidFromUrl(attachment.preview_url || attachment.download_url)
  if (uid) {
    const token = getToken()
    return `/api/file/${uid}?preview=true${token ? `&key=${encodeURIComponent(token)}` : ''}`
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
