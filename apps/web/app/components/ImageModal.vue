<template>
  <UModal
    v-model:open="isOpen"
    :title="image?.fileName"
    :description="image ? `Uploaded ${formatDate(image.uploadedAt)}` : ''"
    :ui="{ content: 'max-w-5xl', body: 'p-0' }"
  >
    <template #body>
      <div
        v-if="image"
        class="flex flex-col items-center bg-zinc-50 dark:bg-zinc-900 p-6 relative min-h-[400px]"
      >
        <img
          ref="imageElement"
          :src="image.downloadUrl"
          :alt="image.fileName"
          class="max-w-full max-h-[75vh] object-contain rounded-lg shadow-lg"
          @load="handleImageLoad"
          @error="handleImageError"
        />
        <div
          v-if="isImageLoading"
          class="absolute inset-0 flex items-center justify-center bg-zinc-50/80 dark:bg-zinc-900/80 backdrop-blur-sm"
        >
          <UIcon
            name="i-lucide-loader-2"
            class="w-12 h-12 text-emerald-600 dark:text-emerald-400 animate-spin"
          />
        </div>
        <div
          v-if="isImageError"
          class="absolute inset-0 flex items-center justify-center bg-zinc-50 dark:bg-zinc-900"
        >
          <div class="text-center">
            <UIcon
              name="i-lucide-alert-circle"
              class="w-12 h-12 text-zinc-400 dark:text-zinc-600 mx-auto mb-4"
            />
            <p class="text-sm text-zinc-500 dark:text-zinc-400">Failed to load image</p>
          </div>
        </div>
      </div>
    </template>

    <template #footer="{ close }">
      <div class="flex justify-end items-center gap-2">
        <UButton
          label="Cancel"
          color="neutral"
          variant="outline"
          @click="close"
        />
        <UButton
          v-if="image && onDownload"
          @click="handleDownload"
          :loading="isDownloading"
          icon="i-lucide-download"
          color="primary"
          size="lg"
        >
          Download Image
        </UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { UploadMetadata } from '~/types/api'

interface Props {
  image: UploadMetadata | null
  modelValue: boolean
  onDownload?: (image: UploadMetadata) => Promise<void>
}

const props = withDefaults(defineProps<Props>(), {
  onDownload: undefined
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const isDownloading = ref(false)
const isImageLoading = ref(true)
const isImageError = ref(false)
const imageElement = ref<HTMLImageElement | null>(null)

const isOpen = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const checkImageLoaded = () => {
  nextTick(() => {
    if (imageElement.value) {
      if (imageElement.value.complete && imageElement.value.naturalHeight !== 0) {
        // Image is already loaded (cached)
        isImageLoading.value = false
        isImageError.value = false
      } else if (imageElement.value.complete && imageElement.value.naturalHeight === 0) {
        // Image failed to load
        isImageLoading.value = false
        isImageError.value = true
      }
      // If not complete, wait for @load or @error event
    }
  })
}

const handleImageLoad = () => {
  isImageLoading.value = false
  isImageError.value = false
}

const handleImageError = () => {
  isImageLoading.value = false
  isImageError.value = true
}

watch(() => props.modelValue, (isOpen) => {
  if (isOpen && props.image) {
    // Reset loading state when modal opens
    isImageLoading.value = true
    isImageError.value = false
    // Check if image is already loaded (cached)
    checkImageLoaded()
  }
})

watch(() => props.image, (image) => {
  if (image) {
    isImageLoading.value = true
    isImageError.value = false
    // Check if image is already loaded (cached) after src updates
    checkImageLoaded()
  }
})

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

const handleDownload = async () => {
  if (props.image && props.onDownload) {
    isDownloading.value = true
    try {
      await props.onDownload(props.image)
    } finally {
      isDownloading.value = false
    }
  }
}
</script>

