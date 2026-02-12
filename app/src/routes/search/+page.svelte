<script lang="ts">
	import { onMount } from 'svelte';
	import { listBooks, listSources, createFavorite, deleteFavorite, listFavorites, ApiClientError } from '$lib/api';
	import { isAuthenticated } from '$lib/stores';
	import type { Book, Source, BookType, Favorite } from '$lib/types';

	let books = $state<Book[]>([]);
	let sources = $state<Source[]>([]);
	let favorites = $state<Favorite[]>([]);
	let favBookIds = $derived(new Set(favorites.map(f => f.book_id)));
	
	let loading = $state(true);
	let error = $state<string | null>(null);
	let page = $state(1);
	let totalPages = $state(1);
	let total = $state(0);
	let limit = $state(20);
	let loadingMore = $state(false);
	let hasMore = $derived(page < totalPages);

	// Filters
	let selectedSource = $state<number | undefined>(undefined);
	let selectedType = $state<BookType | undefined>(undefined);
	let searchQuery = $state('');

	async function loadBooks(pageNum: number = page, append: boolean = false) {
		if (append) {
			loadingMore = true;
		} else {
			loading = true;
		}
		error = null;

		try {
			const response = await listBooks({
				page: pageNum,
				limit,
				source_id: selectedSource,
				type: selectedType,
				search: searchQuery || undefined
			});

			if (append) {
				books = [...books, ...response.items];
			} else {
				books = response.items;
			}
			
			totalPages = response.pagination.total_pages;
			total = response.pagination.total;
			page = pageNum;
		} catch (err) {
			if (err instanceof ApiClientError) {
				error = `Error ${err.status}: ${err.message}`;
			} else {
				error = 'An unexpected error occurred';
			}
			if (!append) {
				books = [];
			}
		} finally {
			loading = false;
			loadingMore = false;
		}
	}

	async function loadSources() {
		try {
			const response = await listSources({ limit: 100 });
			sources = response.items;
		} catch (err) {
			console.error('Failed to load sources:', err);
		}
	}

	async function loadFavorites() {
		if (!$isAuthenticated) return;
		
		try {
			const response = await listFavorites({ limit: 1000 });
			favorites = response.items;
		} catch (err) {
			console.error('Failed to load favorites:', err);
		}
	}

	async function toggleFavorite(book: Book) {
		if (!$isAuthenticated) {
			window.location.href = '/auth/login';
			return;
		}

		try {
			const existingFav = favorites.find(f => f.book_id === book.id);
			
			if (existingFav) {
				await deleteFavorite(existingFav.id);
				favorites = favorites.filter(f => f.id !== existingFav.id);
			} else {
				const newFav = await createFavorite({
					user_id: 1, // Will be set by backend from session
					book_id: book.id
				});
				favorites = [...favorites, newFav];
			}
		} catch (err) {
			console.error('Failed to toggle favorite:', err);
		}
	}

	function loadMore() {
		if (!loadingMore && hasMore) {
			loadBooks(page + 1, true);
		}
	}

	function handleFilterChange() {
		page = 1;
		loadBooks(1, false);
	}

	function clearFilters() {
		selectedSource = undefined;
		selectedType = undefined;
		searchQuery = '';
		page = 1;
		loadBooks(1, false);
	}

	// Intersection Observer for lazy loading
	let sentinel = $state<HTMLElement>();
	onMount(() => {
		loadSources();
		loadFavorites();
		loadBooks();

		// Setup intersection observer for infinite scroll
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && !loadingMore && hasMore) {
					loadMore();
				}
			},
			{ threshold: 0.1 }
		);

		if (sentinel) {
			observer.observe(sentinel);
		}

		return () => {
			if (sentinel) {
				observer.unobserve(sentinel);
			}
		};
	});
</script>

