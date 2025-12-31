export default defineEventHandler(async (event) => {
  const uid = getRouterParam(event, 'uid')
  const query = getQuery(event)
  const preview = query.preview === 'true'

  if (!uid) {
    throw createError({
      statusCode: 400,
      message: 'UID is required'
    })
  }

  // Get the key from the request (should be passed as a query param or header)
  const key = getHeader(event, 'key') || getQuery(event).key as string

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
        message: `Failed to fetch file: ${response.statusText} - ${errorText}`
      })
    }

    // Get content type from response headers (Go backend sets this)
    const contentType = response.headers.get('content-type') || 'application/octet-stream'
    
    // Read the entire response body
    const buffer = await response.arrayBuffer()
    
    // Log for debugging
    console.log(`Fetched file: ${uid}, preview: ${preview}, size: ${buffer.byteLength} bytes, type: ${contentType}`)
    
    // Check if we actually got data
    if (buffer.byteLength === 0 || buffer.byteLength < 10) {
      const errorText = new TextDecoder().decode(buffer.slice(0, 100))
      console.error('Unexpected small response:', errorText)
      throw createError({
        statusCode: 500,
        message: `Unexpected response from backend: ${errorText || 'empty or too small'}`
      })
    }

    // Set headers - must be set before returning
    setHeader(event, 'Content-Type', contentType)
    setHeader(event, 'Cache-Control', 'public, max-age=3600')

    // Convert to Uint8Array for H3 compatibility
    // H3 should automatically detect binary data based on content-type
    return new Uint8Array(buffer)
  } catch (error: any) {
    console.error('Proxy error:', error)
    throw createError({
      statusCode: error.statusCode || 500,
      message: error.message || 'Failed to fetch file'
    })
  }
})

