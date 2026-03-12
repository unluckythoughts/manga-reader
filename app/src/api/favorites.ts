import { apiClient } from './client'
import type { Favorite, PaginatedResponse } from '../types'

export const favoritesApi = {
  list: (page = 1, limit = 20) =>
    apiClient.get<PaginatedResponse<Favorite>>(`/reader/favorites?page=${page}&limit=${limit}`),
  get: (id: number) => apiClient.get<Favorite>(`/reader/favorites/${id}`),
  create: (bookId: number, userId: number) =>
    apiClient.post<Favorite>('/reader/favorites', { book_id: bookId, user_id: userId }),
  update: (id: number, data: { progress?: string; categories?: string }) =>
    apiClient.put<Favorite>(`/reader/favorites/${id}`, data),
  updateProgress: (id: number, chapter: number, level = 0) =>
    apiClient.patch<Favorite>(`/reader/favorites/${id}`, { chapter, level }),
  delete: (id: number) => apiClient.delete<void>(`/reader/favorites/${id}`)
}
