<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { getBook, listChapters, createFavorite, deleteFavorite, listFavorites, listCategories, updateFavorite, ApiClientError } from '$lib/api';
	import { isAuthenticated } from '$lib/stores';
	import type { Book, Chapter, Favorite, Category } from '$lib/types';

	let bookId = $derived(parseInt($page.params.id || '0'));
	let book = $state<Book | null>(null);
	let chapters = $state<Chapter[]>([]);
	let categories = $state<Category[]>([]);
	let favorite = $state<Favorite | null>(null);
	let selectedCategories = $state<string[]>([]);
	
	let loading = $state(true);
	let error = $state<string | null>(null);
	let chaptersPage = $state(1);
	let hasMoreChapters = $state(true);
	let loadingMoreChapters = $state(false);

	async function loadBook() {
		try {
			book = await getBook(bookId);
		} catch (err) {
			if (err instanceof ApiClientError) {
				error = `Error ${err.status}: ${err.message}`;
			} else {
				error = 'Failed to load book';
			}
		}
	}

	async function loadChapters(pageNum: number = 1) {
		if (pageNum === 1) {
			loading = true;
		} else {
			loadingMoreChapters = true;
		}

		try {
			const response = await listChapters({ book_id: bookId, page: pageNum, limit: 50 });
			
			if (pageNum === 1) {
				chapters = response.items;
			} else {
				chapters = [...chapters, ...response.items];
			}
			
			hasMoreChapters = pageNum < response.pagination.total_pages;
			chaptersPage = pageNum;
		} catch (err) {
			console.error('Failed to load chapters:', err);
		} finally {
			loading = false;
			loadingMoreChapters = false;
		}
	}

	async function loadFavoriteStatus() {
		if (!$isAuthenticated) return;
		
		try {
			const response = await listFavorites({ limit: 1000 });
			favorite = response.items.find(f => f.book_id === bookId) || null;
			
			if (favorite && favorite.categories) {
				selectedCategories = favorite.categories;
			}
		} catch (err) {
			console.error('Failed to load favorite status:', err);
		}
	}

	async function loadCategories() {
		if (!$isAuthenticated) return;
		
		try {
			const response = await listCategories({ limit: 100 });
			categories = response.items;
		} catch (err) {
			console.error('Failed to load categories:', err);
		}
	}

	async function toggleFavorite() {
		if (!$isAuthenticated) {
			window.location.href = '/auth/login';
			return;
		}

		try {
			if (favorite) {
				await deleteFavorite(favorite.id);
				favorite = null;
				selectedCategories = [];
			} else {
				favorite = await createFavorite({
					user_id: 1, // Will be set by backend from session
					book_id: bookId
				});
			}
		} catch (err) {
			console.error('Failed to toggle favorite:', err);
		}
	}

	async function updateCategories(categoryName: string) {
		if (!favorite) return;

		const newCategories = selectedCategories.includes(categoryName)
			? selectedCategories.filter(c => c !== categoryName)
			: [...selectedCategories, categoryName];

		selectedCategories = newCategories;

		try {
			await updateFavorite(favorite.id, {
				categories: newCategories.join(',')
			});
		} catch (err) {
			console.error('Failed to update categories:', err);
		}
	}

	function loadMoreChapters() {
		if (!loadingMoreChapters && hasMoreChapters) {
			loadChapters(chaptersPage + 1);
		}
	}

	let sentinel = $state<HTMLElement>();
	
	onMount(() => {
		// Async initialization
		(async () => {
			await Promise.all([
				loadBook(),
				loadChapters(),
				loadFavoriteStatus(),
				loadCategories()
			]);
		})();

		// Setup intersection observer for infinite scroll
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && !loadingMoreChapters && hasMoreChapters) {
					loadMoreChapters();
				}
			},
			{ threshold: 0.1 }
		);

		if (sentinel) {
			observer.observe(sentinel);
		}

		return () => observer.disconnect();
	});
</script>

