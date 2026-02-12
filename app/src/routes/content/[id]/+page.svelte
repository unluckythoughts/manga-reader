<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getChapter, listChapters, updateFavoriteProgress, listFavorites, ApiClientError } from '$lib/api';
	import {isAuthenticated } from '$lib/stores';
	import type { Chapter, Favorite } from '$lib/types';

	let chapterId = $derived(parseInt($page.params.id || '0'));
	let currentChapter = $state<Chapter | null>(null);
	let nextChapter = $state<Chapter | null>(null);
	let loadingNext = $state(false);
	let allChapters = $state<Chapter[]>([]);
	let favorite = $state<Favorite | null>(null);
	
	let loading = $state(true);
	let error = $state<string | null>(null);

	async function loadChapter(id: number) {
		try {
			const chapter = await getChapter(id);
			return chapter;
		} catch (err) {
			if (err instanceof ApiClientError) {
				throw new Error(`Error ${err.status}: ${err.message}`);
			}
			throw new Error('Failed to load chapter');
		}
	}

	async function loadAllChapters(bookId: number) {
		try {
			// Load all chapters for this book to find next/previous
			const response = await listChapters({ book_id: bookId, limit: 1000 });
			allChapters = response.items;
		} catch (err) {
			console.error('Failed to load chapters list:', err);
		}
	}

	async function loadFavoriteStatus(bookId: number) {
		if (!$isAuthenticated) return;
		
		try {
			const response = await listFavorites({ limit: 1000 });
			favorite = response.items.find(f => f.book_id === bookId) || null;
		} catch (err) {
			console.error('Failed to load favorite status:', err);
		}
	}

	async function updateProgress(chapterNumber: string) {
		if (!favorite) return;

		try {
			await updateFavoriteProgress(favorite.id, {
				chapter: parseInt(chapterNumber) || 1
			});
		} catch (err) {
			console.error('Failed to update progress:', err);
		}
	}

	async function loadNextChapter() {
		if (!currentChapter || !allChapters.length) return;

		const currentIndex = allChapters.findIndex(c => c.id === currentChapter!.id);
		if (currentIndex === -1 || currentIndex >= allChapters.length - 1) return;

		const nextChapterInfo = allChapters[currentIndex + 1];
		
		loadingNext = true;
		try {
			nextChapter = await loadChapter(nextChapterInfo.id);
		} catch (err) {
			console.error('Failed to preload next chapter:', err);
		} finally {
			loadingNext = false;
		}
	}

	function getPreviousChapter(): Chapter | null {
		if (!currentChapter || !allChapters.length) return null;
		const currentIndex = allChapters.findIndex(c => c.id === currentChapter!.id);
		return currentIndex > 0 ? allChapters[currentIndex - 1] : null;
	}

	function getNextChapterInfo(): Chapter | null {
		if (!currentChapter || !allChapters.length) return null;
		const currentIndex = allChapters.findIndex(c => c.id === currentChapter!.id);
		return currentIndex < allChapters.length - 1 ? allChapters[currentIndex + 1] : null;
	}

	// Intersection Observer for lazy loading next chapter
	let nextChapterSentinel = $state<HTMLElement>();
	
	onMount(() => {
		// Async initialization
		(async () => {
			try {
				loading = true;
				error = null;

				currentChapter = await loadChapter(chapterId);
				
				if (currentChapter.book_id) {
					await Promise.all([
						loadAllChapters(currentChapter.book_id),
						loadFavoriteStatus(currentChapter.book_id)
					]);

					// Update reading progress
					if (currentChapter.number) {
						updateProgress(currentChapter.number);
					}

					// Preload next chapter
					setTimeout(() => loadNextChapter(), 1000);
				}
			} catch (err: any) {
				error = err.message || 'Failed to load chapter';
			} finally {
				loading = false;
			}
		})();

		// Setup intersection observer for lazy loading next chapter
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && !loadingNext && nextChapter) {
					// Append next chapter content
					if (currentChapter && nextChapter) {
						currentChapter = {
							...nextChapter
						};
						
						// Clear and load the next-next chapter
						nextChapter = null;
						setTimeout(() => loadNextChapter(), 500);
					}
				}
			},
			{ threshold: 0.5, rootMargin: '100px' }
		);

		if (nextChapterSentinel) {
			observer.observe(nextChapterSentinel);
		}

		return () => observer.disconnect();
	});

	let prevChapter = $derived(getPreviousChapter());
	let nextChapterInfo = $derived(getNextChapterInfo());