<svelte:head>
	<title>Search Books - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Header -->
		<div class="mb-6">
			<h1 class="text-3xl font-bold text-gray-900">🔍 Search Books</h1>
			<p class="text-gray-600 mt-1">Discover manga and novels to read</p>
		</div>

		<!-- Filters -->
		<div class="bg-white rounded-lg shadow p-6 mb-6">
			<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
				<!-- Search -->
				<div class="md:col-span-2">
					<label for="search" class="block text-sm font-medium text-gray-700 mb-1">
						Search
					</label>
					<input
						id="search"
						type="text"
						bind:value={searchQuery}
						onchange={handleFilterChange}
						placeholder="Search books..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					/>
				</div>

				<!-- Type Filter -->
				<div>
					<label for="type" class="block text-sm font-medium text-gray-700 mb-1">Type</label>
					<select
						id="type"
						bind:value={selectedType}
						onchange={handleFilterChange}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					>
						<option value={undefined}>All Types</option>
						<option value="manga">Manga</option>
						<option value="novel">Novel</option>
					</select>
				</div>

				<!-- Source Filter -->
				<div>
					<label for="source" class="block text-sm font-medium text-gray-700 mb-1">Source</label>
					<select
						id="source"
						bind:value={selectedSource}
						onchange={handleFilterChange}
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
					>
						<option value={undefined}>All Sources</option>
						{#each sources as source (source.id)}
							<option value={source.id}>{source.name}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Clear Filters -->
			{#if selectedSource || selectedType || searchQuery}
				<div class="mt-4">
					<button
						onclick={clearFilters}
						class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 transition"
					>
						Clear Filters
					</button>
				</div>
			{/if}
		</div>

		<!-- Loading State -->
		{#if loading}
			<div class="text-center py-12">
				<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
				<p class="mt-4 text-gray-600">Loading books...</p>
			</div>

		<!-- Error State -->
		{:else if error}
			<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
				<p class="text-red-800 font-medium">{error}</p>
				<button
					onclick={() => loadBooks(1, false)}
					class="mt-4 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
				>
					Retry
				</button>
			</div>

		<!-- Empty State -->
		{:else if books.length === 0}
			<div class="bg-white rounded-lg shadow p-12 text-center">
				<div class="text-6xl mb-4">📭</div>
				<h3 class="text-xl font-semibold text-gray-900 mb-2">No books found</h3>
				<p class="text-gray-600">Try adjusting your filters or search query</p>
			</div>

		<!-- Books Grid -->
		{:else}
			<div>
				<!-- Results Count -->
				<div class="mb-4 text-sm text-gray-600">
					Showing {books.length} of {total} books
				</div>

				<!-- Grid -->
				<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 mb-8">
					{#each books as book (book.id)}
						<div class="bg-white rounded-lg shadow hover:shadow-lg transition group">
							<a href="/book/{book.id}" class="block">
								<div class="relative aspect-2/3 overflow-hidden rounded-t-lg bg-gray-100">
									{#if book.image_url}
										<img
											src={book.image_url}
											alt={book.title}
											class="w-full h-full object-cover group-hover:scale-105 transition"
										/>
									{:else}
										<div class="w-full h-full flex items-center justify-center text-gray-400 text-4xl">
											📖
										</div>
									{/if}
									<div class="absolute top-2 right-2">
										<span
											class="px-2 py-1 text-xs font-medium rounded {book.type === 'manga'
												? 'bg-blue-100 text-blue-800'
												: 'bg-purple-100 text-purple-800'}"
										>
											{book.type}
										</span>
									</div>
									
									<!-- Favorite Button -->
									<button
										onclick={(e) => { e.preventDefault(); toggleFavorite(book); }}
										class="absolute top-2 left-2 w-8 h-8 bg-white/90 hover:bg-white rounded-full flex items-center justify-center transition"
									>
										<span class="text-xl">{favBookIds.has(book.id) ? '❤️' : '🤍'}</span>
									</button>
								</div>
							</a>
							<div class="p-3">
								<a href="/book/{book.id}">
									<h3 class="font-medium text-sm text-gray-900 line-clamp-2" title={book.title}>
										{book.title}
									</h3>
								</a>
								{#if book.source}
									<p class="text-xs text-gray-500 mt-1">{book.source.name}</p>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Loading More Indicator -->
				{#if loadingMore}
					<div class="text-center py-8">
						<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
						<p class="mt-2 text-gray-600">Loading more...</p>
					</div>
				{/if}

				<!-- Sentinel for Infinite Scroll -->
				{#if hasMore}
					<div bind:this={sentinel} class="h-20"></div>
				{/if}
			</div>
		{/if}
	</div>
</div>

<style>
	.line-clamp-2 {
		display: -webkit-box;
		-webkit-line-clamp: 2;
		line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
	}
</style>
