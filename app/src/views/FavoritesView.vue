<template lang="pug">
.page
  .page-header
    h1 My Favorites
  .spinner(v-if="loading")
  p.error-msg(v-else-if="error") {{ error }}
  template(v-else)
    .favorites-grid(v-if="favorites.length")
      .fav-card(v-for="fav in favorites" :key="fav.ID")
        router-link.fav-card__cover(:to="fav.book_id ? `/books/${fav.book?.ID ?? fav.book_id}` : '#'")
          img(
            v-if="fav.book?.image_url"
            :src="resolveImage(fav.book.image_url, fav.book.source?.domain)"
            :alt="fav.book?.title"
            loading="lazy"
          )
          .fav-card__no-cover(v-else) {{ fav.book?.title?.[0] ?? '?' }}
        .fav-card__body
          .fav-card__meta
            span.badge(v-if="fav.book?.type" :class="`badge--${fav.book.type}`") {{ fav.book.type }}
          router-link.fav-card__title(:to="fav.book?.ID ? `/books/${fav.book.ID}` : '#'")
            | {{ fav.book?.title ?? 'Unknown Book' }}
          p.fav-card__source(v-if="fav.book?.source") {{ fav.book.source.name }}
          .fav-card__progress(v-if="fav.progress?.length")
            span.fav-card__progress-label Progress:
            span {{ fav.progress.join(', ') }}
          .fav-card__actions
            button.btn.btn--danger(@click="removeFavorite(fav.ID)" :disabled="removing === fav.ID")
              span(v-if="removing === fav.ID") Removing…
              span(v-else) Remove
    p.empty-state(v-else) No favorites yet. Browse #[router-link(to="/books") books] to add some!

    .pagination(v-if="totalPages > 1")
      button.pagination__btn(:disabled="page === 1" @click="loadFavorites(page - 1)") ‹ Prev
      button.pagination__btn(
        v-for="p in visiblePages" :key="p"
        :class="{ 'pagination__btn--active': p === page }"
        @click="loadFavorites(p)"
      ) {{ p }}
      button.pagination__btn(:disabled="page === totalPages" @click="loadFavorites(page + 1)") Next ›
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { favoritesApi } from '../api/favorites'
import type { Favorite } from '../types'

function resolveImage(imageUrl: string, domain?: string): string {
  if (imageUrl.startsWith('http')) return imageUrl
  return `https://${domain}${imageUrl}`
}

const favorites = ref<Favorite[]>([])
const loading = ref(false)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)
const removing = ref<number | null>(null)

const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function loadFavorites(p: number) {
  loading.value = true
  error.value = ''
  try {
    const res = await favoritesApi.list(p, 20)
    favorites.value = res.items
    page.value = res.pagination.page
    totalPages.value = res.pagination.total_pages
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load favorites'
  } finally {
    loading.value = false
  }
}

async function removeFavorite(id: number) {
  removing.value = id
  try {
    await favoritesApi.delete(id)
    favorites.value = favorites.value.filter(f => f.ID !== id)
  } catch {
    // ignore
  } finally {
    removing.value = null
  }
}

onMounted(() => loadFavorites(1))
</script>

<style lang="less" scoped>
.page-header {
  margin-bottom: 2rem;

  h1 {
    font-size: @fs-3xl;
    font-weight: 700;
  }
}

.favorites-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1.25rem;
}

.fav-card {
  background-color: @bg-card;
  border: 1px solid @border-color;
  border-radius: @radius-lg;
  overflow: hidden;
  display: flex;
  flex-direction: column;

  &__cover {
    aspect-ratio: 2 / 3;
    overflow: hidden;
    background-color: @bg-secondary;
    flex-shrink: 0;
    display: block;

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
    padding: 0.85rem;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    flex: 1;
  }

  &__meta {
    display: flex;
    gap: 0.4rem;
  }

  &__title {
    font-size: @fs-sm;
    font-weight: 600;
    color: @text-primary;
    text-decoration: none;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;

    &:hover { color: @accent; text-decoration: none; }
  }

  &__source {
    font-size: @fs-xs;
    color: @text-muted;
  }

  &__progress {
    font-size: @fs-xs;
    color: @text-secondary;
    display: flex;
    gap: 0.35rem;
    flex-wrap: wrap;
  }

  &__progress-label {
    font-weight: 600;
    color: @text-muted;
  }

  &__actions {
    margin-top: auto;
    padding-top: 0.5rem;
  }
}
</style>
