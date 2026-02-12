<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import type { RegisterRequest } from '$lib/types';

	let formData = $state<RegisterRequest>({
		username: '',
		email: '',
		password: ''
	});

	let confirmPassword = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let acceptTerms = $state(false);
	let touched = $state({
		username: false,
		email: false,
		password: false,
		confirmPassword: false,
		terms: false
	});

	// Field validation errors
	let fieldErrors = $derived.by(() => {
		const errors: Record<string, string> = {};
		
		if (touched.username && !formData.username) {
			errors.username = 'Username is required';
		} else if (touched.username && formData.username.length < 3) {
			errors.username = 'Username must be at least 3 characters';
		}
		
		if (touched.email && !formData.email) {
			errors.email = 'Email is required';
		} else if (touched.email && formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
			errors.email = 'Please enter a valid email address';
		}
		
		if (touched.password && !formData.password) {
			errors.password = 'Password is required';
		} else if (touched.password && formData.password && formData.password.length < 8) {
			errors.password = 'Password must be at least 8 characters';
		}
		
		if (touched.confirmPassword && !confirmPassword) {
			errors.confirmPassword = 'Please confirm your password';
		} else if (touched.confirmPassword && confirmPassword && formData.password !== confirmPassword) {
			errors.confirmPassword = 'Passwords do not match';
		}
		
		if (touched.terms && !acceptTerms) {
			errors.terms = 'You must agree to the Terms of Service';
		}
		
		return errors;
	});

	let passwordStrength = $derived.by(() => {
		const pwd = formData.password;
		if (!pwd) return { score: 0, text: '', color: '' };
		
		let score = 0;
		if (pwd.length >= 8) score++;
		if (pwd.length >= 12) score++;
		if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score++;
		if (/\d/.test(pwd)) score++;
		if (/[^A-Za-z0-9]/.test(pwd)) score++;

		if (score <= 2) return { score, text: 'Weak', color: 'bg-red-500' };
		if (score <= 3) return { score, text: 'Fair', color: 'bg-yellow-500' };
		if (score <= 4) return { score, text: 'Good', color: 'bg-blue-500' };
		return { score, text: 'Strong', color: 'bg-green-500' };
	});

	let passwordsMatch = $derived(formData.password && formData.password === confirmPassword);
	let canSubmit = $derived(
		!!formData.username &&
		formData.username.length >= 3 &&
		!!formData.email &&
		/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email) &&
		!!formData.password &&
		formData.password.length >= 8 &&
		!!confirmPassword &&
		passwordsMatch &&
		acceptTerms &&
		!$auth.loading
	);

	function markTouched(field: keyof typeof touched) {
		touched[field] = true;
	}

	function togglePassword() {
		showPassword = !showPassword;
	}

	function toggleConfirmPassword() {
		showConfirmPassword = !showConfirmPassword;
	}

	async function handleSubmit() {
		// Mark all fields as touched to show validation errors
		touched.username = true;
		touched.email = true;
		touched.password = true;
		touched.confirmPassword = true;
		touched.terms = true;

		if (!canSubmit) return;

		try {
			await auth.register(formData.username, formData.email, formData.password);
			// Redirect to home
			window.location.href = '/';
		} catch (error) {
			// Error is already set in the store
			console.error('Registration failed:', error);
		}
	}
</script>

