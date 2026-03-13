<template lang="pug">
.chapter-page
  .spinner(v-if="loading")
  p.error-msg(v-else-if="error") {{ error }}
  template(v-else-if="chapter")
    .chapter-nav(:class="{ 'chapter-nav--hidden': !showNav }")
      router-link.btn.btn--ghost(
        v-if="chapter.book_id"
        :to="`/books/${chapter.book_id}`"
      ) ← Back to Book
      h2.chapter-nav__title
        router-link.chapter-nav__crumb.chapter-nav__crumb--book(
          v-if="chapter.book_id && chapter.book?.title"
          :to="`/books/${chapter.book_id}`"
        ) {{ chapter.book.title }}
        span.chapter-nav__crumb-sep /&nbsp;
        span {{ chapter.number + '. ' + chapter.title }}
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

      //- Loading indicator for next chapter
      .next-chapter-loading(v-if="loadingNextChapter")
        .spinner
        p Loading next chapter...

      .chapter-end-sentinel(ref="chapterEndSentinel" aria-hidden="true")

    .chapter-bottom
      router-link.btn.btn--ghost(
        v-if="chapter.book_id"
        :to="`/books/${chapter.book_id}`"
      ) ← Back to Book
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { chaptersApi } from '../api/chapters'
import type { Chapter } from '../types'

const route = useRoute()
const router = useRouter()

const chapter = ref<Chapter | null>(null)
const nextChapter = ref<Chapter | null>(null)
const loading = ref(false)
const error = ref('')
const showNav = ref(true)
const loadingNextChapter = ref(false)
const chapters = ref<Chapter[]>([])
const chapterEndSentinel = ref<HTMLElement | null>(null)

let nextChapterInFlight = false
let endObserver: IntersectionObserver | null = null
let userHasScrolled = false

let lastScrollY = 0
let upScrollDistance = 0

function handleScroll() {
  const currentY = window.scrollY
  const delta = currentY - lastScrollY

  if (currentY > 120) {
    userHasScrolled = true
  }

  if (currentY <= 20) {
    showNav.value = true
    upScrollDistance = 0
    lastScrollY = currentY
    return
  }

  if (delta > 4) {
    // Scrolling down: hide quickly.
    showNav.value = false
    upScrollDistance = 0
  } else if (delta < -2) {
    // Scrolling up: reveal after a small upward movement.
    upScrollDistance += -delta
    if (upScrollDistance >= 40) {
      showNav.value = true
      upScrollDistance = 0
    }
  }

  lastScrollY = currentY
}

const isManga = computed(() => chapter.value?.book?.type === 'manga')
void isManga

async function fetchChaptersList() {
  if (!chapter.value?.book_id) return
  try {
    const response = await chaptersApi.list(chapter.value.book_id, 1, 1000)
    chapters.value = response.items
  } catch (e: unknown) {
    console.error('Failed to fetch chapters list:', e)
  }
}

function getNextChapter() {
  const current = chapter.value
  if (!current || chapters.value.length === 0) return null

  const ordered = [...chapters.value].sort((a, b) => {
    const aNum = Number.parseFloat(a.number ?? '')
    const bNum = Number.parseFloat(b.number ?? '')
    const aRank = Number.isFinite(aNum) ? aNum : Number.MAX_SAFE_INTEGER
    const bRank = Number.isFinite(bNum) ? bNum : Number.MAX_SAFE_INTEGER

    if (aRank !== bRank) return aRank - bRank
    return a.ID - b.ID
  })

  const currentIndex = ordered.findIndex(c => c.ID === current.ID)
  if (currentIndex >= 0) {
    for (let i = currentIndex + 1; i < ordered.length; i++) {
      if (ordered[i].ID !== current.ID) {
        return ordered[i]
      }
    }
  }

  return null
}

async function loadChapter() {
  loading.value = true
  error.value = ''
  nextChapterInFlight = false
  userHasScrolled = false
  try {
    chapter.value = await chaptersApi.get(Number(route.params.id))
    await fetchChaptersList()
    const next = getNextChapter()
    nextChapter.value = next
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Failed to load chapter'
  } finally {
    loading.value = false
  }
}

watch([chapter, loading], async ([currentChapter, isLoading]) => {
  if (!currentChapter || isLoading) return
  await nextTick()
  setupEndObserver()
})

async function loadNextChapter() {
  if (!chapter.value) return
  if (nextChapterInFlight) return
  if (loadingNextChapter.value) return

  const candidate = getNextChapter()
  nextChapter.value = candidate
  if (!candidate || candidate.ID === chapter.value.ID) {
    nextChapter.value = null
    return
  }

  nextChapterInFlight = true
  loadingNextChapter.value = true
  try {
    const next = await chaptersApi.get(candidate.ID)

    // Safety guard: never replace current content with itself.
    if (!chapter.value || next.ID === chapter.value.ID) {
      nextChapter.value = null
      return
    }

    chapter.value = next
    await router.push(`/chapters/${next.ID}`)

    // Fetch new chapters list and find next
    await fetchChaptersList()
    const newNext = getNextChapter()
    nextChapter.value = newNext

    // Scroll to top of new chapter
    window.scrollTo(0, 0)
    showNav.value = true
  } catch (e: unknown) {
    console.error('Failed to load next chapter:', e)
  } finally {
    nextChapterInFlight = false
    loadingNextChapter.value = false
  }
}

function setupEndObserver() {
  if (!chapterEndSentinel.value) return

  endObserver?.disconnect()
  endObserver = new IntersectionObserver(
    (entries) => {
      const hitEnd = entries.some((entry) => entry.isIntersecting)
      if (hitEnd && userHasScrolled) {
        void loadNextChapter()
      }
    },
    {
      root: null,
      rootMargin: '0px 0px 280px 0px',
      threshold: 0.01
    }
  )

  endObserver.observe(chapterEndSentinel.value)
}

onMounted(() => {
  lastScrollY = window.scrollY
  window.addEventListener('scroll', handleScroll, { passive: true })
  void loadChapter()
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
  endObserver?.disconnect()
})
</script>

<style lang="less" scoped>
.chapter-page {
  width: 100%;
  max-width: none;
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
  position: sticky;
  top: 0;
  z-index: 5;
  background: @bg-primary;
  transition: transform 0.25s ease;

  &--hidden {
    transform: translateY(calc(-100% - 1rem));
  }

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

  &__crumb {
    color: @text-secondary;
    text-decoration: none;

    &:hover {
      color: @accent;
      text-decoration: underline;
    }
  }

  &__crumb--book {
    margin-right: 0.35rem;
  }

  &__crumb--chapter {
    margin-left: 0.35rem;
    color: @text-primary;
  }

  &__crumb-sep {
    color: @text-muted;
    font-weight: 500;
  }
}

.chapter-body {
  margin-bottom: 2rem;
}

.next-chapter-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 3rem 1rem;
  text-align: center;
  color: @text-secondary;

  p {
    margin: 0;
    font-size: @fs-base;
  }
}

.chapter-end-sentinel {
  width: 100%;
  height: 1px;
}

// ── Manga reader ──────────────────────────────────────────────────
.manga-reader {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;
  width: 100%;

  &__page {
    width: 100%;
    height: auto;
    display: block;
  }
}

// ── Novel reader ──────────────────────────────────────────────────
.novel-reader {
  width: 100%;
  margin: 0;

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
