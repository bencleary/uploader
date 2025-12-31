# Tests

This directory contains all tests for the `apps/web` application.

## Structure

```
__tests__/
├── setup.ts              # Global test setup and mocks
├── composables/          # Tests for composables
│   ├── useAuth.test.ts
│   ├── useChat.test.ts
│   ├── useUploadApi.test.ts
│   └── useUploads.test.ts
└── components/           # Tests for Vue components
    └── ChatInput.test.ts
```

## Running Tests

```bash
# Run all tests
pnpm test

# Run tests in watch mode
pnpm test --watch

# Run tests with UI
pnpm test:ui

# Run tests with coverage
pnpm test:coverage
```

## Test Setup

The `setup.ts` file configures:
- localStorage mocking
- crypto.randomUUID mocking
- import.meta.client mocking
- Global test utilities

## Writing Tests

### Testing Composables

Composables are tested by importing them directly and testing their behavior:

```typescript
import { useAuth } from '~/composables/useAuth'

describe('useAuth', () => {
  it('should do something', () => {
    const auth = useAuth()
    // test implementation
  })
})
```

### Testing Components

Components are tested using `@vue/test-utils`:

```typescript
import { mount } from '@vue/test-utils'
import MyComponent from '~/components/MyComponent.vue'

describe('MyComponent', () => {
  it('should render', () => {
    const wrapper = mount(MyComponent, {
      props: { /* props */ }
    })
    // assertions
  })
})
```

## Notes

- All Nuxt auto-imports (useState, computed, readonly, etc.) are mocked in individual test files
- localStorage is automatically cleared between tests
- import.meta.client is set to `true` by default in tests

