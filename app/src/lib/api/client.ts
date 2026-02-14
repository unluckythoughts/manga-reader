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
const ENABLE_AUTH = import.meta.env.ENABLE_AUTH !== 'false' && import.meta.env.VITE_ENABLE_AUTH !== 'false';

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
	options: FetchOptions = {},
	isRetry: boolean = false
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
		
		// Parse response body first to check for securecookie errors
		let responseData: any;
		const contentType = response.headers.get('content-type');
		
		if (response.status === 204) {
			responseData = {};
		} else if (contentType?.includes('application/json')) {
			try {
				responseData = await response.json();
			} catch {
				responseData = null;
			}
		}
		
		// Check for securecookie errors in response (even if ok: false)
		if (responseData && !isRetry) {
			const hasSecureCookieError = 
				isSecureCookieError(responseData.error) || 
				isSecureCookieError(responseData.message);
			
			if (hasSecureCookieError) {
				console.warn('Detected securecookie error - re-authenticating');
				
				// Re-authenticate to get a fresh cookie (overwrites the old one)
				// Note: Can't clear HttpOnly cookies from JavaScript - they must be cleared server-side
				if (!ENABLE_AUTH) {
					console.log('Auto-logging in with dummy credentials...');
					try {
						const loginResponse = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
							method: 'POST',
							credentials: 'include',
							headers: { 'Content-Type': 'application/json' },
							body: JSON.stringify({ username: 'dummy@example.com', password: 'dummy' })
						});
						
						if (!loginResponse.ok) {
							console.error('Auto-login failed with status:', loginResponse.status);
							throw new Error('Auto-login failed');
						}
						
						console.log('Auto-login successful, retrying request');
					} catch (loginError) {
						console.error('Auto-login failed:', loginError);
						throw loginError;
					}
				}
				
				// Retry the request once with fresh auth cookie
				return apiRequest<T>(endpoint, options, true);
			}
		}
		
		// Handle non-OK responses
		if (!response.ok) {
			const errorData: ApiError = responseData || {
				error: response.statusText,
				status: response.status
			};
			
			throw new ApiClientError(response.status, errorData.error, errorData.message);
		}

		return responseData as T;
	} catch (error) {
		if (error instanceof ApiClientError) {
			throw error;
		}
		throw new ApiClientError(0, 'Network error', (error as Error).message);
	}
}