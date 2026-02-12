<script lang="ts">
	import { onMount } from 'svelte';
	import { listBooks, listSources, ApiClientError } from '$lib/api';
	import type { Book, Source, BookType } from '$lib/types';

	let books = $state<Book[]>([]);
	let sources = $state<Source[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let page = $state(1);
	let totalPages = $state(1);
	let total = $state(0);
	let limit = $state(20);

	// Filters
	let selectedSource = $state<number | undefined>(undefined);
	let selectedType = $state<BookType | undefined>(undefined);
	let searchQuery = $state('');

	async function loadBooks() {
		loading = true;
		error = null;

		try {
			const response = await listBooks({
				page,
				limit,
				source_id: selectedSource,
				type: selectedType,
				search: searchQuery || undefined
			});

			books = response.items;
			totalPages = response.pagination.total_pages;
			total = response.pagination.total;
		} catch (err) {
			if (err instanceof ApiClientError) {
				error = `Error ${err.status}: ${err.message}`;
			} else {
				error = 'An unexpected error occurred';
			}
			books = [];
		} finally {
			loading = false;
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

	onMount(() => {
		loadSources();
		loadBooks();
	});

	function nextPage() {
		if (page < totalPages) {
			page++;
			loadBooks();
		}
	}

	function prevPage() {
		if (page > 1) {
			page--;
			loadBooks();
		}
	}

	function handleFilterChange() {
		page = 1;
		loadBooks();
	}

	function clearFilters() {
		selectedSource = undefined;
		selectedType = undefined;
		searchQuery = '';
		page = 1;
		loadBooks();
	}
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Header -->
	<header class="bg-white shadow-sm">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
			<h1 class="text-3xl font-bold text-gray-900">📚 Book Reader</h1>
			<p class="text-gray-600 mt-1">Discover and read your favorite manga and novels</p>
		</div>
	</header>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
					onclick={loadBooks}
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
						<div class="bg-white rounded-lg shadow hover:shadow-lg transition group cursor-pointer">
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
							</div>
							<div class="p-3">
								<h3 class="font-medium text-sm text-gray-900 line-clamp-2" title={book.title}>
									{book.title}
								</h3>
								{#if book.source}
									<p class="text-xs text-gray-500 mt-1">{book.source.name}</p>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination -->
				{#if totalPages > 1}
					<div class="flex items-center justify-center gap-4 py-6">
						<button
							onclick={prevPage}
							disabled={page === 1}
							class="px-4 py-2 bg-white border border-gray-300 rounded-md font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							← Previous
						</button>

						<span class="text-gray-700">
							Page <span class="font-semibold">{page}</span> of
							<span class="font-semibold">{totalPages}</span>
						</span>

						<button
							onclick={nextPage}
							disabled={page === totalPages}
							class="px-4 py-2 bg-white border border-gray-300 rounded-md font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition"
						>
							Next →
						</button>
					</div>
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
