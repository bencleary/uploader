function firstQueryValue(value: string | string[] | undefined): string | undefined {
  if (Array.isArray(value)) {
    return value[0]
  }
  return value
}

function getErrorMessage(error: unknown): string {
  if (error && typeof error === 'object') {
    const record = error as Record<string, unknown>
    if (typeof record.message === 'string') {
      return record.message
    }
  }
  return 'Failed to fetch file'
}

function getErrorStatusCode(error: unknown): number {
  if (error && typeof error === 'object') {
    const record = error as Record<string, unknown>
    if (typeof record.statusCode === 'number') {
      return record.statusCode
    }
  }
  return 500
}

export default defineEventHandler(async (event) => {
  const uid = getRouterParam(event, 'uid')
  const query = getQuery(event)
  const preview = firstQueryValue(query.preview as string | string[] | undefined) === 'true'

  if (!uid) {
    throw createError({
      statusCode: 400,
      message: 'UID is required'
    })
  }

  // Get the key from the request (should be passed as a query param or header)
  const keyFromHeader = getHeader(event, 'key')
  const keyFromQuery = firstQueryValue(query.key as string | string[] | undefined)
  const key = keyFromHeader || keyFromQuery

  if (!key) {
    throw createError({
      statusCode: 401,
      message: 'Authentication key is required'
    })
  }

  const config = useRuntimeConfig()
  const baseURL = config.public.apiBaseUrl || 'http://localhost:1323'
  const url = `${baseURL}/file/${uid}${preview ? '?preview=true' : ''}`

  try {
    const response = await fetch(url, {
      headers: {
        key
      }
    })

    if (!response.ok) {
      const errorText = await response.text()
      console.error(`Backend error (${response.status}):`, errorText)
      throw createError({
        statusCode: response.status,
        message: `Failed to fetch file: ${response.statusText}`
      })
    }

    // Get content type from response headers (Go backend sets this)
    const contentType = response.headers.get('content-type') || 'application/octet-stream'
    
    // Read the entire response body
    const buffer = await response.arrayBuffer()
    
    // Check if we actually got data
    if (buffer.byteLength === 0) {
      throw createError({
        statusCode: 500,
        message: 'Unexpected empty response from backend'
      })
    }

    // Set headers - must be set before returning
    setHeader(event, 'Content-Type', contentType)
    // Key can be provided via query param (e.g. <img src="...">), so never allow shared caching.
    setHeader(event, 'Cache-Control', 'private, no-store')
    setHeader(event, 'Pragma', 'no-cache')

    // Convert to Uint8Array for H3 compatibility
    // H3 should automatically detect binary data based on content-type
    return new Uint8Array(buffer)
  } catch (error: unknown) {
    console.error('Proxy error:', error)
    throw createError({
      statusCode: getErrorStatusCode(error),
      message: getErrorMessage(error)
    })
  }
})
