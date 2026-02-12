<script lang="ts">
	import { resetPassword } from '$lib/api/auth';

	let email = $state('');
	let loading = $state(false);
	let error = $state<string | null>(null);
	let success = $state(false);

	async function handleSubmit() {
		loading = true;
		error = null;

		try {
			await resetPassword(email);
			success = true;
		} catch (err: any) {
			error = err.message || 'Failed to send reset link';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Forgot Password - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center px-4 py-12">
	<div class="w-full max-w-md">
		<div class="bg-white rounded-2xl shadow-2xl p-8">
			<!-- Logo/Title -->
			<div class="text-center mb-8">
				<div class="text-5xl mb-3">🔑</div>
				<h1 class="text-3xl font-bold text-gray-900">Forgot Password?</h1>
				<p class="text-gray-600 mt-2">We'll send you a reset link</p>
			</div>

			{#if success}
				<!-- Success Message -->
				<div class="bg-green-50 border border-green-200 rounded-lg p-6 text-center">
					<div class="text-4xl mb-3">✅</div>
					<h3 class="text-lg font-semibold text-green-900 mb-2">Check Your Email</h3>
					<p class="text-green-700 text-sm mb-4">
						We've sent a password reset link to <strong>{email}</strong>
					</p>
					<p class="text-green-600 text-xs">
						If you don't see the email, check your spam folder.
					</p>
				</div>

				<div class="mt-6 space-y-3">
					<a
						href="/auth/login"
						class="block w-full py-3 px-4 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 text-center transition"
					>
						Back to Login
					</a>
					<button
						onclick={() => { success = false; email = ''; }}
						class="block w-full py-3 px-4 border border-gray-300 text-gray-700 font-semibold rounded-lg hover:bg-gray-50 text-center transition"
					>
						Try Another Email
					</button>
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
					<div>
						<label for="email" class="block text-sm font-medium text-gray-700 mb-2">
							Email Address
						</label>
						<input
							id="email"
							type="email"
							bind:value={email}
							required
							disabled={loading}
							placeholder="your@email.com"
							class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
						/>
						<p class="mt-2 text-xs text-gray-500">
							Enter the email address associated with your account
						</p>
					</div>

					<button
						type="submit"
						disabled={loading || !email}
						class="w-full py-3 px-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-lg hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-100"
					>
						{loading ? 'Sending...' : 'Send Reset Link'}
					</button>
				</form>

				<div class="mt-6 text-center">
					<a href="/auth/login" class="text-blue-600 hover:text-blue-700 font-medium text-sm">
						← Back to Login
					</a>
				</div>
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
