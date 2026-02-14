import { writable } from 'svelte/store';
import type { User } from '../types';
import * as authApi from '../api/auth';
import { ApiClientError } from '../api/client';

interface AuthState {
	user: User | null;
	loading: boolean;
	error: string | null;
}

const ENABLE_AUTH = import.meta.env.ENABLE_AUTH !== 'false' && import.meta.env.VITE_ENABLE_AUTH !== 'false';

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		loading: false,
		error: null
	});

	return {
		subscribe,

		/**
		 * Initialize auth state by fetching current user
		 */
		async init() {
			update(state => ({ ...state, loading: true, error: null }));
			
			// If auth is disabled, auto-login with dummy credentials first
			if (!ENABLE_AUTH) {
				try {
					console.log('Auth disabled, auto-logging in with dummy credentials...');
					const { user } = await authApi.login({ 
						username: 'dummy@example.com', 
						password: 'dummy' 
					});
					console.log('Auto-login successful');
					set({ user, loading: false, error: null });
					return user;
				} catch (error: unknown) {
					console.error('Auto-login failed:', error);
					// Continue to try getting user anyway
				}
			}
			
			try {
				const user = await authApi.getUser();
				set({ user, loading: false, error: null });
				return user;
			} catch (error: unknown) {
				if (error instanceof ApiClientError && error.status === 401) {
					// Not authenticated, this is expected
					set({ user: null, loading: false, error: null });
				} else if (error instanceof ApiClientError && 
				          (error.message.toLowerCase().includes('securecookie') || 
				           error.error.toLowerCase().includes('securecookie'))) {
					// Securecookie error - cookies were already cleared and retried
					console.info('Session invalid, logged out');
					set({ user: null, loading: false, error: null });
				} else {
					const errorMessage = error instanceof ApiClientError 
						? error.message 
						: 'Failed to load user';
					set({ user: null, loading: false, error: errorMessage });
				}
				return null;
			}
		},

		/**
		 * Login with username and password
		 */
		async login(username: string, password: string) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const { user } = await authApi.login({ username, password });
				set({ user, loading: false, error: null });
				return user;
			} catch (error: unknown) {
				const errorMessage = error instanceof ApiClientError 
					? error.message 
					: 'Login failed';
				update(state => ({ ...state, loading: false, error: errorMessage }));
				throw error;
			}
		},

		/**
		 * Register a new user
		 */
		async register(name: string, email: string, password: string) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const { user } = await authApi.register({ name, email, password });
				set({ user, loading: false, error: null });
				return user;
			} catch (error: unknown) {
				const errorMessage = error instanceof ApiClientError 
					? error.message 
					: 'Registration failed';
				update(state => ({ ...state, loading: false, error: errorMessage }));
				throw error;
			}
		},

		/**
		 * Login with Google OAuth
		 */
		async googleLogin(token: string) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const { user } = await authApi.googleOAuthLogin(token);
				set({ user, loading: false, error: null });
				return user;
			} catch (error: unknown) {
				const errorMessage = error instanceof ApiClientError 
					? error.message 
					: 'Google login failed';
				update(state => ({ ...state, loading: false, error: errorMessage }));
				throw error;
			}
		},

		/**
		 * Logout current user
		 */
		async logout() {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				await authApi.logout();
				set({ user: null, loading: false, error: null });
			} catch (error: unknown) {
				const errorMessage = error instanceof ApiClientError 
					? error.message 
					: 'Logout failed';
				update(state => ({ ...state, loading: false, error: errorMessage }));
				throw error;
			}
		},

		/**
		 * Update current user
		 */
		async updateUser(data: { username?: string; email?: string }) {
			update(state => ({ ...state, loading: true, error: null }));
			
			try {
				const user = await authApi.updateUser(data);
				set({ user, loading: false, error: null });
				return user;
			} catch (error: unknown) {
				const errorMessage = error instanceof ApiClientError 
					? error.message 
					: 'Update failed';
				update(state => ({ ...state, loading: false, error: errorMessage }));
				throw error;
			}
		},

		/**
		 * Clear error message
		 */
		clearError() {
			update(state => ({ ...state, error: null }));
		},

		/**
		 * Reset store to initial state
		 */
		reset() {
			set({ user: null, loading: false, error: null });
		}
	};
}

export const auth = createAuthStore();

// Helper derived store for checking if user is authenticated
export const isAuthenticated = writable(false);
auth.subscribe(state => {
	isAuthenticated.set(state.user !== null);
});
