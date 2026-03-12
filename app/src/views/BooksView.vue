<template lang="pug">
.page
  .books-header
    h1.books-header__title Books
    .books-header__filters
      select(v-model="selectedSource")
        option(value="") All Sources
        option(v-for="src in sources" :key="src.ID" :value="src.ID") {{ src.name }}
      input.books-header__search(
        type="search"
        v-model="search"
        placeholder="Filter by title…"
      )

  .books-grid(v-if="books.length")
    router-link.book-card(
      v-for="book in filtered"
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
      .book-card__body
        span.badge(:class="`badge--${book.type}`") {{ book.type }}
        h3.book-card__title {{ book.title }}
        p.book-card__source(v-if="book.source") {{ book.source.name }}

  p.empty-state(v-else-if="!loading && !error") No books found.
  p.error-msg(v-if="error") {{ error }}

  .load-sentinel(ref="sentinel")
  .spinner(v-if="loading")
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { booksApi } from '../api/books'
import { sourcesApi } from '../api/sources'
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
const search = ref('')
const sentinel = ref<HTMLElement | null>(null)
let io: IntersectionObserver | null = null

const filtered = computed(() =>
  search.value
    ? books.value.filter(b => b.title.toLowerCase().includes(search.value.toLowerCase()))
    : books.value
)

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
    const res = await booksApi.list(page.value + 1, 20, sourceId)
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

onMounted(async () => {
  loadSources()
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
    aspect-ratio: 2 / 3;
    overflow: hidden;
    background-color: @bg-secondary;
    flex-shrink: 0;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.4s ease;
    }

    &:hover img {
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
}

.load-sentinel {
  height: 1px;
  margin-top: 1rem;
}
</style>
