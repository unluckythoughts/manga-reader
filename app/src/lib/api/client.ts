// Base API client with fetch wrapper

export interface ApiError {
	error: string;
	message?: string;
	status: number;
}

export class ApiClientError extends Error {
	constructor(
		public status: number,
		public error: string,
		message?: string
	) {
		super(message || error);
		this.name = 'ApiClientError';
	}
}

export interface FetchOptions extends RequestInit {
	params?: Record<string, string | number | boolean | undefined>;
}

/**
 * Base API client configuration
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Clear all cookies - useful when encountering securecookie errors
 */
function clearAllCookies(): void {
	const cookies = document.cookie.split(';');
	
	for (const cookie of cookies) {
		const eqPos = cookie.indexOf('=');
		const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim();
		
		// Clear for all paths and domains
		document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;
		document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; domain=${window.location.hostname};`;
	}
}

/**
 * Check if error is a securecookie validation error
 */
function isSecureCookieError(error: string | undefined): boolean {
	return error?.toLowerCase().includes('securecookie') || 
	       error?.toLowerCase().includes('the value is not valid') ||
	       false;
}

/**
 * Build URL with query parameters
 */
function buildUrl(endpoint: string, params?: Record<string, string | number | boolean | undefined>): string {
	// If endpoint is already a full URL, use it as base
	let url: URL;
	if (endpoint.startsWith('http://') || endpoint.startsWith('https://')) {
		url = new URL(endpoint);
	} else {
		// Construct full URL from base + endpoint
		const fullUrl = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
		url = new URL(fullUrl);
	}
	
	if (params) {
		Object.entries(params).forEach(([key, value]) => {
			if (value !== undefined && value !== null && value !== '') {
				url.searchParams.append(key, String(value));
			}
		});
	}
	
	return url.toString();
}

/**
 * Generic fetch wrapper with error handling
 */
export async function apiRequest<T>(
	endpoint: string,
	options: FetchOptions = {}
): Promise<T> {
	const { params, ...fetchOptions } = options;
	
	const url = buildUrl(`${API_BASE_URL}${endpoint}`, params);
	
	const config: RequestInit = {
		...fetchOptions,
		credentials: 'include', // Include cookies for auth
		headers: {
			'Content-Type': 'application/json',
			...fetchOptions.headers
		}
	};

	try {
		const response = await fetch(url, config);
		
		// Handle non-OK responses
		if (!response.ok) {
			let errorData: ApiError;
			try {
				errorData = await response.json();
			} catch {
				errorData = {
					error: response.statusText,
					status: response.status
				};
			}
			
			// Check for securecookie errors and clear cookies
			if (isSecureCookieError(errorData.error) || isSecureCookieError(errorData.message)) {
				console.warn('Detected securecookie error, clearing all cookies');
				clearAllCookies();
			}
			
			throw new ApiClientError(response.status, errorData.error, errorData.message);
		}

		// Handle empty responses (204 No Content)
		if (response.status === 204) {
			return {} as T;
		}

		// Parse JSON response
		return await response.json();
	} catch (error) {
		if (error instanceof ApiClientError) {
			throw error;
		}
		throw new ApiClientError(0, 'Network error', (error as Error).message);
	}
}
