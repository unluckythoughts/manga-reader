export interface Pagination {
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination: Pagination
}

export type BookType = 'manga' | 'novel'

export interface Source {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  name: string
  domain: string
  icon_url?: string
}

export interface Chapter {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  url: string
  title: string
  book_id?: number
  number?: string
  content?: string[]
  upload_date?: string
  completed: boolean
  downloaded: boolean
  book?: Book
}

export interface Book {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  url: string
  title: string
  type: BookType
  image_url?: string
  synopsis?: string
  slug?: string
  source_id?: number
  source?: Source
  chapters?: Chapter[]
}

export interface Favorite {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  user_id?: number
  book_id?: number
  progress?: string[]
  categories?: string[]
  book?: Book
}

export interface Category {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  name: string
}

export interface User {
  ID: number
  name: string
  email: string
  role: number
}

export interface LoginResponse {
  token: string
}
