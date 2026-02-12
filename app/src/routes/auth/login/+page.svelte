<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import type { LoginRequest } from '$lib/types';

	let formData = $state<LoginRequest>({
		username: '',
		password: ''
	});

	let showPassword = $state(false);
	let rememberMe = $state(false);
	let touched = $state({
		username: false,
		password: false
	});

	// Field validation errors
	let fieldErrors = $derived.by(() => {
		const errors: Record<string, string> = {};
		
		if (touched.username && !formData.username) {
			errors.username = 'Email or username is required';
		}
		
		if (touched.password && !formData.password) {
			errors.password = 'Password is required';
		}
		
		return errors;
	});

	function markTouched(field: keyof typeof touched) {
		touched[field] = true;
	}

	function togglePassword(e: MouseEvent) {
		e.preventDefault();
		e.stopPropagation();
		showPassword = !showPassword;
	}

	async function handleSubmit() {
		try {
			await auth.login(formData.username, formData.password);
			// Redirect to home
			window.location.href = '/';
		} catch (error) {
			// Error is already set in the store
			console.error('Login failed:', error);
		}
	}

	async function handleGoogleLogin() {
		// TODO: Implement Google OAuth flow
		// This would typically open a popup or redirect to Google OAuth
		alert('Google OAuth not yet configured. Set up clientId in backend.');
	}
</script>

<svelte:head>
	<title>Login - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center px-4 py-12">
	<div class="w-full max-w-md">
		<div class="bg-white rounded-2xl shadow-2xl p-8">
			<!-- Logo/Title -->
			<div class="text-center mb-8">
				<div class="text-5xl mb-3">📚</div>
				<h1 class="text-3xl font-bold text-gray-900">Welcome Back</h1>
				<p class="text-gray-600 mt-2">Sign in to continue reading</p>
			</div>

			<!-- Error Message -->
			{#if $auth.error}
				<div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start gap-3">
					<span class="text-red-600 text-xl">⚠️</span>
					<div class="flex-1">
						<p class="text-red-800 text-sm">{$auth.error}</p>
					</div>
					<button
						onclick={() => auth.clearError()}
						class="text-red-600 hover:text-red-800 text-xl leading-none"
					>
						×
					</button>
				</div>
			{/if}

			<!-- Login Form -->
			<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-5">
				<!-- Username/Email -->
				<div>
					<label for="username" class="block text-sm font-medium text-gray-700 mb-2">
						Username or Email
					</label>
					<input
						id="username"
						type="text"
						bind:value={formData.username}
					onblur={() => markTouched('username')}
					required
					disabled={$auth.loading}
					placeholder="Enter your username or email"
					class="w-full px-4 py-3 border {fieldErrors.username ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
				/>
				{#if fieldErrors.username}
					<p class="mt-1 text-sm text-red-600">{fieldErrors.username}</p>
				{/if}
			</div>

			<!-- Password -->
			<div>
				<label for="password" class="block text-sm font-medium text-gray-700 mb-2">
					Password
				</label>
				<div class="relative">
					<input
						id="password"
						type={showPassword ? 'text' : 'password'}
						bind:value={formData.password}
						onblur={() => markTouched('password')}
						required
						disabled={$auth.loading}
						placeholder="Enter your password"
						class="w-full px-4 py-3 pr-12 border {fieldErrors.password ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
					/>
					<button
						type="button"
						onclick={togglePassword}
						disabled={$auth.loading}
						class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 text-xl"
					>
						{showPassword ? '👁️' : '👁️‍🗨️'}
					</button>
				</div>
				{#if fieldErrors.password}
					<p class="mt-1 text-sm text-red-600">{fieldErrors.password}</p>
				{/if}
			</div>

			<!-- Remember Me & Forgot Password -->
			<div class="flex items-center justify-between">
				<label class="flex items-center">
					<input
						type="checkbox"
						bind:checked={rememberMe}
						disabled={$auth.loading}
						class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
					/>
					<span class="ml-2 text-sm text-gray-700">Remember me</span>
				</label>
				<a href="/auth/forgot-password" class="text-sm text-blue-600 hover:text-blue-700 font-medium">
					Forgot password?
				</a>
			</div>

			<!-- Submit Button -->
			<button
				type="submit"
				disabled={$auth.loading}
				class="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-lg hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-100"
			>
				{$auth.loading ? 'Signing in...' : 'Sign In'}
			</button>
		</form>

			<!-- Divider -->
			<div class="relative my-6">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-gray-300"></div>
				</div>
				<div class="relative flex justify-center text-sm">
					<span class="px-4 bg-white text-gray-500">Or continue with</span>
				</div>
			</div>

			<!-- Google OAuth Button -->
			<button
				type="button"
				onclick={handleGoogleLogin}
				disabled={$auth.loading}
				class="w-full py-3 px-4 border-2 border-gray-300 rounded-lg font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-3"
			>
				<svg class="w-5 h-5" viewBox="0 0 24 24">
					<path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
					<path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
					<path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
					<path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
				</svg>
				Sign in with Google
			</button>

			<!-- Register Link -->
			<div class="mt-6 text-center">
				<p class="text-gray-600">
					Don't have an account?
					<a href="/auth/register" class="text-blue-600 hover:text-blue-700 font-semibold">
						Sign up
					</a>
				</p>
			</div>
		</div>

		<!-- Back to Home -->
		<div class="text-center mt-6">
			<a href="/" class="text-white hover:text-gray-200 text-sm font-medium">
				← Back to Home
			</a>
		</div>
	</div>
</div>
