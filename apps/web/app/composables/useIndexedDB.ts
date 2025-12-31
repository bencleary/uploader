const DB_NAME = 'uploader_db'
const DB_VERSION = 1
const STORE_NAME = 'storage'

let dbPromise: Promise<IDBDatabase> | null = null

function getDatabase(): Promise<IDBDatabase> {
  if (!dbPromise) {
    dbPromise = new Promise((resolve, reject) => {
      if (typeof indexedDB === 'undefined') {
        reject(new Error('IndexedDB is not supported'))
        return
      }

      const request = indexedDB.open(DB_NAME, DB_VERSION)

      request.onerror = () => {
        reject(request.error)
      }

      request.onsuccess = () => {
        resolve(request.result)
      }

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result
        if (!db.objectStoreNames.contains(STORE_NAME)) {
          db.createObjectStore(STORE_NAME)
        }
      }
    })
  }

  return dbPromise
}

// Export function to reset database promise for testing
export function resetIndexedDB() {
  dbPromise = null
}

export const useIndexedDB = () => {
  const get = async <T>(key: string): Promise<T | null> => {
    if (!import.meta.client) {
      return null
    }

    try {
      const db = await getDatabase()
      return new Promise<T | null>((resolve, reject) => {
        const transaction = db.transaction(STORE_NAME, 'readonly')
        const store = transaction.objectStore(STORE_NAME)
        const request = store.get(key)

        request.onerror = () => {
          reject(request.error)
        }

        request.onsuccess = () => {
          resolve(request.result || null)
        }
      })
    } catch (error) {
      console.error('IndexedDB get error:', error)
      return null
    }
  }

  const set = async <T>(key: string, value: T): Promise<void> => {
    if (!import.meta.client) {
      return
    }

    try {
      const db = await getDatabase()
      return new Promise<void>((resolve, reject) => {
        const transaction = db.transaction(STORE_NAME, 'readwrite')
        const store = transaction.objectStore(STORE_NAME)
        // Convert to plain object/array by serializing to JSON and parsing back
        // This removes Vue reactivity proxies and ensures IndexedDB can clone the data
        const plainValue = JSON.parse(JSON.stringify(value))
        const request = store.put(plainValue, key)

        request.onerror = () => {
          reject(request.error)
        }

        request.onsuccess = () => {
          resolve()
        }
      })
    } catch (error) {
      console.error('IndexedDB set error:', error)
      throw error
    }
  }

  const remove = async (key: string): Promise<void> => {
    if (!import.meta.client) {
      return
    }

    try {
      const db = await getDatabase()
      return new Promise<void>((resolve, reject) => {
        const transaction = db.transaction(STORE_NAME, 'readwrite')
        const store = transaction.objectStore(STORE_NAME)
        const request = store.delete(key)

        request.onerror = () => {
          reject(request.error)
        }

        request.onsuccess = () => {
          resolve()
        }
      })
    } catch (error) {
      console.error('IndexedDB delete error:', error)
      throw error
    }
  }

  return {
    get,
    set,
    remove
  }
}

