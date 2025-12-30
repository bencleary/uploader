<template>
  <div class="flex gap-4 group">
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
          class="relative group/image cursor-pointer"
          @click="$emit('image-click', attachment)"
        >
          <div class="rounded-xl overflow-hidden border border-zinc-200 dark:border-zinc-700 shadow-md hover:shadow-lg transition-shadow max-w-sm relative">
            <img
              :src="getPreviewUrl(attachment)"
              :alt="attachment.file_name"
              class="max-h-80 object-cover w-full h-full max-w-sm"
              @load="imageLoadingStates[idx] = false"
              @error="imageLoadingStates[idx] = false; imageErrorStates[idx] = true"
            />
            <div
              v-if="imageLoadingStates[idx] !== false"
              class="absolute inset-0 flex items-center justify-center bg-zinc-100/80 dark:bg-zinc-900/80 backdrop-blur-sm"
            >
              <UIcon
                name="i-lucide-loader-2"
                class="w-8 h-8 text-emerald-600 dark:text-emerald-400 animate-spin"
              />
            </div>
            <div
              v-if="imageErrorStates[idx]"
              class="absolute inset-0 flex items-center justify-center bg-zinc-100 dark:bg-zinc-900"
            >
              <div class="text-center">
                <UIcon
                  name="i-lucide-alert-circle"
                  class="w-8 h-8 text-zinc-400 dark:text-zinc-600 mx-auto mb-2"
                />
                <p class="text-xs text-zinc-500 dark:text-zinc-400">Failed to load image</p>
              </div>
            </div>
          </div>
          <p class="text-xs text-zinc-500 dark:text-zinc-400 mt-2 px-1">{{ attachment.file_name }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ChatMessage, UploadResponse } from '~/types/api'

interface Props {
  message: ChatMessage
  getPreviewUrl: (attachment: UploadResponse) => string
}

const props = defineProps<Props>()

defineEmits<{
  'image-click': [attachment: UploadResponse]
}>()

const imageLoadingStates = reactive<Record<number, boolean>>({})
const imageErrorStates = reactive<Record<number, boolean>>({})

watch(() => props.message.attachments, (attachments) => {
  if (attachments) {
    attachments.forEach((_, idx) => {
      imageLoadingStates[idx] = true
      imageErrorStates[idx] = false
    })
  }
}, { immediate: true })

const formatTime = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  }).format(date)
}
</script>

