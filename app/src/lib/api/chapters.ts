import { apiRequest } from './client';
import type { Chapter, PaginatedResponse, ListChaptersParams } from '../types';

/**
 * List chapters with optional filters
 */
export async function listChapters(params?: ListChaptersParams): Promise<PaginatedResponse<Chapter>> {
	return apiRequest('/api/v1/reader/chapters', {
		method: 'GET',
		params: params as Record<string, string | number | boolean>
	});
}

/**
 * Get a single chapter by ID
 */
export async function getChapter(id: number): Promise<Chapter> {
	return apiRequest(`/api/v1/reader/chapters/${id}`, {
		method: 'GET'
	});
}
