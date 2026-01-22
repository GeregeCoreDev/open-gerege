# API Documentation

This project uses a unified API client wrapper around `axios` (or `fetch`) and custom React hooks for data fetching and mutations.

## Core API Client

The core client is located in `lib/api.ts`. It handles:
- Base URL configuration
- Authentication token injection (if applicable)
- Unified error handling
- Response transformation

### Usage

```typescript
import { api } from '@/lib/api';

// GET request
const data = await api.get<User>('/users/me');

// POST request
const newUser = await api.post<User>('/users', { name: 'John' });
```

## React Hooks

We use custom hooks that wrap the API client to provide state management (loading, error, success) and caching.

### `useApiQuery`

Used for fetching data (GET requests). Supports caching and deduplication.

```typescript
import { useApiQuery } from '@/hooks/useApiQuery';

const { data, isLoading, error, refetch } = useApiQuery<User>({
  path: '/users/me',
  query: { include: 'roles' }, // Optional query params
  enabled: true, // Conditional fetching
  cacheTtl: 5000, // Time to live in ms
});
```

### `useApiMutation`

Used for modifying data (POST, PUT, DELETE).

```typescript
import { useApiMutation } from '@/hooks/useApiQuery';

const { mutate, isLoading } = useApiMutation<User, UpdateUserDto>({
  path: '/users/me',
  method: 'PATCH',
  onSuccess: (data) => console.log('Updated!', data),
  invalidateKeys: ['api:/users/me:'], // Invalidate cache after success
});

// Trigger the mutation
const handleSubmit = (formData) => {
  mutate(formData);
};
```

## Error Handling

Errors are standardized using the `logger` utility. API errors typically return an object with a uniform structure, but the hooks abstract this into an `error` object (instance of `Error`).
