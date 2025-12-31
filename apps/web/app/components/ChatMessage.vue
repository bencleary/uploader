<template>
  <div 
    :class="[
      'flex gap-4 group',
      isFromUser ? 'flex-row' : 'flex-row-reverse ml-auto max-w-[85%]'
    ]"
  >
    <UAvatar
      :alt="isFromUser ? 'You' : 'Other User'"
      :src="isFromUser 
        ? 'https://i.pinimg.com/736x/de/8e/b6/de8eb693d13454cfa41d81d0d8bd2778.jpg'
        : 'https://api.dicebear.com/7.x/avataaars/svg?seed=OtherUser'"
      size="md"
      :class="[
        'ring-2 bg-zinc-50 dark:bg-zinc-900 shrink-0',
        isFromUser ? 'ring-emerald-500/20' : 'ring-blue-500/20'
      ]"
    />
    <div 
      :class="[
        'min-w-0',
        isFromUser ? 'flex-1' : 'flex-1 flex flex-col items-end'
      ]"
    >
      <div 
        :class="[
          'flex items-center gap-2 mb-2',
          isFromUser ? '' : 'flex-row-reverse'
        ]"
      >
        <span class="font-semibold text-sm text-zinc-900 dark:text-zinc-50">
          {{ isFromUser ? 'You' : 'Other User' }}
        </span>
        <span class="text-xs text-zinc-500 dark:text-zinc-400">
          {{ formatTime(message.timestamp) }}
        </span>
      </div>
      <div
        v-if="message.text"
        :class="[
          'px-4 py-2.5 rounded-2xl shadow-sm mb-2',
          isFromUser 
            ? 'inline-block rounded-tl-sm bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700'
            : 'inline-block rounded-tr-sm bg-blue-500 dark:bg-blue-600 text-white'
        ]"
      >
        <p 
          :class="[
            'text-sm whitespace-pre-wrap',
            isFromUser ? 'text-zinc-900 dark:text-zinc-50' : 'text-white'
          ]"
        >
          {{ message.text }}
        </p>
      </div>
      <div
        v-if="message.attachments && message.attachments.length > 0"
        class="space-y-3"
        :class="isFromUser ? '' : 'flex flex-col items-end'"
      >
        <div
          v-for="(attachment, idx) in message.attachments"
          :key="idx"
          class="relative group/image cursor-pointer"
          @click="$emit('image-click', attachment)"
        >
          <div class="rounded-xl overflow-hidden border border-zinc-200 dark:border-zinc-700 shadow-md hover:shadow-lg transition-shadow max-w-sm">
            <img
              :src="getPreviewUrl(attachment)"
              :alt="attachment.file_name"
              class="max-h-80 object-cover w-full h-full max-w-sm"
            />
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

// Default to true (user message) if not specified for backward compatibility
const isFromUser = computed(() => props.message.isFromUser !== false)

const formatTime = (date: Date): string => {
  return new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    minute: '2-digit'
  }).format(date)
}
</script>

