# 🚀 Quick Start Guide

## Current Status

✅ Frontend is running on `http://localhost:5173`
✅ API client is fully configured
✅ UI is displaying a Books browser with filters and pagination

## What You're Seeing

The app is now showing a proper UI with:
- 📚 Books grid with cover images
- 🔍 Search functionality
- 🏷️ Type filter (Manga/Novel)
- 🌐 Source filter
- 📄 Pagination

## Next Steps

### 1. Start the Backend Server

The frontend is trying to connect to `http://localhost:8080`. You need to start your Go backend:

```bash
# From the project root
go run main.go
```

Or if using Docker:
```bash
cd deploy
docker-compose up
```

### 2. Expected Behavior

Once the backend is running:
- Books will load from your database
- Filters will work
- You can browse through pages

### 3. Current Error

If you're seeing an error in the browser console like:
- "Failed to fetch" or "Network error"
- "CORS error"

This means the backend isn't running yet on port 8080.

## Backend Configuration

If your backend runs on a different port, update the `.env` file:

```bash
# app/.env
VITE_API_BASE_URL=http://localhost:YOUR_PORT
```

Then restart the Vite dev server (Ctrl+C and `npm run dev` again).

## Features Available

### Current Page (`/`)
- Browse all books
- Search by title
- Filter by type (manga/novel)
- Filter by source
- Pagination

### API Integration
All these functions are available via `$lib/api`:
- `listBooks()` - Browse books
- `getBook(id)` - Get book details
- `listChapters()` - Get chapters
- `listFavorites()` - Get user favorites
- `createFavorite()` - Add to favorites
- `login()` / `register()` - Authentication

## Troubleshooting

### "No books found"
- Backend needs to have books in the database
- Check backend is running: `curl http://localhost:8080/_status`

### CORS Error
- Backend needs CORS configured to allow `localhost:5173`
- Or run with `--host` flag: `npm run dev -- --host`

### TypeScript Errors
- Run `npm run check` to see any type errors
- All should be clean ✓

## Creating New Pages

To add more pages, create new routes:

```bash
# Example: Book detail page
# Create: src/routes/books/[id]/+page.svelte

<script lang="ts">
  import { page } from '$app/stores';
  import { getBook } from '$lib/api';
  
  let bookId = $page.params.id;
  // Load book data...
</script>
```

## Environment Variables

Current configuration in `app/.env`:
```
VITE_API_BASE_URL=http://localhost:8080
```

This tells the frontend where to find your backend API.