<svelte:head>
	<title>{book?.title || 'Loading...'} - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if loading && !book}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			<p class="mt-4 text-gray-600">Loading...</p>
		</div>
	{:else if error}
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
			<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<p class="text-red-800 font-medium">{error}</p>
				<a href="/search" class="mt-4 inline-block px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition">
					Back to Search
				</a>
			</div>
		</div>
	{:else if book}
		<!-- Book Header -->
		<div class="bg-white border-b">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
				<div class="flex flex-col md:flex-row gap-8">
					<!-- Cover Image -->
				<div class="shrink-0">
						<div class="w-48 aspect-2/3 bg-gray-100 rounded-lg overflow-hidden shadow-lg">
							{#if book.image_url}
								<img src={book.image_url} alt={book.title} class="w-full h-full object-cover" />
							{:else}
								<div class="w-full h-full flex items-center justify-center text-gray-400 text-6xl">
									📖
								</div>
							{/if}
						</div>
					</div>

					<!-- Book Info -->
					<div class="flex-1">
						<div class="flex items-start justify-between gap-4">
							<div class="flex-1">
								<span class="inline-block px-3 py-1 text-sm font-medium rounded {book.type === 'manga' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'} mb-2">
									{book.type}
								</span>
								<h1 class="text-3xl font-bold text-gray-900 mb-2">{book.title}</h1>
								{#if book.source}
									<p class="text-gray-600 mb-4">
										<span class="font-medium">Source:</span> {book.source.name}
									</p>
								{/if}
							</div>

							<!-- Favorite Button -->
							<button
								onclick={toggleFavorite}
								class="flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition {favorite ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
							>
								<span class="text-2xl">{favorite ? '❤️' : '🤍'}</span>
								{favorite ? 'Favorited' : 'Add to Favorites'}
							</button>
						</div>

						<!-- Synopsis -->
						{#if book.synopsis}
							<div class="mt-4">
								<h3 class="font-semibold text-gray-900 mb-2">Synopsis</h3>
								<p class="text-gray-700 leading-relaxed">{book.synopsis}</p>
							</div>
						{/if}

						<!-- Categories (if favorited) -->
						{#if favorite && categories.length > 0}
							<div class="mt-6">
								<h3 class="font-semibold text-gray-900 mb-3">Categories</h3>
								<div class="flex flex-wrap gap-2">
									{#each categories as category (category.id)}
										<button
											onclick={() => updateCategories(category.name || '')}
											class="px-3 py-1.5 rounded-lg text-sm font-medium transition {selectedCategories.includes(category.name || '') ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
										>
											{selectedCategories.includes(category.name || '') ? '✓ ' : ''}{category.name}
										</button>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>

		<!-- Chapters List -->
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<h2 class="text-2xl font-bold text-gray-900 mb-4">
				Chapters ({chapters.length})
			</h2>

			{#if chapters.length === 0 && !loading}
				<div class="bg-white rounded-lg shadow p-12 text-center">
					<div class="text-6xl mb-4">📄</div>
					<h3 class="text-xl font-semibold text-gray-900 mb-2">No chapters available</h3>
					<p class="text-gray-600">Check back later for updates</p>
				</div>
			{:else}
				<div class="bg-white rounded-lg shadow divide-y">
					{#each chapters as chapter (chapter.id)}
						<a
							href="/content/{chapter.id}"
							class="block px-6 py-4 hover:bg-gray-50 transition"
						>
							<div class="flex items-center justify-between">
								<div class="flex-1">
									<h3 class="font-medium text-gray-900">{chapter.title}</h3>
									{#if chapter.number}
										<p class="text-sm text-gray-500 mt-1">Chapter {chapter.number}</p>
									{/if}
								</div>
								{#if chapter.upload_date}
									<p class="text-sm text-gray-500">
										{new Date(chapter.upload_date).toLocaleDateString()}
									</p>
								{/if}
							</div>
						</a>
					{/each}
				</div>

				<!-- Loading More Indicator -->
				{#if loadingMoreChapters}
					<div class="text-center py-8">
						<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
					</div>
				{/if}

				<!-- Sentinel for Infinite Scroll -->
				{#if hasMoreChapters}
					<div bind:this={sentinel} class="h-20"></div>
				{/if}
			{/if}
		</div>
	{/if}
</div>
