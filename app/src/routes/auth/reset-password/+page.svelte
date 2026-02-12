<script lang="ts">
	import { page } from '$app/stores';
	import { updatePassword } from '$lib/api/auth';

	// Get token from URL query params
	let token = $state($page.url.searchParams.get('token') || '');
	let password = $state('');
	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let loading = $state(false);
	let error = $state<string | null>(null);
	let success = $state(false);

	let passwordsMatch = $derived(password && password === confirmPassword);
	let canSubmit = $derived(token && password && confirmPassword && passwordsMatch && !loading);

	function togglePassword(e: MouseEvent) {
		e.preventDefault();
		e.stopPropagation();
		showPassword = !showPassword;
	}

	function toggleConfirmPassword(e: MouseEvent) {
		e.preventDefault();
		e.stopPropagation();
		showConfirmPassword = !showConfirmPassword;
	}

	async function handleSubmit() {
		if (!canSubmit) return;

		loading = true;
		error = null;

		try {
			await updatePassword({ token, password });
			success = true;
		} catch (err: any) {
			error = err.message || 'Failed to reset password';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Reset Password - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center px-4 py-12">
	<div class="w-full max-w-md">
		<div class="bg-white rounded-2xl shadow-2xl p-8">
			<!-- Logo/Title -->
			<div class="text-center mb-8">
				<div class="text-5xl mb-3">🔐</div>
				<h1 class="text-3xl font-bold text-gray-900">Reset Password</h1>
				<p class="text-gray-600 mt-2">Enter your new password</p>
			</div>

			{#if !token}
				<!-- Invalid Token -->
				<div class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
					<div class="text-4xl mb-3">⚠️</div>
					<h3 class="text-lg font-semibold text-red-900 mb-2">Invalid Reset Link</h3>
					<p class="text-red-700 text-sm mb-4">
						This password reset link is invalid or has expired.
					</p>
				</div>

				<div class="mt-6">
					<a
						href="/auth/forgot-password"
						class="block w-full py-3 px-4 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 text-center transition"
					>
						Request New Link
					</a>
				</div>

			{:else if success}
				<!-- Success Message -->
				<div class="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
					<div class="text-4xl mb-3">✅</div>
					<h3 class="text-lg font-semibold text-green-900 mb-2">Password Reset Successful!</h3>
					<p class="text-green-700 text-sm mb-4">
						Your password has been successfully reset. You can now log in with your new password.
					</p>
				</div>

				<div class="mt-6">
					<a
						href="/auth/login"
						class="block w-full py-3 px-4 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 text-center transition"
					>
						Go to Login
					</a>
				</div>

			{:else}
				<!-- Error Message -->
				{#if error}
					<div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6 flex items-start gap-3">
						<span class="text-red-600 text-xl">⚠️</span>
						<div class="flex-1">
							<p class="text-red-800 text-sm">{error}</p>
						</div>
						<button
							onclick={() => error = null}
							class="text-red-600 hover:text-red-800 text-xl leading-none"
						>
							×
						</button>
					</div>
				{/if}

				<!-- Form -->
				<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-5">
					<!-- New Password -->
					<div>
						<label for="password" class="block text-sm font-medium text-gray-700 mb-2">
							New Password
						</label>
						<div class="relative">
							<input
								id="password"
								type={showPassword ? 'text' : 'password'}
								bind:value={password}
								required
								disabled={loading}
								placeholder="Enter new password"
								minlength="8"
								class="w-full px-4 py-3 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
							/>
							<button
								type="button"
							onclick={togglePassword}
								disabled={loading}
								class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 text-xl"
							>
								{showPassword ? '👁️' : '👁️‍🗨️'}
							</button>
						</div>
						<p class="mt-1 text-xs text-gray-500">
							Must be at least 8 characters
						</p>
					</div>

					<!-- Confirm Password -->
					<div>
						<label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
							Confirm New Password
						</label>
						<div class="relative">
							<input
								id="confirmPassword"
								type={showConfirmPassword ? 'text' : 'password'}
								bind:value={confirmPassword}
								required
								disabled={loading}
								placeholder="Re-enter new password"
								class="w-full px-4 py-3 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
							/>
							<button
								type="button"
								onclick={toggleConfirmPassword}
								disabled={loading}
								class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 text-xl"
							>
								{showConfirmPassword ? '👁️' : '👁️‍🗨️'}
							</button>
						</div>
						
						{#if confirmPassword && password}
							<p class="mt-2 text-sm {passwordsMatch ? 'text-green-600' : 'text-red-600'}">
								{passwordsMatch ? '✓ Passwords match' : '✗ Passwords do not match'}
							</p>
						{/if}
					</div>

					<button
						type="submit"
						disabled={!canSubmit}
						class="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-lg hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-100"
					>
						{loading ? 'Resetting...' : 'Reset Password'}
					</button>
				</form>
			{/if}
		</div>

		<!-- Back to Home -->
		<div class="text-center mt-6">
			<a href="/" class="text-white hover:text-gray-200 text-sm font-medium">
				← Back to Home
			</a>
		</div>
	</div>
</div>
