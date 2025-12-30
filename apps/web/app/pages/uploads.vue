<template>
  <div class="min-h-screen bg-zinc-50 dark:bg-zinc-950">
    <!-- Header -->
    <div class="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 px-6 py-6">
      <div class="max-w-7xl mx-auto flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center">
            <UIcon name="i-lucide-images" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" />
          </div>
          <div>
            <h1 class="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Gallery</h1>
            <p class="text-sm text-zinc-500 dark:text-zinc-400 mt-0.5">
              {{ uploads.length }} {{ uploads.length === 1 ? 'image' : 'images' }}
            </p>
          </div>
        </div>
        <UButton
          to="/"
          icon="i-lucide-arrow-left"
          variant="ghost"
          color="zinc"
          size="lg"
        >
          Back to Chat
        </UButton>
      </div>
    </div>

    <!-- Content -->
    <div class="max-w-7xl mx-auto px-6 py-8">
      <!-- Empty State -->
      <div
        v-if="uploads.length === 0"
        class="flex flex-col items-center justify-center py-20"
      >
        <div class="w-24 h-24 rounded-2xl bg-emerald-500/10 flex items-center justify-center mb-6">
          <UIcon name="i-lucide-image-off" class="w-12 h-12 text-emerald-600 dark:text-emerald-400" />
        </div>
        <h2 class="text-2xl font-semibold text-zinc-900 dark:text-zinc-50 mb-2">No uploads yet</h2>
        <p class="text-zinc-500 dark:text-zinc-400 mb-6 max-w-sm text-center">
          Upload images in the chat to see them displayed here in a beautiful gallery
        </p>
        <UButton
          to="/"
          icon="i-lucide-arrow-left"
          color="emerald"
          size="lg"
        >
          Go to Chat
        </UButton>
      </div>

      <!-- Photo Grid -->
      <div
        v-else
        class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4"
      >
        <div
          v-for="upload in uploads"
          :key="upload.uid"
          class="group cursor-pointer"
          @click="openModal(upload)"
        >
          <UCard class="overflow-hidden border-zinc-200 dark:border-zinc-800 hover:border-emerald-300 dark:hover:border-emerald-700 transition-all hover:shadow-xl">
            <div class="aspect-square relative overflow-hidden bg-zinc-100 dark:bg-zinc-900">
              <img
                :src="upload.previewUrl"
                :alt="upload.fileName"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              />
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
    </div>

    <!-- Full Image Modal -->
    <UModal
      v-model="isModalOpen"
      :ui="{ width: 'max-w-5xl', padding: 'p-0' }"
    >
      <UCard
        v-if="selectedUpload"
        class="overflow-hidden border-0"
      >
        <template #header>
          <div class="flex items-center justify-between px-6 py-4 border-b border-zinc-200 dark:border-zinc-800">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg bg-emerald-500/10 flex items-center justify-center">
                <UIcon name="i-lucide-image" class="w-5 h-5 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <h3 class="font-semibold text-zinc-900 dark:text-zinc-50">{{ selectedUpload.fileName }}</h3>
                <p class="text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">
                  Uploaded {{ formatDate(selectedUpload.uploadedAt) }}
                </p>
              </div>
            </div>
            <UButton
              icon="i-lucide-x"
              variant="ghost"
              color="zinc"
              @click="isModalOpen = false"
            />
          </div>
        </template>

        <div class="flex flex-col items-center bg-zinc-50 dark:bg-zinc-900 p-6">
          <img
            :src="selectedUpload.downloadUrl"
            :alt="selectedUpload.fileName"
            class="max-w-full max-h-[75vh] object-contain rounded-lg shadow-lg"
          />
        </div>

        <template #footer>
          <div class="flex justify-end items-center px-6 py-4 border-t border-zinc-200 dark:border-zinc-800">
            <UButton
              @click="handleDownload(selectedUpload)"
              :loading="isDownloading"
              icon="i-lucide-download"
              color="emerald"
              size="lg"
            >
              Download Image
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { UploadMetadata } from '~/types/api'

const { uploads } = useUploads()
const { getFile } = useUploadApi()

const isModalOpen = ref(false)
const selectedUpload = ref<UploadMetadata | null>(null)
const isDownloading = ref(false)

const openModal = (upload: UploadMetadata) => {
  selectedUpload.value = upload
  isModalOpen.value = true
}

const handleDownload = async (upload: UploadMetadata) => {
  isDownloading.value = true
  try {
    const blob = await getFile(upload.uid, false)
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = upload.fileName
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('Download failed:', error)
    alert('Failed to download file. Please try again.')
  } finally {
    isDownloading.value = false
  }
}

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

