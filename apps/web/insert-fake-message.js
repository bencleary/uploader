// DevTools script to insert a fake reply message from a fictitious user
// Paste this entire script into the browser console and run it
// 
// Usage:
//   insertFakeMessage() // Uses default message
//   insertFakeMessage('Custom message text here') // Uses custom message

async function insertFakeMessage(messageText = 'This is a reply from a fictitious user! How can I help you today?') {
  const STORAGE_KEY = 'uploader_chat_messages';
  const DB_NAME = 'uploader_db';
  const DB_VERSION = 1;
  const STORE_NAME = 'storage';

  // Helper function to get IndexedDB database
  function getDatabase() {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION);
      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result);
      request.onupgradeneeded = (event) => {
        const db = event.target.result;
        if (!db.objectStoreNames.contains(STORE_NAME)) {
          db.createObjectStore(STORE_NAME);
        }
      };
    });
  }

  // Helper function to get messages from IndexedDB
  async function getMessages() {
    const db = await getDatabase();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(STORE_NAME, 'readonly');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.get(STORAGE_KEY);
      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve(request.result || []);
    });
  }

  // Helper function to save messages to IndexedDB
  async function saveMessages(messages) {
    const db = await getDatabase();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(STORE_NAME, 'readwrite');
      const store = transaction.objectStore(STORE_NAME);
      // Convert to plain objects (remove Vue reactivity proxies if any)
      const plainMessages = JSON.parse(JSON.stringify(messages));
      const request = store.put(plainMessages, STORAGE_KEY);
      request.onerror = () => reject(request.error);
      request.onsuccess = () => resolve();
    });
  }

  // Create a fake message from another user (right-aligned)
  const fakeMessage = {
    id: crypto.randomUUID(),
    text: messageText,
    timestamp: new Date().toISOString(), // Store as ISO string for IndexedDB
    attachments: undefined,
    isFromUser: false // This is from another user, so it will be right-aligned
  };

  try {
    // Get current messages
    const currentMessages = await getMessages();
    
    // Add the fake message
    const updatedMessages = [...currentMessages, fakeMessage];
    
    // Save to IndexedDB
    await saveMessages(updatedMessages);
    console.log('✅ Message saved to IndexedDB:', fakeMessage);
    
    // Try to update the Nuxt state by finding the Vue app instance
    // In Nuxt 3, we can try to access the app via the DOM
    let stateUpdated = false;
    
    // Try to access the Nuxt app instance through various methods
    const app = window.__NUXT__?.app || window.$nuxt || document.querySelector('#__nuxt')?.__vue_app__;
    
    if (app) {
      try {
        // Try to access the useState for chat_messages
        // Nuxt 3 stores state in a specific way
        const nuxtState = app.$state || app.state || app._instance?.setupState;
        
        if (nuxtState) {
          // Look for the chat_messages state
          // In Nuxt 3, useState creates a reactive ref
          const messagesState = nuxtState.chat_messages;
          
          if (messagesState && Array.isArray(messagesState)) {
            // Convert timestamp back to Date for the state
            const messageWithDate = {
              ...fakeMessage,
              timestamp: new Date(fakeMessage.timestamp)
            };
            messagesState.push(messageWithDate);
            stateUpdated = true;
            console.log('✅ State updated directly! Message should appear immediately.');
          } else if (messagesState && messagesState.value) {
            // If it's a ref, update the value
            const messageWithDate = {
              ...fakeMessage,
              timestamp: new Date(fakeMessage.timestamp)
            };
            messagesState.value.push(messageWithDate);
            stateUpdated = true;
            console.log('✅ State updated directly! Message should appear immediately.');
          }
        }
      } catch (e) {
        console.log('Could not update state directly:', e);
      }
    }
    
    if (!stateUpdated) {
      // If we couldn't update the state, offer to reload
      console.log('⚠️ Could not update state directly.');
      console.log('💡 The message has been saved to IndexedDB.');
      console.log('💡 Reload the page (F5 or Cmd+R) to see the new message.');
      console.log('💡 Or navigate away and back to trigger a state reload.');
    }
    
    return fakeMessage;
  } catch (error) {
    console.error('❌ Error inserting fake message:', error);
    throw error;
  }
}

// Make it available globally and auto-run with default message
// You can call insertFakeMessage('custom text') to use a custom message
window.insertFakeMessage = insertFakeMessage;

// Auto-run with default message
insertFakeMessage();

