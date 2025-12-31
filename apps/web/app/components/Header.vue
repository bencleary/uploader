<template>
        <div class="border-b border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-4 shrink-0">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center">
              <UIcon name="i-lucide-message-circle" class="w-6 h-6 text-emerald-600 dark:text-emerald-400" />
            </div>
            <div>
              <h1 class="text-3xl font-bold text-zinc-900 dark:text-zinc-50">{{ title }}</h1>
            </div>
          </div>
          <div class="flex gap-2">
            <slot name="actions" />
          </div>
        </div>
      </div>
</template>

<script setup lang="ts">
    const props = defineProps<{
        onLogout: () => void,
        title: string,
        icon: string,
    }>()

    const icon = computed(() => {
        return props.icon || 'i-lucide-message-circle'
    })

    const title = computed(() => {
        return props.title || 'Chat'
    })
</script>