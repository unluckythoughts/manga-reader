import { apiRequest } from './client';
import type {
	Favorite,
	PaginatedResponse,
	ListFavoritesParams,
	CreateFavoriteRequest,
	UpdateFavoriteRequest,
	UpdateFavoriteProgressRequest
} from '../types';

/**
 * List favorites with optional filters
 */
export async function listFavorites(params?: ListFavoritesParams): Promise<PaginatedResponse<Favorite>> {
	return apiRequest('/api/v1/reader/favorites', {
		method: 'GET',
		params: params as Record<string, string | number | boolean>
	});
}

/**
 * Get a single favorite by ID
 */
export async function getFavorite(id: number): Promise<Favorite> {
	return apiRequest(`/api/v1/reader/favorites/${id}`, {
		method: 'GET'
	});
}

/**
 * Create a new favorite
 */
export async function createFavorite(data: CreateFavoriteRequest): Promise<Favorite> {
	return apiRequest('/api/v1/reader/favorites', {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

/**
 * Update an existing favorite
 */
export async function updateFavorite(id: number, data: UpdateFavoriteRequest): Promise<Favorite> {
	return apiRequest(`/api/v1/reader/favorites/${id}`, {
		method: 'PUT',
		body: JSON.stringify(data)
	});
}

/**
 * Update favorite reading progress
 */
export async function updateFavoriteProgress(id: number, data: UpdateFavoriteProgressRequest): Promise<Favorite> {
	return apiRequest(`/api/v1/reader/favorites/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(data)
	});
}

/**
 * Delete a favorite
 */
export async function deleteFavorite(id: number): Promise<void> {
	return apiRequest(`/api/v1/reader/favorites/${id}`, {
		method: 'DELETE'
	});
}
