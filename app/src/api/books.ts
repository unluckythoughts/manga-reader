import { apiClient } from './client'
import type { Book, PaginatedResponse } from '../types'

export const booksApi = {
  list: (page = 1, limit = 20, sourceId?: number) => {
    const params = new URLSearchParams({ page: String(page), limit: String(limit) })
    if (sourceId) params.set('source_id', String(sourceId))
    return apiClient.get<PaginatedResponse<Book>>(`/reader/books?${params}`)
  },
  get: (id: number) => apiClient.get<Book>(`/reader/books/${id}`)
}
