<template>
  <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-2 md:gap-4">
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
          />
          <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <template #footer>
          <div class="p-2">
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

defineProps<Props>()

defineEmits<{
  'image-click': [upload: UploadMetadata]
}>()

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

