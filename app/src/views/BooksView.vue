<template lang="pug">
.page
  .books-header
    h1.books-header__title Books
    .books-header__filters
      .books-header__type-filters
        label.books-header__type-label
          input(type="checkbox" value="manga" v-model="selectedTypes")
          | Manga
        label.books-header__type-label
          input(type="checkbox" value="novel" v-model="selectedTypes")
          | Novel
      select(v-model="selectedSource")
        option(value="") All Sources
        option(v-for="src in sources" :key="src.ID" :value="src.ID") {{ src.name }}
      input.books-header__search(
        type="search"
        v-model="searchDraft"
        placeholder="Filter by title…"
        @keyup.enter="commitSearch"
        @search="onSearchClear"
      )

  .books-grid(v-if="books.length")
    router-link.book-card(
      v-for="book in books"
      :key="book.ID"
      :to="`/books/${book.ID}`"
    )
      .book-card__cover
        img(
          v-if="book.image_url"
          :src="resolveImage(book.image_url, book.source?.domain)"
          :alt="book.title"
          loading="lazy"
        )
        .book-card__no-cover(v-else) {{ book.title[0] }}
        img.book-card__source-icon(
          v-if="book.source?.icon_url"
          :src="resolveImage(book.source.icon_url, book.source.domain)"
          :alt="book.source.name"
          :title="book.source.name"
        )
      span.book-card__ribbon(:class="`badge badge--${book.type}`") {{ book.type }}
      button.book-card__fav(
          :class="{ 'book-card__fav--active': favMap[book.ID] }"
          @click.prevent.stop="toggleFavorite(book.ID)"
          :disabled="favPending === book.ID"
          :title="favMap[book.ID] ? 'Remove from favorites' : 'Add to favorites'"
        ) {{ favMap[book.ID] ? '♥' : '♡' }}
      .book-card__body
        h3.book-card__title {{ book.title }}

  p.empty-state(v-else-if="!loading && !error") No books found.
  p.error-msg(v-if="error") {{ error }}

  .load-sentinel(ref="sentinel")
  .spinner(v-if="loading")
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { booksApi } from '../api/books'
import { sourcesApi } from '../api/sources'
import { favoritesApi } from '../api/favorites'
import type { Book, Source } from '../types'

function resolveImage(imageUrl: string, domain?: string): string {
  if (imageUrl.startsWith('http')) return imageUrl
  return `https://${domain}${imageUrl}`
}

const books = ref<Book[]>([])
const sources = ref<Source[]>([])
const loading = ref(false)
const error = ref('')
const page = ref(0)
const totalPages = ref(1)
const selectedSource = ref<number | ''>('')
const selectedTypes = ref<string[]>(['manga', 'novel'])
const searchDraft = ref('')
const search = ref('')
const sentinel = ref<HTMLElement | null>(null)
const favMap = ref<Record<number, number>>({}) // bookId -> favoriteId
const favPending = ref<number | null>(null)
let io: IntersectionObserver | null = null

const hasMore = computed(() => page.value < totalPages.value)

async function loadNextPage() {
  console.log('[lazy] loadNextPage called — loading:', loading.value, 'hasMore:', hasMore.value, 'page:', page.value, 'totalPages:', totalPages.value)
  if (loading.value || !hasMore.value) {
    console.log('[lazy] loadNextPage SKIPPED')
    return
  }
  loading.value = true
  error.value = ''
  try {
    const sourceId = selectedSource.value ? Number(selectedSource.value) : undefined
    const bookType = selectedTypes.value.length === 1 ? selectedTypes.value[0] : undefined
    const res = await booksApi.list(page.value + 1, 20, sourceId, search.value || undefined, bookType)
    console.log('[lazy] API response pagination:', res.pagination)
    books.value.push(...res.items)
    page.value = res.pagination.page
    totalPages.value = res.pagination.total_pages
    console.log('[lazy] After load — page:', page.value, 'totalPages:', totalPages.value, 'hasMore:', hasMore.value)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load books'
  } finally {
    loading.value = false
  }
}

function observeSentinel() {
  io?.disconnect()
  io = null
  console.log('[lazy] observeSentinel — sentinel:', !!sentinel.value, 'hasMore:', hasMore.value)
  if (!sentinel.value || !hasMore.value) {
    console.log('[lazy] observeSentinel ABORTED — no sentinel or no more pages')
    return
  }
  io = new IntersectionObserver(
    async ([entry]) => {
      console.log('[lazy] IO fired — isIntersecting:', entry.isIntersecting, 'boundingRect:', entry.boundingClientRect)
      if (!entry.isIntersecting) return
      io?.disconnect()
      io = null
      await loadNextPage()
      observeSentinel()
    },
    { rootMargin: '0px 0px 400px 0px' }
  )
  io.observe(sentinel.value)
  console.log('[lazy] IO attached to sentinel')
}

function commitSearch() {
  search.value = searchDraft.value
  resetAndLoad()
}

function onSearchClear(e: Event) {
  if ((e.target as HTMLInputElement).value === '') {
    search.value = ''
    resetAndLoad()
  }
}

