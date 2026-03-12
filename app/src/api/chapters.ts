import { apiClient } from './client'
import type { Chapter, PaginatedResponse } from '../types'

export const chaptersApi = {
  list: (bookId?: number, page = 1, limit = 100) => {
    const params = new URLSearchParams({ page: String(page), limit: String(limit) })
    if (bookId) params.set('book_id', String(bookId))
    return apiClient.get<PaginatedResponse<Chapter>>(`/reader/chapters?${params}`)
  },
  get: (id: number) => apiClient.get<Chapter>(`/reader/chapters/${id}`)
}
