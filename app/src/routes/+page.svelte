<script lang="ts">
	import { onMount } from 'svelte';
	import { listFavorites, ApiClientError } from '$lib/api';
	import { auth, isAuthenticated } from '$lib/stores';
	import type { Favorite } from '$lib/types';

	let favorites = $state<Favorite[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let page = $state(1);
	let hasMore = $state(true);
	let loadingMore = $state(false);

	// Group favorites by category
	let groupedFavorites = $derived.by(() => {
		const groups: Record<string, Favorite[]> = {
			'Reading': [],
			'Plan to Read': [],
			'Completed': [],
			'Other': []
		};

		favorites.forEach(fav => {
			const category = fav.categories?.[0] || 'Other';
			if (!groups[category]) {
				groups[category] = [];
			}
			groups[category].push(fav);
		});

		return groups;
	});

	async function loadFavorites(pageNum: number = 1) {
		if (pageNum === 1) {
			loading = true;
		} else {
			loadingMore = true;
		}
		error = null;

		try {
			const response = await listFavorites({ page: pageNum, limit: 20 });
			
			if (pageNum === 1) {
				favorites = response.items;
			} else {
				favorites = [...favorites, ...response.items];
			}
			
			hasMore = page < response.pagination.total_pages;
		} catch (err) {
			if (err instanceof ApiClientError) {
				if (err.status === 401) {
					error = null; // Don't show error for not authenticated
				} else {
					error = `Error ${err.status}: ${err.message}`;
				}
			} else {
				error = 'An unexpected error occurred';
			}
			favorites = [];
		} finally {
			loading = false;
			loadingMore = false;
		}
	}

	function loadMore() {
		if (!loadingMore && hasMore) {
			page++;
			loadFavorites(page);
		}
	}

	// Intersection Observer for lazy loading
	let sentinel = $state<HTMLElement>();
	onMount(() => {
		// Only load if authenticated
		if ($isAuthenticated) {
			loadFavorites();
		} else {
			loading = false;
		}

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
	<title>My Library - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Hero Section -->
	<div class="bg-gradient-to-r from-blue-600 to-purple-600 text-white">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
			<h1 class="text-4xl font-bold mb-2">📚 My Library</h1>
			<p class="text-blue-100 text-lg">Continue reading your favorite books</p>
		</div>
	</div>

	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
		<!-- Not Authenticated -->
		{#if !$isAuthenticated}
			<div class="bg-white rounded-lg shadow-lg p-12 text-center max-w-2xl mx-auto">
				<div class="text-6xl mb-4">🔒</div>
				<h2 class="text-2xl font-bold text-gray-900 mb-4">Sign in to view your library</h2>
				<p class="text-gray-600 mb-8">
					Create an account or sign in to save your favorite books and track your reading progress.
				</p>
				<div class="flex gap-4 justify-center">
					<a
						href="/auth/login"
						class="px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition"
					>
						Sign In
					</a>
					<a
						href="/auth/register"
						class="px-6 py-3 border-2 border-blue-600 text-blue-600 font-semibold rounded-lg hover:bg-blue-50 transition"
					>
						Create Account
					</a>
				</div>
				<div class="mt-8">
					<a href="/search" class="text-blue-600 hover:text-blue-700 font-medium">
						Browse books without signing in →
					</a>
				</div>
			</div>

		<!-- Loading State -->
		{:else if loading}
			<div class="text-center py-12">
				<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
				<p class="mt-4 text-gray-600">Loading your library...</p>
			</div>

		<!-- Error State -->
		{:else if error}
			<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center max-w-2xl mx-auto">
				<p class="text-red-800 font-medium">{error}</p>
				<button
					onclick={() => loadFavorites(1)}
					class="mt-4 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition"
				>
					Retry
				</button>
			</div>

		<!-- Empty State -->
		{:else if favorites.length === 0}
			<div class="bg-white rounded-lg shadow p-12 text-center max-w-2xl mx-auto">
				<div class="text-6xl mb-4">📚</div>
				<h3 class="text-xl font-semibold text-gray-900 mb-2">Your library is empty</h3>
				<p class="text-gray-600 mb-6">
					Start adding books to your favorites to see them here!
				</p>
				<a
					href="/search"
					class="inline-block px-6 py-3 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition"
				>
					Discover Books
				</a>
			</div>

		<!-- Favorites Grid by Category -->
		{:else}
			<div class="space-y-8">
				{#each Object.entries(groupedFavorites) as [category, items]}
					{#if items.length > 0}
						<div>
							<div class="flex items-center justify-between mb-4">
								<h2 class="text-2xl font-bold text-gray-900">{category}</h2>
								<span class="text-sm text-gray-500">{items.length} books</span>
							</div>

							<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
								{#each items as favorite (favorite.id)}
									{#if favorite.book}
										<a
											href="/book/{favorite.book.id}"
											class="bg-white rounded-lg shadow hover:shadow-lg transition group"
										>
											<div class="relative aspect-2/3 overflow-hidden rounded-t-lg bg-gray-100">
												{#if favorite.book.image_url}
													<img
														src={favorite.book.image_url}
														alt={favorite.book.title}
														class="w-full h-full object-cover group-hover:scale-105 transition"
													/>
												{:else}
													<div class="w-full h-full flex items-center justify-center text-gray-400 text-4xl">
														📖
													</div>
												{/if}
												
												<!-- Progress Badge -->
												{#if favorite.progress && favorite.progress.length > 0}
													<div class="absolute bottom-2 left-2 right-2">
														<div class="bg-black/75 text-white text-xs px-2 py-1 rounded">
															Ch. {favorite.progress[0]}
														</div>
													</div>
												{/if}
											</div>
											<div class="p-3">
												<h3 class="font-medium text-sm text-gray-900 line-clamp-2" title={favorite.book.title}>
													{favorite.book.title}
												</h3>
											</div>
										</a>
									{/if}
								{/each}
							</div>
						</div>
					{/if}
				{/each}
			</div>

			<!-- Loading More Indicator -->
			{#if loadingMore}
				<div class="text-center py-8">
					<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
				</div>
			{/if}

			<!-- Sentinel for Infinite Scroll -->
			{#if hasMore}
				<div bind:this={sentinel} class="h-20"></div>
			{/if}
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
