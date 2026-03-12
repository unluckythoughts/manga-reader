import { apiClient } from './client'
import type { Book, PaginatedResponse } from '../types'

export const booksApi = {
  list: (page = 1, limit = 20, sourceId?: number, search?: string, bookType?: string) => {
    const params = new URLSearchParams({ page: String(page), limit: String(limit) })
    if (sourceId) params.set('source_id', String(sourceId))
    if (search) params.set('search', search)
    if (bookType) params.set('book_type', bookType)
    return apiClient.get<PaginatedResponse<Book>>(`/reader/books?${params}`)
  },
  get: (id: number, force = false) => {
    const params = new URLSearchParams()
    if (force) params.set('force', 'true')
    const query = params.toString()
    return apiClient.get<Book>(`/reader/books/${id}${query ? `?${query}` : ''}`)
  }
}