<svelte:head>
	<title>Register - Book Reader</title>
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center px-4 py-12">
	<div class="w-full max-w-md">
		<div class="bg-white rounded-2xl shadow-2xl p-8">
			<!-- Logo/Title -->
			<div class="text-center mb-8">
				<div class="text-5xl mb-3">📚</div>
				<h1 class="text-3xl font-bold text-gray-900">Create Account</h1>
				<p class="text-gray-600 mt-2">Join us and start reading</p>
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

			<!-- Register Form -->
			<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-5">
				<!-- Username -->
				<div>
					<label for="username" class="block text-sm font-medium text-gray-700 mb-2">
						Username
					</label>
					<input
						id="username"
						type="text"
						bind:value={formData.username}
						onblur={() => markTouched('username')}
						required
						disabled={$auth.loading}
						placeholder="Choose a username"
						class="w-full px-4 py-3 border {fieldErrors.username ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
					/>
					{#if fieldErrors.username}
						<p class="mt-1 text-sm text-red-600">{fieldErrors.username}</p>
					{/if}
				</div>

				<!-- Email -->
				<div>
					<label for="email" class="block text-sm font-medium text-gray-700 mb-2">
						Email Address
					</label>
					<input
						id="email"
						type="email"
						bind:value={formData.email}
						onblur={() => markTouched('email')}
						required
						disabled={$auth.loading}
						placeholder="your@email.com"
						class="w-full px-4 py-3 border {fieldErrors.email ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
					/>
					{#if fieldErrors.email}
						<p class="mt-1 text-sm text-red-600">{fieldErrors.email}</p>
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
							placeholder="Create a strong password"
							class="w-full px-4 py-3 pr-12 border {fieldErrors.password ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
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
					
					<!-- Password Strength Indicator -->
					{#if formData.password}
						<div class="mt-2">
							<div class="flex items-center gap-2">
								<div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
									<div
										class="{passwordStrength.color} h-full transition-all duration-300"
										style="width: {(passwordStrength.score / 5) * 100}%"
									></div>
								</div>
								<span class="text-xs font-medium text-gray-600">
									{passwordStrength.text}
								</span>
							</div>
						</div>
					{/if}
				</div>

				<!-- Confirm Password -->
				<div>
					<label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
						Confirm Password
					</label>
					<div class="relative">
						<input
							id="confirmPassword"
							type={showConfirmPassword ? 'text' : 'password'}
							bind:value={confirmPassword}
							onblur={() => markTouched('confirmPassword')}
							required
							disabled={$auth.loading}
							placeholder="Re-enter your password"
							class="w-full px-4 py-3 pr-12 border {fieldErrors.confirmPassword ? 'border-red-500' : 'border-gray-300'} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition disabled:bg-gray-100 disabled:cursor-not-allowed"
						/>
						<button
							type="button"
							onclick={toggleConfirmPassword}
							disabled={$auth.loading}
							class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 text-xl"
						>
							{showConfirmPassword ? '👁️' : '👁️‍🗨️'}
						</button>
					</div>
					
					{#if fieldErrors.confirmPassword}
						<p class="mt-1 text-sm text-red-600">{fieldErrors.confirmPassword}</p>
					{:else if confirmPassword && formData.password && passwordsMatch}
						<p class="mt-1 text-sm text-green-600">✓ Passwords match</p>
					{/if}
				</div>

				<!-- Terms & Conditions -->
				<div>
					<div class="flex items-start">
						<input
							id="terms"
							type="checkbox"
							bind:checked={acceptTerms}
							onchange={() => markTouched('terms')}
							disabled={$auth.loading}
							class="w-4 h-4 text-purple-600 border-gray-300 rounded focus:ring-purple-500 mt-1"
						/>
						<label for="terms" class="ml-2 text-sm text-gray-700">
							I agree to the
							<a href="/terms" class="text-purple-600 hover:text-purple-700 font-medium">Terms of Service</a>
							and
							<a href="/privacy" class="text-purple-600 hover:text-purple-700 font-medium">Privacy Policy</a>
						</label>
					</div>
					{#if fieldErrors.terms}
						<p class="mt-1 text-sm text-red-600">{fieldErrors.terms}</p>
					{/if}
				</div>
				<!-- Submit Button -->
				<button
					type="submit"
					disabled={$auth.loading}
					class="w-full py-3 px-4 bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold rounded-lg hover:from-purple-700 hover:to-pink-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-100"
				>
					{$auth.loading ? 'Creating Account...' : 'Create Account'}
				</button>
			</form>

			<!-- Login Link -->
			<div class="mt-6 text-center">
				<p class="text-gray-600">
					Already have an account?
					<a href="/auth/login" class="text-purple-600 hover:text-purple-700 font-semibold">
						Sign in
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
