import { apiClient } from './client'
import type { Source, PaginatedResponse } from '../types'

export const sourcesApi = {
  list: (page = 1, limit = 20) =>
    apiClient.get<PaginatedResponse<Source>>(`/reader/sources?page=${page}&limit=${limit}`),
  get: (id: number) => apiClient.get<Source>(`/reader/sources/${id}`)
}
