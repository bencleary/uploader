<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
    <div
      v-for="upload in uploads"
      :key="upload.uid"
      class="group cursor-pointer"
      @click="$emit('image-click', upload)"
    >
      <UCard class="overflow-hidden border-zinc-200 dark:border-zinc-800 hover:border-emerald-300 dark:hover:border-emerald-700 transition-all hover:shadow-xl bg-white dark:bg-zinc-900">
        <div class="aspect-square relative overflow-hidden bg-zinc-100 dark:bg-zinc-900">
          <img
            :src="upload.previewUrl"
            :alt="upload.fileName"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
            @load="imageLoadingStates[upload.uid] = false"
            @error="imageLoadingStates[upload.uid] = false; imageErrorStates[upload.uid] = true"
          />
          <div
            v-if="imageLoadingStates[upload.uid] !== false"
            class="absolute inset-0 flex items-center justify-center bg-zinc-100/80 dark:bg-zinc-900/80 backdrop-blur-sm z-10"
          >
            <UIcon
              name="i-lucide-loader-2"
              class="w-8 h-8 text-emerald-600 dark:text-emerald-400 animate-spin"
            />
          </div>
          <div
            v-if="imageErrorStates[upload.uid]"
            class="absolute inset-0 flex items-center justify-center bg-zinc-100 dark:bg-zinc-900 z-10"
          >
            <div class="text-center">
              <UIcon
                name="i-lucide-alert-circle"
                class="w-8 h-8 text-zinc-400 dark:text-zinc-600 mx-auto mb-2"
              />
              <p class="text-xs text-zinc-500 dark:text-zinc-400">Failed to load</p>
            </div>
          </div>
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <template #footer>
          <div class="p-3">
            <p class="text-sm font-medium text-zinc-900 dark:text-zinc-50 truncate mb-1">
              {{ upload.fileName }}
            </p>
            <p class="text-xs text-zinc-500 dark:text-zinc-400">
              {{ formatDate(upload.uploadedAt) }}
            </p>
          </div>
        </template>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UploadMetadata } from '~/types/api'

interface Props {
  uploads: UploadMetadata[]
}

const props = defineProps<Props>()

defineEmits<{
  'image-click': [upload: UploadMetadata]
}>()

const imageLoadingStates = reactive<Record<string, boolean>>({})
const imageErrorStates = reactive<Record<string, boolean>>({})

watch(() => props.uploads, (uploads) => {
  uploads.forEach((upload) => {
    imageLoadingStates[upload.uid] = true
    imageErrorStates[upload.uid] = false
  })
}, { immediate: true })

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit'
  }).format(date)
}
</script>

