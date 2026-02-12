<script lang="ts">
	import { onMount } from 'svelte';
	import { listBooks, ApiClientError, type Book, type PaginatedResponse } from '$lib/api';

	let books: Book[] = $state([]);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let page = $state(1);
	let totalPages = $state(1);

	async function loadBooks() {
		loading = true;
		error = null;
		
		try {
			const response: PaginatedResponse<Book> = await listBooks({
				page,
				limit: 20,
				type: 'manga'
			});
			
			books = response.items;
			totalPages = response.pagination.total_pages;
		} catch (err) {
			if (err instanceof ApiClientError) {
				error = `Error ${err.status}: ${err.message}`;
			} else {
				error = 'An unexpected error occurred';
			}
		} finally {
			loading = false;
		}
	}

	onMount(() => {
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
</script>

<div class="books-container">
	<h2>Books</h2>

	{#if loading}
		<p>Loading books...</p>
	{:else if error}
		<p class="error">{error}</p>
	{:else if books.length === 0}
		<p>No books found.</p>
	{:else}
		<div class="books-grid">
			{#each books as book (book.id)}
				<div class="book-card">
					{#if book.image_url}
						<img src={book.image_url} alt={book.title} />
					{/if}
					<h3>{book.title}</h3>
					<p class="book-type">{book.type}</p>
					{#if book.synopsis}
						<p class="synopsis">{book.synopsis.slice(0, 150)}...</p>
					{/if}
				</div>
			{/each}
		</div>

		<div class="pagination">
			<button onclick={prevPage} disabled={page === 1}>Previous</button>
			<span>Page {page} of {totalPages}</span>
			<button onclick={nextPage} disabled={page === totalPages}>Next</button>
		</div>
	{/if}
</div>

<style>
	.books-container {
		padding: 2rem;
	}

	.books-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
		margin: 1rem 0;
	}

	.book-card {
		border: 1px solid #ddd;
		border-radius: 8px;
		padding: 1rem;
		transition: transform 0.2s;
	}

	.book-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	}

	.book-card img {
		width: 100%;
		height: 250px;
		object-fit: cover;
		border-radius: 4px;
	}

	.book-card h3 {
		margin: 0.5rem 0;
		font-size: 1rem;
	}

	.book-type {
		color: #666;
		font-size: 0.875rem;
		text-transform: capitalize;
	}

	.synopsis {
		font-size: 0.875rem;
		color: #555;
		margin-top: 0.5rem;
	}

	.pagination {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
		margin-top: 2rem;
	}

	.pagination button {
		padding: 0.5rem 1rem;
		border: 1px solid #ddd;
		background: white;
		border-radius: 4px;
		cursor: pointer;
	}

	.pagination button:hover:not(:disabled) {
		background: #f5f5f5;
	}

	.pagination button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.error {
		color: #d32f2f;
		padding: 1rem;
		background: #ffebee;
		border-radius: 4px;
	}
</style>
