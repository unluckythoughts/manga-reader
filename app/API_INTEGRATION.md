# Book Reader - Frontend App

This is the frontend application for the Book Reader, built with SvelteKit and TypeScript.

## Project Structure

```
src/
├── lib/
│   ├── api/              # API client for backend integration
│   │   ├── client.ts     # Base fetch wrapper
│   │   ├── auth.ts       # Authentication endpoints
│   │   ├── books.ts      # Books endpoints
│   │   ├── chapters.ts   # Chapters endpoints
│   │   ├── sources.ts    # Sources endpoints
│   │   ├── categories.ts # Categories endpoints
│   │   ├── favorites.ts  # Favorites endpoints
│   │   └── index.ts      # Main export
│   ├── components/       # Reusable Svelte components
│   ├── types.ts          # TypeScript type definitions
│   └── index.ts          # Library exports
├── routes/               # SvelteKit routes
│   ├── +layout.svelte
│   └── +page.svelte
└── app.html              # HTML template
```

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file (copy from `.env.example`):
   ```bash
   cp .env.example .env
   ```

3. Configure the API base URL in `.env`:
   ```
   VITE_API_BASE_URL=http://localhost:8080
   ```

## Development

Start the development server:

```bash
npm run dev
```

The app will be available at `http://localhost:5173`

## API Integration

The app includes a fully typed API client that matches the backend Go routes. All API functions are available through the `$lib/api` import.

### Quick Start

```typescript
import { listBooks, getBook, createFavorite } from '$lib/api';
import type { Book, PaginatedResponse } from '$lib/api';

// List books
const response: PaginatedResponse<Book> = await listBooks({
  page: 1,
  limit: 20,
  type: 'manga'
});

// Get a specific book
const book = await getBook(123);

// Add to favorites
await createFavorite({
  user_id: 1,
  book_id: book.id
});
```

### Available Endpoints

#### Authentication (if enabled)
- `login(data)` - Login with credentials
- `register(data)` - Register new user
- `logout()` - Logout current user
- `getUser()` - Get current user info
- `updateUser(data)` - Update user profile
- `changePassword(data)` - Change password
- `resetPassword(email)` - Request password reset
- `googleOAuthLogin(token)` - Login with Google

#### Books
- `listBooks(params?)` - List books with pagination
- `getBook(id)` - Get book details

#### Chapters
- `listChapters(params?)` - List chapters
- `getChapter(id)` - Get chapter content

#### Sources
- `listSources(params?)` - List book sources
- `getSource(id)` - Get source details

#### Categories
- `listCategories(params?)` - List categories
- `getCategory(id)` - Get category
- `createCategory(data)` - Create category
- `updateCategory(id, data)` - Update category
- `deleteCategory(id)` - Delete category

#### Favorites
- `listFavorites(params?)` - List user favorites
- `getFavorite(id)` - Get favorite
- `createFavorite(data)` - Add to favorites
- `updateFavorite(id, data)` - Update favorite
- `updateFavoriteProgress(id, data)` - Update reading progress
- `deleteFavorite(id)` - Remove from favorites

### Error Handling

All API calls throw `ApiClientError` on failure:

```typescript
import { ApiClientError } from '$lib/api';

try {
  const books = await listBooks();
} catch (error) {
  if (error instanceof ApiClientError) {
    if (error.status === 401) {
      // Redirect to login
    } else if (error.status === 404) {
      // Handle not found
    }
    console.error(error.message);
  }
}
```

### Example Component

See [BooksExample.svelte](src/lib/components/BooksExample.svelte) for a complete example of using the API client in a Svelte component with:
- Loading states
- Error handling
- Pagination
- TypeScript types

## Building

Build the production version:

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

## Type Checking

Run TypeScript type checking:

```bash
npm run check
```

Watch mode:

```bash
npm run check:watch
```

## Linting & Formatting

Check code style:

```bash
npm run lint
```

Format code:

```bash
npm run format
```

## Tech Stack

- **SvelteKit** - Full-stack framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **ESLint & Prettier** - Code quality

## API Client Features

- ✅ Fully typed TypeScript interfaces matching Go models
- ✅ Automatic JSON serialization/deserialization
- ✅ Cookie-based authentication support
- ✅ Comprehensive error handling
- ✅ Query parameter handling
- ✅ Pagination support
- ✅ RESTful API design

## Notes

- The API client automatically includes credentials (cookies) for authentication
- All timestamps are in ISO 8601 format
- IDs are numeric
- Pagination uses page-based navigation (page/limit)
