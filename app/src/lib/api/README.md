# API Client

This directory contains the TypeScript API client for the Book Reader application, automatically generated from the Go backend routes.

## Structure

- `client.ts` - Base API client with fetch wrapper and error handling
- `auth.ts` - Authentication and user management endpoints
- `books.ts` - Book management endpoints
- `chapters.ts` - Chapter management endpoints
- `sources.ts` - Source management endpoints
- `categories.ts` - Category management endpoints
- `favorites.ts` - Favorites management endpoints
- `index.ts` - Exports all API functions and types

## Usage

### Basic Import

```typescript
import { listBooks, getBook, login, createFavorite } from '$lib/api';
```

### Examples

#### Authentication

```typescript
// Login
const { user, token } = await login({
  username: 'user@example.com',
  password: 'password123'
});

// Register
const { user, token } = await register({
  username: 'newuser',
  email: 'user@example.com',
  password: 'password123'
});

// Get current user
const user = await getUser();

// Logout
await logout();
```

#### Books

```typescript
// List books with pagination
const response = await listBooks({
  page: 1,
  limit: 20,
  type: 'manga',
  source_id: 1,
  search: 'one piece'
});

// Get a specific book
const book = await getBook(123);
```

#### Chapters

```typescript
// List chapters for a book
const chapters = await listChapters({
  book_id: 123,
  page: 1,
  limit: 50
});

// Get chapter content
const chapter = await getChapter(456);
```

#### Favorites

```typescript
// List user's favorites
const favorites = await listFavorites({
  category: 'reading',
  page: 1,
  limit: 20
});

// Add to favorites
const favorite = await createFavorite({
  user_id: 1,
  book_id: 123
});

// Update reading progress
await updateFavoriteProgress(favorite.id, {
  chapter: 10,
  level: 5
});

// Remove from favorites
await deleteFavorite(favorite.id);
```

#### Categories

```typescript
// List categories
const categories = await listCategories();

// Create category
const category = await createCategory({ name: 'Reading' });

// Update category
await updateCategory(category.id, { name: 'Currently Reading' });

// Delete category
await deleteCategory(category.id);
```

## Error Handling

All API functions throw `ApiClientError` on failure:

```typescript
import { ApiClientError } from '$lib/api';

try {
  const book = await getBook(123);
} catch (error) {
  if (error instanceof ApiClientError) {
    console.error(`API Error ${error.status}: ${error.error}`);
    // Handle specific status codes
    if (error.status === 404) {
      // Book not found
    }
  }
}
```

## Configuration

Set the API base URL in your `.env` file:

```
VITE_API_BASE_URL=http://localhost:8080
```

If not set, it defaults to the same origin as the app.

## Authentication

The API client automatically includes credentials (cookies) with all requests. When authentication is enabled on the backend, the auth middleware will protect the `/api/v1/reader/*` routes.

## Types

All TypeScript types are exported from `$lib/types` and re-exported from `$lib/api`:

```typescript
import type { Book, Chapter, Favorite, PaginatedResponse } from '$lib/api';
```
