export default defineAppConfig({
  ui: {
    primary: 'emerald',
    gray: 'zinc',
    modal: {
      slots: {
        overlay: 'fixed inset-0 bg-zinc-950/75 dark:bg-zinc-950/75',
        content: 'bg-white dark:bg-zinc-900 divide-y divide-zinc-200 dark:divide-zinc-800 flex flex-col focus:outline-none',
        header: 'flex items-center gap-1.5 p-4 sm:px-6 min-h-16 bg-white dark:bg-zinc-900',
        body: 'flex-1 p-4 sm:p-6 bg-white dark:bg-zinc-900',
        footer: 'flex items-center gap-1.5 p-4 sm:px-6 bg-white dark:bg-zinc-900',
        title: 'text-zinc-900 dark:text-zinc-50 font-semibold',
        description: 'mt-1 text-zinc-500 dark:text-zinc-400 text-sm',
        close: 'absolute top-4 end-4 text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-50'
      },
      variants: {
        transition: {
          true: {
            overlay: 'data-[state=open]:animate-[fade-in_200ms_ease-out] data-[state=closed]:animate-[fade-out_200ms_ease-in]',
            content: 'data-[state=open]:animate-[scale-in_200ms_ease-out] data-[state=closed]:animate-[scale-out_200ms_ease-in]'
          }
        },
        fullscreen: {
          true: {
            content: 'inset-0'
          },
          false: {
            content: 'w-[calc(100vw-2rem)] max-w-lg rounded-lg shadow-lg ring-1 ring-zinc-200 dark:ring-zinc-800'
          }
        },
        overlay: {
          true: {
            overlay: 'bg-zinc-950/75 dark:bg-zinc-950/75'
          }
        },
        scrollable: {
          true: {
            overlay: 'overflow-y-auto',
            content: 'relative'
          },
          false: {
            content: 'fixed',
            body: 'overflow-y-auto'
          }
        }
      },
      compoundVariants: [
        {
          scrollable: true,
          fullscreen: false,
          class: {
            overlay: 'grid place-items-center p-4 sm:py-8'
          }
        },
        {
          scrollable: false,
          fullscreen: false,
          class: {
            content: 'top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 max-h-[calc(100dvh-2rem)] sm:max-h-[calc(100dvh-4rem)] overflow-hidden'
          }
        }
      ]
    }
  }
})
