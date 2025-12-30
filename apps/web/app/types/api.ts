export interface UploadResponse {
  file_name: string
  preview_url: string
  download_url: string
  uploaded_at: string
}

export interface AuthResponse {
  token: string
  message?: string
}

export interface ChatMessage {
  id: string
  text: string
  timestamp: Date
  attachments?: UploadResponse[]
}

export interface UploadMetadata {
  uid: string
  fileName: string
  previewUrl: string
  downloadUrl: string
  uploadedAt: string
}

