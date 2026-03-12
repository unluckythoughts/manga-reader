<template lang="pug">
.page
  .page-header
    h1 Sources
  .spinner(v-if="loading")
  p.error-msg(v-else-if="error") {{ error }}
  template(v-else)
    .sources-grid(v-if="sources.length")
      router-link.source-card(
        v-for="src in sources"
        :key="src.ID"
        :to="`/books?source_id=${src.ID}`"
      )
        .source-card__icon
          img(
            v-if="src.icon_url"
            :src="resolveImage(src.icon_url, src.domain)"
            :alt="src.name"
            loading="lazy"
          )
          span.source-card__icon-fallback(v-else) {{ src.name[0] }}
        .source-card__body
          h3.source-card__name {{ src.name }}
          a.source-card__domain(
            :href="`https://${src.domain}`"
            target="_blank"
            rel="noopener noreferrer"
            @click.stop
          ) {{ src.domain }}
    p.empty-state(v-else) No sources found.

    .pagination(v-if="totalPages > 1")
      button.pagination__btn(:disabled="page === 1" @click="loadSources(page - 1)") ‹ Prev
      button.pagination__btn(
        v-for="p in visiblePages" :key="p"
        :class="{ 'pagination__btn--active': p === page }"
        @click="loadSources(p)"
      ) {{ p }}
      button.pagination__btn(:disabled="page === totalPages" @click="loadSources(page + 1)") Next ›
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { sourcesApi } from '../api/sources'
import type { Source } from '../types'

function resolveImage(imageUrl: string, domain?: string): string {
  if (imageUrl.startsWith('http')) return imageUrl
  return `https://${domain}${imageUrl}`
}

const sources = ref<Source[]>([])
const loading = ref(false)
const error = ref('')
const page = ref(1)
const totalPages = ref(1)

const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

async function loadSources(p: number) {
  loading.value = true
  error.value = ''
  try {
    const res = await sourcesApi.list(p, 20)
    sources.value = res.items
    page.value = res.pagination.page
    totalPages.value = res.pagination.total_pages
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load sources'
  } finally {
    loading.value = false
  }
}

onMounted(() => loadSources(1))
</script>

<style lang="less" scoped>
.page-header {
  margin-bottom: 2rem;

  h1 {
    font-size: @fs-3xl;
    font-weight: 700;
  }
}

.sources-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1rem;
}

.source-card {
  background-color: @bg-card;
  border: 1px solid @border-color;
  border-radius: @radius-lg;
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  text-decoration: none;
  color: @text-primary;
  transition: transform @transition-base, box-shadow @transition-base;

  &:hover {
    transform: translateY(-3px);
    box-shadow: @shadow-md;
    text-decoration: none;
    color: @text-primary;
  }

  &__icon {
    width: 48px;
    height: 48px;
    border-radius: @radius-md;
    overflow: hidden;
    background-color: @bg-secondary;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;

    img {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
  }

  &__icon-fallback {
    font-size: 1.5rem;
    font-weight: 700;
    color: @text-muted;
    text-transform: uppercase;
  }

  &__body {
    flex: 1;
    overflow: hidden;
  }

  &__name {
    font-size: @fs-base;
    font-weight: 600;
    margin-bottom: 0.2rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__domain {
    font-size: @fs-xs;
    color: @text-muted;
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;

    &:hover {
      color: @accent;
      text-decoration: underline;
    }
  }
}
</style>
