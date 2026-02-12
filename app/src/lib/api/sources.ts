import { apiRequest } from './client';
import type { Source, PaginatedResponse, ListSourcesParams } from '../types';

/**
 * List sources with optional filters
 */
export async function listSources(params?: ListSourcesParams): Promise<PaginatedResponse<Source>> {
	return apiRequest('/api/v1/reader/sources', {
		method: 'GET',
		params: params as Record<string, string | number | boolean>
	});
}

/**
 * Get a single source by ID
 */
export async function getSource(id: number): Promise<Source> {
	return apiRequest(`/api/v1/reader/sources/${id}`, {
		method: 'GET'
	});
}