async function resetAndLoad() {
  io?.disconnect()
  io = null
  books.value = []
  page.value = 0
  totalPages.value = 1
  await loadNextPage()
  observeSentinel()
}

watch(selectedSource, resetAndLoad)
watch(selectedTypes, resetAndLoad)

onMounted(async () => {
  loadSources()
  loadFavorites()
  await loadNextPage()
  observeSentinel()
})

onBeforeUnmount(() => {
  io?.disconnect()
  io = null
})

async function loadSources() {
  try {
    const res = await sourcesApi.list(1, 100)
    sources.value = res.items
  } catch {
    // non-critical
  }
}

async function loadFavorites() {
  try {
    const res = await favoritesApi.list(1, 1000)
    const map: Record<number, number> = {}
    for (const fav of res.items) {
      if (fav.book_id) map[fav.book_id] = fav.ID
    }
    favMap.value = map
  } catch {
    // non-critical
  }
}

async function toggleFavorite(bookId: number) {
  favPending.value = bookId
  try {
    if (favMap.value[bookId]) {
      await favoritesApi.delete(favMap.value[bookId])
      const updated = { ...favMap.value }
      delete updated[bookId]
      favMap.value = updated
    } else {
      const fav = await favoritesApi.create(bookId)
      favMap.value = { ...favMap.value, [bookId]: fav.ID }
    }
  } catch {
    // ignore
  } finally {
    favPending.value = null
  }
}
</script>

<style lang="less" scoped>
.books-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 2rem;

  &__title {
    font-size: @fs-3xl;
    font-weight: 700;
  }

  &__type-filters {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }

  &__type-label {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: @fs-sm;
    color: @text-primary;
    cursor: pointer;
    user-select: none;

    input[type='checkbox'] {
      accent-color: @accent;
      cursor: pointer;
    }
  }

  &__filters {
    display: flex;
    gap: 0.75rem;
    flex-wrap: wrap;

    select {
      padding: 0.5rem 0.85rem;
      background-color: @bg-input;
      border: 1px solid @border-color;
      border-radius: @radius-md;
      color: @text-primary;
      font-size: @fs-sm;
      outline: none;
      cursor: pointer;

      &:focus { border-color: @accent; }
    }
  }

  &__search {
    padding: 0.5rem 0.85rem;
    background-color: @bg-input;
    border: 1px solid @border-color;
    border-radius: @radius-md;
    color: @text-primary;
    font-size: @fs-sm;
    outline: none;
    width: 200px;

    &:focus { border-color: @accent; }
    &::placeholder { color: @text-muted; }
  }
}

.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 1.25rem;
}

.book-card {
  position: relative;
  background-color: @bg-card;
  border: 1px solid @border-color;
  border-radius: @radius-lg;
  overflow: hidden;
  text-decoration: none;
  color: inherit;
  transition: transform @transition-base, box-shadow @transition-base;
  display: flex;
  flex-direction: column;

  &:hover {
    transform: translateY(-4px);
    box-shadow: @shadow-md;
    text-decoration: none;
    color: inherit;
  }

  &__cover {
    position: relative;
    aspect-ratio: 2 / 3;
    overflow: hidden;
    background-color: @bg-secondary;
    flex-shrink: 0;

    img:not(.book-card__source-icon) {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.4s ease;
    }

    &:hover img:not(.book-card__source-icon) {
      transform: scale(1.04);
    }
  }

  &__no-cover {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 3rem;
    font-weight: 700;
    color: @text-muted;
    background-color: @bg-input;
  }

  &__body {
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    flex: 1;
  }

  &__title {
    font-size: @fs-sm;
    font-weight: 600;
    color: @text-primary;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  &__source {
    font-size: @fs-xs;
    color: @text-muted;
    margin-top: auto;
  }

  &__source-icon {
    position: absolute;
    bottom: 0.4rem;
    right: 0.4rem;
    width: 1.75rem;
    height: 1.75rem;
    border-radius: @radius-sm;
    object-fit: contain;
    background: rgba(0, 0, 0, 0.55);
    padding: 0.15rem;
    z-index: 2;
  }

  &__ribbon {
    position: absolute;
    top: 0.5rem;
    left: -1.75rem;
    width: 6rem;
    text-align: center;
    transform: rotate(-45deg);
    padding: 0.2rem 0.2rem;
    font-size: @fs-xs;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    z-index: 2;
    pointer-events: none;
    background-color: #000 !important;
    color: #aaa !important;
    border-radius: 0%;
  }

  &__fav {
    position: absolute;
    top: 0.4rem;
    right: 0.4rem;
    background: rgba(0, 0, 0, 0.55);
    border: none;
    border-radius: 50%;
    width: 2.5rem;
    height: 2.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2rem;
    line-height: 1;
    color: @text-muted;
    cursor: pointer;
    transition: color 0.2s ease;
    color: #e53e3e;
    z-index: 1;

    &--active {
      opacity: 1;
    }

    &:disabled { cursor: default; }
  }

}


.load-sentinel {
  height: 1px;
  margin-top: 1rem;
}
</style>
