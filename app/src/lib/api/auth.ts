import { apiRequest } from './client';
import type {
	User,
	LoginRequest,
	RegisterRequest,
	UpdateUserRequest,
	ChangePasswordRequest,
	UpdatePasswordRequest
} from '../types';

/**
 * Login with username and password
 */
export async function login(data: LoginRequest): Promise<{ user: User; token: string }> {
	return apiRequest('/api/v1/auth/login', {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

/**
 * Register a new user
 */
export async function register(data: RegisterRequest): Promise<{ user: User; token: string }> {
	return apiRequest('/api/v1/auth/register', {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

/**
 * Logout current user
 */
export async function logout(): Promise<void> {
	return apiRequest('/api/v1/auth/logout', {
		method: 'POST'
	});
}

/**
 * Login with Google OAuth
 */
export async function googleOAuthLogin(token: string): Promise<{ user: User; token: string }> {
	return apiRequest('/api/v1/oauth/login/google', {
		method: 'POST',
		body: JSON.stringify({ token })
	});
}

/**
 * Send verification token to email or phone
 * @param target - 'email' or 'phone'
 */
export async function sendVerificationToken(target: 'email' | 'phone'): Promise<void> {
	return apiRequest(`/api/v1/auth/verify/${target}`, {
		method: 'PATCH'
	});
}

/**
 * Verify token for email or phone
 * @param target - 'email' or 'phone'
 * @param token - Verification token
 */
export async function verifyToken(target: 'email' | 'phone', token: string): Promise<void> {
	return apiRequest(`/api/v1/auth/verify/${target}/${token}`, {
		method: 'GET'
	});
}

/**
 * Get current user information
 */
export async function getUser(): Promise<User> {
	return apiRequest('/api/v1/user', {
		method: 'GET'
	});
}

/**
 * Update current user information
 */
export async function updateUser(data: UpdateUserRequest): Promise<User> {
	return apiRequest('/api/v1/user', {
		method: 'PUT',
		body: JSON.stringify(data)
	});
}

/**
 * Change user password
 */
export async function changePassword(data: ChangePasswordRequest): Promise<void> {
	return apiRequest('/api/v1/user/change-password', {
		method: 'PATCH',
		body: JSON.stringify(data)
	});
}

/**
 * Request password reset (sends reset link via email)
 */
export async function resetPassword(email: string): Promise<void> {
	return apiRequest('/api/v1/user/reset-password', {
		method: 'GET',
		params: { email }
	});
}

/**
 * Update password with reset token
 */
export async function updatePassword(data: UpdatePasswordRequest): Promise<void> {
	return apiRequest('/api/v1/user/update-password', {
		method: 'PATCH',
		body: JSON.stringify(data)
	});
}
