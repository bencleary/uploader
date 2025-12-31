<template>
  <div class="h-full w-full overflow-y-auto">
        <!-- Demo Warning Banner -->
    <Header title="Gallery" icon="i-lucide-images" @logout="">
      <template #actions>
        <UButton
          to="/"
          icon="i-lucide-arrow-left"
          variant="ghost"
          color="neutral"
          size="lg"
        >
          Back to Chat
        </UButton>
      </template>
    </Header>

    <div class="p-4">
      <UAlert
      color="warning"
      variant="soft"
      icon="i-lucide-alert-triangle"
      title="Demo Environment"
      description="This is a demo authentication service. The encryption key is generated client-side and should not be used in production."
      class=""
      :close-button="{ icon: 'i-lucide-x', color: 'amber', variant: 'link', 'aria-label': 'Close' }"
    />
    </div>

    <!-- Content -->
    <div class="px-4 py-4">
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
          color="primary"
          size="lg"
        >
          Go to Chat
        </UButton>
      </div>

      <!-- Photo Grid -->
      <ImageGrid
        v-else
        :uploads="Array.from(uploads)"
        @image-click="openModal"
      />
    </div>

    <!-- Full Image Modal -->
    <ImageModal
      v-model="isModalOpen"
      :image="selectedUpload"
      :on-download="handleDownload"
    />
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

</script>