</script>

<svelte:head>
	<title>{currentChapter?.title || 'Loading...'} - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if loading && !currentChapter}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600">Loading chapter...</p>
		</div>
	{:else if error}
		<div class="max-w-4xl mx-auto px-4 py-12">
			<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<p class="text-red-800 font-medium">{error}</p>
				<button onclick={() => window.history.back()} class="mt-4 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition">
					Go Back
				</button>
			</div>
		</div>
	{:else if currentChapter}
		<!-- Chapter Header -->
		<div class="bg-white border-b sticky top-0 z-10 shadow-sm">
			<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-4">
						{#if currentChapter.book}
							<a
								href="/book/{currentChapter.book.id}"
								class="text-blue-600 hover:text-blue-700 font-medium"
							>
								← Back to Book
							</a>
						{/if}
						<div>
							<h1 class="text-lg font-bold text-gray-900">{currentChapter.title}</h1>
							{#if currentChapter.number}
								<p class="text-sm text-gray-500">Chapter {currentChapter.number}</p>
							{/if}
						</div>
					</div>

					<!-- Navigation Buttons -->
					<div class="flex items-center gap-2">
						{#if prevChapter}
							<a
								href="/content/{prevChapter.id}"
								class="px-4 py-2 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 transition font-medium"
							>
								← Previous
							</a>
						{/if}
						{#if nextChapterInfo}
							<a
								href="/content/{nextChapterInfo.id}"
								class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition font-medium"
							>
								Next →
							</a>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Chapter Content -->
		<div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<div class="bg-white rounded-lg shadow-lg p-8">
				<!-- Current Chapter Content -->
				{#if currentChapter.content && currentChapter.content.length > 0}
					<div class="prose prose-lg max-w-none">
						{#if currentChapter.type === 'manga' || currentChapter.book?.type === 'manga'}
							<!-- Manga: Display images -->
							<div class="space-y-2">
								{#each currentChapter.content as imageUrl}
									<img
										src={imageUrl}
										alt="Page"
										class="w-full h-auto"
										loading="lazy"
									/>
								{/each}
							</div>
						{:else}
							<!-- Novel: Display text -->
							<div class="space-y-4">
								{#each currentChapter.content as paragraph}
									<p class="text-gray-800 leading-relaxed">{paragraph}</p>
								{/each}
							</div>
						{/if}
					</div>
				{:else}
					<div class="text-center py-12">
						<div class="text-6xl mb-4">📄</div>
						<h3 class="text-xl font-semibold text-gray-900 mb-2">Content not available</h3>
						<p class="text-gray-600">This chapter's content hasn't been loaded yet</p>
					</div>
				{/if}

				<!-- Chapter Navigation Footer -->
				<div class="mt-12 pt-8 border-t flex items-center justify-between">
					{#if prevChapter}
						<a
							href="/content/{prevChapter.id}"
							class="flex items-center gap-2 px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition font-medium"
						>
							← Previous Chapter
						</a>
					{:else}
						<div></div>
					{/if}

					{#if currentChapter.book}
						<a
							href="/book/{currentChapter.book.id}"
							class="px-6 py-3 border-2 border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition font-medium"
						>
							All Chapters
						</a>
					{/if}

					{#if nextChapterInfo}
						<a
							href="/content/{nextChapterInfo.id}"
							class="flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition font-medium"
						>
							Next Chapter →
						</a>
					{:else}
						<div></div>
					{/if}
				</div>

				<!-- Lazy Load Next Chapter Sentinel -->
				{#if nextChapterInfo}
					<div bind:this={nextChapterSentinel} class="h-20 mt-8"></div>
					
					{#if loadingNext}
						<div class="text-center py-8">
							<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
							<p class="mt-2 text-gray-600">Loading next chapter...</p>
						</div>
					{/if}
				{/if}

				<!-- Next Chapter Preview (if loaded) -->
				{#if nextChapter}
					<div class="mt-12 pt-8 border-t">
						<div class="bg-blue-50 border-2 border-blue-200 rounded-lg p-6">
							<h3 class="text-lg font-semibold text-blue-900 mb-2">
								Next: {nextChapter.title}
							</h3>
							<p class="text-blue-700 text-sm mb-4">
								Chapter {nextChapter.number || ''}
							</p>
							<p class="text-blue-600 text-sm">
								Scroll down to automatically continue to the next chapter...
							</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
