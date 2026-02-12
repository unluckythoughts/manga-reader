import { apiRequest } from './client';
import type {
	Category,
	PaginatedResponse,
	ListCategoriesParams,
	CreateCategoryRequest,
	UpdateCategoryRequest
} from '../types';

/**
 * List categories with optional filters
 */
export async function listCategories(params?: ListCategoriesParams): Promise<PaginatedResponse<Category>> {
	return apiRequest('/api/v1/reader/categories', {
		method: 'GET',
		params: params as Record<string, string | number | boolean>
	});
}

/**
 * Get a single category by ID
 */
export async function getCategory(id: number): Promise<Category> {
	return apiRequest(`/api/v1/reader/categories/${id}`, {
		method: 'GET'
	});
}

/**
 * Create a new category
 */
export async function createCategory(data: CreateCategoryRequest): Promise<Category> {
	return apiRequest('/api/v1/reader/categories', {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

/**
 * Update an existing category
 */
export async function updateCategory(id: number, data: UpdateCategoryRequest): Promise<Category> {
	return apiRequest(`/api/v1/reader/categories/${id}`, {
		method: 'PUT',
		body: JSON.stringify(data)
	});
}

/**
 * Delete a category
 */
export async function deleteCategory(id: number): Promise<void> {
	return apiRequest(`/api/v1/reader/categories/${id}`, {
		method: 'DELETE'
	});
}
