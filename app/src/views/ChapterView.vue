<template lang="pug">
.chapter-page
  .spinner(v-if="loading")
  p.error-msg(v-else-if="error") {{ error }}
  template(v-else-if="chapter")
    .chapter-nav
      router-link.btn.btn--ghost(
        v-if="chapter.book_id"
        :to="`/books/${chapter.book_id}`"
      ) ← Back to Book
      h2.chapter-nav__title {{ chapter.title }}
      .chapter-nav__spacer

    .chapter-body
      //- Manga: render each content item as an image
      template(v-if="isManga && chapter.content?.length")
        .manga-reader
          img.manga-reader__page(
            v-for="(src, idx) in chapter.content.split('::;;::')"
            :key="idx"
            :src="src"
            :alt="`Page ${idx + 1}`"
            loading="lazy"
          )

      //- Novel: render each content item as a paragraph
      template(v-else-if="!isManga && chapter.content?.length")
        .novel-reader
          p.novel-reader__para(
            v-for="(para, idx) in chapter.content.split('::;;::')"
            :key="idx"
          ) {{ para }}

      p.empty-state(v-else) No content available for this chapter.

    .chapter-bottom
      router-link.btn.btn--ghost(
        v-if="chapter.book_id"
        :to="`/books/${chapter.book_id}`"
      ) ← Back to Book
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { chaptersApi } from '../api/chapters'
import type { Chapter } from '../types'

const route = useRoute()

const chapter = ref<Chapter | null>(null)
const loading = ref(false)
const error = ref('')

const isManga = computed(() => chapter.value?.book?.type === 'manga')

async function loadChapter() {
  loading.value = true
  error.value = ''
  try {
    chapter.value = await chaptersApi.get(Number(route.params.id))
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load chapter'
  } finally {
    loading.value = false
  }
}

onMounted(loadChapter)
</script>

<style lang="less" scoped>
.chapter-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 1.5rem 1rem;
}

.chapter-nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid @border-color;

  &__title {
    font-size: @fs-lg;
    font-weight: 600;
    flex: 1;
    text-align: center;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__spacer {
    width: 120px; // mirror the back button width
  }
}

.chapter-body {
  margin-bottom: 2rem;
}

// ── Manga reader ──────────────────────────────────────────────────
.manga-reader {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;

  &__page {
    width: 100%;
    max-width: 800px;
    height: auto;
    display: block;
  }
}

// ── Novel reader ──────────────────────────────────────────────────
.novel-reader {
  max-width: 720px;
  margin: 0 auto;

  &__para {
    font-size: @fs-lg;
    line-height: 1;
    color: @text-primary;
    margin-bottom: 1.25em;
  }
}

.chapter-bottom {
  display: flex;
  justify-content: flex-start;
  padding-top: 1rem;
  border-top: 1px solid @border-color;
}
</style>
