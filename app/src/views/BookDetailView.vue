<template lang="pug">
.page
  .spinner(v-if="loading")
  p.error-msg(v-else-if="error") {{ error }}
  template(v-else-if="book")
    .detail-hero
      .detail-hero__cover
        img(
          v-if="book.image_url"
          :src="resolveImage(book.image_url, book.source?.domain)"
          :alt="book.title"
        )
        .detail-hero__no-cover(v-else) {{ book.title[0] }}
      .detail-hero__info
        .detail-hero__meta
          span.badge(:class="`badge--${book.type}`") {{ book.type }}
          span.detail-hero__source(v-if="book.source") {{ book.source.name }}
        h1.detail-hero__title {{ book.title }}
        p.detail-hero__synopsis(v-if="book.synopsis") {{ book.synopsis }}
        p.detail-hero__synopsis.detail-hero__synopsis--muted(v-else) No synopsis available.
        .detail-hero__actions
          button.btn.btn--ghost(
            @click="refreshBook"
            :disabled="refreshLoading"
            title="Refresh book data from source"
          ) {{ refreshLoading ? 'Refreshing...' : '↻ Refresh' }}
          button.btn.btn--primary(
            v-if="!favorited"
            @click="addFavorite"
            :disabled="favLoading"
          ) ♡ Add to Favorites
          button.btn.btn--danger(
            v-else
            @click="removeFavorite"
            :disabled="favLoading"
          ) ✕ Remove Favorite
          p.error-msg(v-if="favError") {{ favError }}

    .chapters-section
      .chapters-section__header
        h2 Chapters
        span.chapters-section__count(v-if="chapters.length") {{ chapters.length }} chapter{{ chapters.length !== 1 ? 's' : '' }}
      .spinner(v-if="loading")
      p.error-msg(v-else-if="!chapters.length") No chapters available.
      .chapters-list(v-else)
        router-link.chapter-item(
          v-for="ch in chapters"
          :key="ch.ID"
          :to="`/chapters/${ch.ID}`"
        )
          .chapter-item__left
            span.chapter-item__num(v-if="ch.number") Ch. {{ ch.number }}
            span.chapter-item__title {{ ch.title }}
          .chapter-item__right
            span.chapter-item__date(v-if="ch.upload_date") {{ formatDate(ch.upload_date) }}
            span.chapter-item__status(v-if="ch.completed") ✓
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { booksApi } from '../api/books'
import { favoritesApi } from '../api/favorites'
import type { Book, Favorite } from '../types'

const route = useRoute()

function resolveImage(imageUrl: string, domain?: string): string {
  if (imageUrl.startsWith('http')) return imageUrl
  return `https://${domain}${imageUrl}`
}
void resolveImage

const book = ref<Book | null>(null)
const loading = ref(false)
const error = ref('')
const favorited = ref<Favorite | null>(null)
const favLoading = ref(false)
const favError = ref('')
const refreshLoading = ref(false)

const chapters = computed(() =>
  [...(book.value?.chapters ?? [])].sort((a, b) => parseFloat(b.number ?? '0') - parseFloat(a.number ?? '0'))
)
void chapters

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString(undefined, {
    year: 'numeric', month: 'short', day: 'numeric'
  })
}
void formatDate

async function loadBook() {
  loading.value = true
  error.value = ''
  try {
    book.value = await booksApi.get(Number(route.params.id))
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load book'
  } finally {
    loading.value = false
  }
}

async function refreshBook() {
  refreshLoading.value = true
  error.value = ''
  try {
    book.value = await booksApi.get(Number(route.params.id), true)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to refresh book'
  } finally {
    refreshLoading.value = false
  }
}
void refreshBook

async function loadFavorites() {
  try {
    const res = await favoritesApi.list(1, 200)
    favorited.value =
      res.items.find(f => f.book_id === Number(route.params.id) || f.book?.ID === Number(route.params.id)) ?? null
  } catch {
    // non-critical
  }
}

async function addFavorite() {
  favLoading.value = true
  favError.value = ''
  try {
    const fav = await favoritesApi.create(Number(route.params.id))
    favorited.value = fav
  } catch (e: unknown) {
    favError.value = e instanceof Error ? e.message : 'Failed to add favorite'
  } finally {
    favLoading.value = false
  }
}
void addFavorite

async function removeFavorite() {
  if (!favorited.value) return
  favLoading.value = true
  try {
    await favoritesApi.delete(favorited.value.ID)
    favorited.value = null
  } catch {
    // ignore
  } finally {
    favLoading.value = false
  }
}
void removeFavorite

onMounted(() => {
  loadBook()
  loadFavorites()
})
</script>

<style lang="less" scoped>
.detail-hero {
  display: flex;
  gap: 2rem;
  margin-bottom: 3rem;
  align-items: flex-start;

  @media (max-width: 640px) {
    flex-direction: column;
  }

  &__cover {
    flex-shrink: 0;
    width: 220px;
    border-radius: @radius-lg;
    overflow: hidden;
    box-shadow: @shadow-md;
    background-color: @bg-secondary;
    aspect-ratio: 2 / 3;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    @media (max-width: 640px) {
      width: 100%;
      max-width: 220px;
      margin: 0 auto;
    }
  }

  &__no-cover {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 4rem;
    font-weight: 700;
    color: @text-muted;
    background-color: @bg-input;
  }

  &__info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  &__meta {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  &__source {
    font-size: @fs-sm;
    color: @text-secondary;
  }

  &__title {
    font-size: @fs-3xl;
    font-weight: 700;
    line-height: 1.25;
  }

  &__synopsis {
    font-size: @fs-base;
    color: @text-secondary;
    line-height: 1.75;

    &--muted {
      color: @text-muted;
      font-style: italic;
    }
  }

  &__actions {
    display: flex;
    gap: 0.75rem;
    margin-top: 0.5rem;
    flex-wrap: wrap;
  }
}

.chapters-section {
  &__header {
    display: flex;
    align-items: baseline;
    gap: 1rem;
    margin-bottom: 1rem;
    border-bottom: 1px solid @border-color;
    padding-bottom: 0.75rem;

    h2 {
      font-size: @fs-2xl;
      font-weight: 700;
    }
  }

  &__count {
    font-size: @fs-sm;
    color: @text-muted;
  }
}

.chapters-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.chapter-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid @border-color;
  text-decoration: none;
  color: @text-primary;
  transition: background-color @transition-fast;
  gap: 1rem;

  &:hover {
    background-color: @bg-card;
    text-decoration: none;
    color: @text-primary;
  }

  &__left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    overflow: hidden;
  }

  &__num {
    font-size: @fs-xs;
    font-weight: 700;
    color: @accent;
    white-space: nowrap;
    background-color: @accent-muted;
    padding: 0.1rem 0.5rem;
    border-radius: @radius-sm;
  }

  &__title {
    font-size: @fs-sm;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__right {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-shrink: 0;
  }

  &__date {
    font-size: @fs-xs;
    color: @text-muted;
  }

  &__status {
    font-size: @fs-xs;
    color: @success;
    font-weight: 700;
  }
}
</style>
