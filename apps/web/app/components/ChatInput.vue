<template>
  <div class="border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-4">
    <form
      @submit.prevent="$emit('submit')"
      class="flex gap-3 items-end"
    >
      <div class="flex-1 relative">
        <UInput
          :model-value="messageText"
          @update:model-value="$emit('update:messageText', $event)"
          placeholder="Type a message..."
          :disabled="isUploading"
          size="lg"
          color="zinc"
          :ui="{ root: 'rounded-xl' }"
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
        :ui="{ root: 'rounded-xl' }"
      />
      <UButton
        type="submit"
        :loading="isUploading"
        :disabled="!messageText.trim() && !selectedFile"
        color="emerald"
        size="lg"
        :ui="{ root: 'rounded-xl' }"
      >
        <UIcon name="i-lucide-send" class="w-4 h-4" />
      </UButton>
    </form>
    <div
      v-if="selectedFile"
      :class="[
        'mt-3 flex items-center gap-2 px-3 py-2 border rounded-lg transition-colors',
        isUploading
          ? 'bg-emerald-100 dark:bg-emerald-950/40 border-emerald-300 dark:border-emerald-800'
          : 'bg-emerald-50 dark:bg-emerald-950/20 border-emerald-200 dark:border-emerald-900'
      ]"
    >
      <UIcon
        :name="isUploading ? 'i-lucide-loader-2' : 'i-lucide-file-image'"
        :class="[
          'w-4 h-4 text-emerald-600 dark:text-emerald-400',
          isUploading && 'animate-spin'
        ]"
      />
      <span class="text-sm text-emerald-900 dark:text-emerald-200 flex-1 truncate">{{ selectedFile.name }}</span>
      <UButton
        v-if="!isUploading"
        icon="i-lucide-x"
        size="xs"
        variant="ghost"
        color="emerald"
        @click="$emit('clear-file')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  messageText: string
  selectedFile: File | null
  isUploading: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  'update:messageText': [value: string]
  'file-selected': [file: File]
  'clear-file': []
  'submit': []
}>()

const fileInput = ref<HTMLInputElement | null>(null)

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    if (file && file.type.startsWith('image/')) {
      emit('file-selected', file)
    } else if (file) {
      alert('Please select an image file')
    }
  }
}
</script>

