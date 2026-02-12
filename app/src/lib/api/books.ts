import { apiRequest } from './client';
import type { Book, PaginatedResponse, ListBooksParams } from '../types';

/**
 * List books with optional filters
 */
export async function listBooks(params?: ListBooksParams): Promise<PaginatedResponse<Book>> {
	return apiRequest('/api/v1/reader/books', {
		method: 'GET',
		params: params as Record<string, string | number | boolean>
	});
}

/**
 * Get a single book by ID
 */
export async function getBook(id: number): Promise<Book> {
	return apiRequest(`/api/v1/reader/books/${id}`, {
		method: 'GET'
	});
}
