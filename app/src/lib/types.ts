// Core types matching Go models

export type BookType = 'manga' | 'novel';

export interface Book {
	id: number;
	url: string;
	title: string;
	type: BookType;
	image_url?: string;
	synopsis?: string;
	slug?: string;
	other_id?: string;
	source_id?: number;
	updated_at: string;
	deleted_at?: string;
	source?: Source;
	chapters?: Chapter[];
}

export interface Chapter {
	id: number;
	url: string;
	title: string;
	book_id?: number;
	number?: string;
	content?: string[];
	upload_date?: string;
	completed: boolean;
	downloaded: boolean;
	other_id?: string;
	updated_at: string;
	deleted_at?: string;
	book?: Book;
}

export interface Source {
	id: number;
	name: string;
	domain: string;
	icon_url?: string;
	updated_at: string;
	deleted_at?: string;
	books?: Book[];
}

export interface Category {
	id: number;
	name?: string;
	updated_at: string;
	deleted_at?: string;
}

export interface Favorite {
	id: number;
	user_id?: number;
	book_id?: number;
	progress?: string[];
	categories?: string[];
	updated_at: string;
	deleted_at?: string;
	user?: User;
	book?: Book;
}

export interface User {
	id: number;
	username: string;
	email: string;
	role: number;
	verified: boolean;
	created_at: string;
	updated_at: string;
}

// Pagination types
export interface Pagination {
	page: number;
	limit: number;
	total: number;
	total_pages: number;
}

export interface PaginatedResponse<T> {
	items: T[];
	pagination: Pagination;
}

// Request types
export interface CreateBookRequest {
	url: string;
	title: string;
	type: BookType;
	image_url?: string;
	synopsis?: string;
	slug?: string;
	other_id?: string;
	source_id?: number;
}

export interface CreateChapterRequest {
	url: string;
	title: string;
	book_id: number;
	number?: string;
	content?: string;
	upload_date?: string;
	completed?: boolean;
	downloaded?: boolean;
	other_id?: string;
}

export interface CreateFavoriteRequest {
	user_id: number;
	book_id: number;
	progress?: string;
	categories?: string;
}

export interface UpdateFavoriteRequest {
	progress?: string;
	categories?: string;
}

export interface UpdateFavoriteProgressRequest {
	chapter: number;
	level?: number;
}

export interface CreateCategoryRequest {
	name: string;
}

export interface UpdateCategoryRequest {
	name: string;
}

// Auth request types
export interface LoginRequest {
	username: string;
	password: string;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
}

export interface UpdateUserRequest {
	username?: string;
	email?: string;
}

export interface ChangePasswordRequest {
	old_password: string;
	new_password: string;
}

export interface UpdatePasswordRequest {
	token: string;
	password: string;
}

// Query parameters
export interface ListBooksParams {
	page?: number;
	limit?: number;
	source_id?: number;
	type?: BookType;
	search?: string;
}

export interface ListChaptersParams {
	page?: number;
	limit?: number;
	book_id?: number;
}

export interface ListSourcesParams {
	page?: number;
	limit?: number;
}

export interface ListCategoriesParams {
	page?: number;
	limit?: number;
}

export interface ListFavoritesParams {
	page?: number;
	limit?: number;
	category?: string;
}
