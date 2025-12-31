import { beforeEach, vi } from 'vitest'

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {}

  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString()
    },
    removeItem: (key: string) => {
      delete store[key]
    },
    clear: () => {
      store = {}
    }
  }
})()

// Mock IndexedDB
const indexedDBStore: Record<string, any> = {}
const objectStoreNamesSet = new Set<string>()

const createIDBRequest = <T>(getResult: () => T): IDBRequest<T> => {
  const request = {
    result: undefined as T,
    error: null,
    readyState: 'pending' as IDBRequestReadyState,
    source: null,
    transaction: null,
    onerror: null as ((event: Event) => void) | null,
    onsuccess: null as ((event: Event) => void) | null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  } as IDBRequest<T>

  // Use setTimeout to simulate async behavior
  setTimeout(() => {
    request.result = getResult()
    request.readyState = 'done'
    if (request.onsuccess) {
      request.onsuccess({ target: request } as Event)
    }
  }, 0)

  return request
}

const indexedDBMock = {
  open: (name: string, version?: number): IDBOpenDBRequest => {
    const request = {
      result: null,
      error: null,
      readyState: 'pending' as IDBRequestReadyState,
      source: null,
      transaction: null,
      onerror: null as ((event: IDBVersionChangeEvent) => void) | null,
      onsuccess: null as ((event: IDBVersionChangeEvent) => void) | null,
      onblocked: null as ((event: IDBVersionChangeEvent) => void) | null,
      onupgradeneeded: null as ((event: IDBVersionChangeEvent) => void) | null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn()
    } as IDBOpenDBRequest

    // Simulate database creation
    setTimeout(() => {
      const db = {
        name,
        version: version || 1,
        objectStoreNames: {
          contains: (name: string) => objectStoreNamesSet.has(name),
          length: objectStoreNamesSet.size,
          [Symbol.iterator]: () => objectStoreNamesSet[Symbol.iterator]()
        },
        createObjectStore: (name: string) => {
          objectStoreNamesSet.add(name)
          return {
            name,
            keyPath: null,
            indexNames: { length: 0 },
            autoIncrement: false,
            createIndex: vi.fn(),
            deleteIndex: vi.fn(),
            get: vi.fn(),
            getKey: vi.fn(),
            getAll: vi.fn(),
            getAllKeys: vi.fn(),
            count: vi.fn(),
            put: vi.fn(),
            add: vi.fn(),
            delete: vi.fn(),
            clear: vi.fn(),
            openCursor: vi.fn(),
            openKeyCursor: vi.fn()
          }
        },
        transaction: (storeNames: string | string[], mode?: IDBTransactionMode) => {
          const stores = Array.isArray(storeNames) ? storeNames : [storeNames]
          return {
            objectStore: (name: string) => {
              if (!stores.includes(name)) {
                throw new Error(`Object store '${name}' not in transaction`)
              }
              return {
                get: (key: string) => {
                  return createIDBRequest(() => indexedDBStore[key] || undefined)
                },
                put: (value: any, key: string) => {
                  return createIDBRequest(() => {
                    indexedDBStore[key] = value
                    return key
                  })
                },
                delete: (key: string) => {
                  return createIDBRequest(() => {
                    delete indexedDBStore[key]
                    return undefined
                  })
                }
              }
            },
            mode: mode || 'readonly',
            db,
            error: null,
            onabort: null,
            oncomplete: null,
            onerror: null,
            abort: vi.fn(),
            objectStoreNames: { length: stores.length },
            addEventListener: vi.fn(),
            removeEventListener: vi.fn(),
            dispatchEvent: vi.fn()
          } as IDBTransaction
        },
        close: vi.fn(),
        deleteObjectStore: (name: string) => {
          objectStoreNamesSet.delete(name)
        },
        createIndex: vi.fn()
      } as IDBDatabase

      // Set result on request first (must be set before onupgradeneeded is called)
      request.result = db
      request.readyState = 'done'
      
      // Create event objects with explicit target property
      // The handler accesses event.target.result, so target must be the request
      const upgradeEvent: any = {}
      Object.defineProperty(upgradeEvent, 'target', {
        value: request,
        enumerable: true,
        writable: false,
        configurable: false
      })
      Object.defineProperty(upgradeEvent, 'currentTarget', {
        value: request,
        enumerable: true,
        writable: false,
        configurable: false
      })
      
      const successEvent: any = {}
      Object.defineProperty(successEvent, 'target', {
        value: request,
        enumerable: true,
        writable: false,
        configurable: false
      })
      Object.defineProperty(successEvent, 'currentTarget', {
        value: request,
        enumerable: true,
        writable: false,
        configurable: false
      })
      
      // Call onupgradeneeded first (before onsuccess) - this is when object stores are created
      // The handler will access event.target.result to get the database
      if (request.onupgradeneeded) {
        request.onupgradeneeded(upgradeEvent as IDBVersionChangeEvent)
      }
      // Then call onsuccess - this is when the database is ready
      if (request.onsuccess) {
        request.onsuccess(successEvent as IDBVersionChangeEvent)
      }
    }, 0)

    return request
  },
  deleteDatabase: vi.fn(),
  cmp: vi.fn()
} as IDBFactory

// Mock window.crypto.randomUUID
Object.defineProperty(global, 'crypto', {
  value: {
    randomUUID: () => {
      return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = (Math.random() * 16) | 0
        const v = c === 'x' ? r : (r & 0x3) | 0x8
        return v.toString(16)
      })
    }
  },
  writable: true
})

// Mock import.meta.client
Object.defineProperty(import.meta, 'client', {
  value: true,
  writable: true
})

beforeEach(() => {
  // Clear localStorage before each test
  localStorageMock.clear()
  
  // Clear IndexedDB store and object stores before each test
  Object.keys(indexedDBStore).forEach(key => delete indexedDBStore[key])
  objectStoreNamesSet.clear()
  
  // Reset all mocks
  vi.clearAllMocks()
  
  // Reset import.meta.client
  Object.defineProperty(import.meta, 'client', {
    value: true,
    writable: true
  })
})

// Make localStorage available globally (only in browser-like environments)
if (typeof window !== 'undefined') {
  Object.defineProperty(window, 'localStorage', {
    value: localStorageMock,
    writable: true
  })
  
  // Make indexedDB available globally
  Object.defineProperty(window, 'indexedDB', {
    value: indexedDBMock,
    writable: true
  })
}

// Also make indexedDB available on global for Node environment
Object.defineProperty(global, 'indexedDB', {
  value: indexedDBMock,
  writable: true
})

