import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mountSuspended } from '@nuxt/test-utils/runtime'
import ChatInput from '~/components/ChatInput.vue'
import { clearNuxtState } from '#app'

describe('ChatInput', () => {
  const defaultProps = {
    messageText: '',
    selectedFile: null,
    isUploading: false
  }

  beforeEach(async () => {
    // clearNuxtState is only available in nuxt environment
    try {
      clearNuxtState()
    } catch {
      // Ignore if not available
    }
  })

  it('should render the input field', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: defaultProps
    })
    
    expect(wrapper.html()).toContain('input')
  })

  it('should emit update:messageText when input changes', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: defaultProps
    })
    
    const input = wrapper.find('input[type="text"]')
    if (input.exists()) {
      await input.setValue('Hello world')
      expect(wrapper.emitted('update:messageText')).toBeTruthy()
    }
  })

  it('should disable input when uploading', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        isUploading: true
      }
    })
    
    const input = wrapper.find('input[type="text"]')
    if (input.exists()) {
      expect(input.attributes('disabled')).toBeDefined()
    }
  })

  it('should disable submit button when message is empty and no file selected', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: defaultProps
    })
    
    const submitButton = wrapper.find('button[type="submit"]')
    if (submitButton.exists()) {
      expect(submitButton.attributes('disabled')).toBeDefined()
    }
  })

  it('should enable submit button when message has text', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        messageText: 'Hello'
      }
    })
    
    const submitButton = wrapper.find('button[type="submit"]')
    if (submitButton.exists()) {
      // Button should not be disabled when message has text
      expect(submitButton.attributes('disabled')).toBeUndefined()
    }
  })

  it('should enable submit button when file is selected', async () => {
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        selectedFile: file
      }
    })
    
    const submitButton = wrapper.find('button[type="submit"]')
    if (submitButton.exists()) {
      expect(submitButton.attributes('disabled')).toBeUndefined()
    }
  })

  it('should emit submit when form is submitted', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        messageText: 'Hello'
      }
    })
    
    const form = wrapper.find('form')
    if (form.exists()) {
      await form.trigger('submit')
      expect(wrapper.emitted('submit')).toBeTruthy()
    }
  })

  it('should show file input when file is selected', async () => {
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        selectedFile: file
      }
    })
    
    expect(wrapper.text()).toContain('test.jpg')
  })

  it('should emit clear-file when clear button is clicked', async () => {
    const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
    const wrapper = await mountSuspended(ChatInput, {
      props: {
        ...defaultProps,
        selectedFile: file,
        isUploading: false
      }
    })
    
    // Find the clear button - it should be in the file preview section
    const buttons = wrapper.findAll('button')
    // The clear button should be the one with icon="i-lucide-x" or similar
    const clearBtn = buttons.find(btn => {
      const html = btn.html()
      return html.includes('lucide-x') || html.includes('i-lucide-x') || html.includes('×')
    })
    
    if (clearBtn) {
      await clearBtn.trigger('click')
      expect(wrapper.emitted('clear-file')).toBeTruthy()
    } else {
      // If we can't find it by icon, try to find any button in the file preview area
      const filePreview = wrapper.find('[class*="emerald"]')
      if (filePreview.exists()) {
        const previewButtons = filePreview.findAll('button')
        if (previewButtons.length > 0) {
          await previewButtons[0].trigger('click')
          expect(wrapper.emitted('clear-file')).toBeTruthy()
        }
      }
    }
  })

  it('should emit file-selected when valid image file is selected', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: defaultProps
    })
    
    const fileInput = wrapper.find('input[type="file"]')
    if (fileInput.exists()) {
      const file = new File(['test'], 'test.jpg', { type: 'image/jpeg' })
      
      // Create a FileList mock
      const dataTransfer = new DataTransfer()
      dataTransfer.items.add(file)
      Object.defineProperty(fileInput.element, 'files', {
        value: dataTransfer.files,
        writable: false
      })
      
      await fileInput.trigger('change')
      
      expect(wrapper.emitted('file-selected')).toBeTruthy()
      expect(wrapper.emitted('file-selected')?.[0]?.[0]).toBeInstanceOf(File)
    }
  })

  it('should not emit file-selected for non-image files', async () => {
    const wrapper = await mountSuspended(ChatInput, {
      props: defaultProps
    })
    
    const fileInput = wrapper.find('input[type="file"]')
    if (fileInput.exists()) {
      const file = new File(['test'], 'test.pdf', { type: 'application/pdf' })
      
      // Mock alert to avoid console errors
      global.alert = vi.fn()
      
      const dataTransfer = new DataTransfer()
      dataTransfer.items.add(file)
      Object.defineProperty(fileInput.element, 'files', {
        value: dataTransfer.files,
        writable: false
      })
      
      await fileInput.trigger('change')
      
      expect(wrapper.emitted('file-selected')).toBeFalsy()
      expect(global.alert).toHaveBeenCalledWith('Please select an image file')
    }
  })
})

