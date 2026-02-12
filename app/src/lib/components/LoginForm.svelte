<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import type { LoginRequest } from '$lib/types';

	let formData = $state<LoginRequest>({
		username: '',
		password: ''
	});

	let showPassword = $state(false);

	async function handleSubmit() {
		try {
			await auth.login(formData.username, formData.password);
			// Redirect to home or dashboard
			window.location.href = '/';
		} catch (error) {
			// Error is already set in the store
			console.error('Login failed:', error);
		}
	}
</script>

<div class="login-container">
	<div class="login-card">
		<h1>Login</h1>

		{#if $auth.error}
			<div class="error-message">
				{$auth.error}
				<button onclick={() => auth.clearError()}>✕</button>
			</div>
		{/if}

		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
			<div class="form-group">
				<label for="username">Username or Email</label>
				<input
					id="username"
					type="text"
					bind:value={formData.username}
					required
					disabled={$auth.loading}
					placeholder="Enter your username or email"
				/>
			</div>

			<div class="form-group">
				<label for="password">Password</label>
				<div class="password-input">
					<input
						id="password"
						type={showPassword ? 'text' : 'password'}
						bind:value={formData.password}
						required
						disabled={$auth.loading}
						placeholder="Enter your password"
					/>
					<button
						type="button"
						class="toggle-password"
						onclick={() => showPassword = !showPassword}
						disabled={$auth.loading}
					>
						{showPassword ? '👁️' : '👁️‍🗨️'}
					</button>
				</div>
			</div>

			<button type="submit" class="submit-btn" disabled={$auth.loading}>
				{$auth.loading ? 'Logging in...' : 'Login'}
			</button>
		</form>

		<div class="footer-links">
			<a href="/auth/register">Don't have an account? Register</a>
			<a href="/auth/forgot-password">Forgot password?</a>
		</div>
	</div>
</div>

<style>
	.login-container {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	}

	.login-card {
		background: white;
		padding: 2rem;
		border-radius: 12px;
		box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
		width: 100%;
		max-width: 400px;
	}

	h1 {
		margin: 0 0 1.5rem 0;
		text-align: center;
		color: #333;
	}

	.error-message {
		background: #ffebee;
		color: #c62828;
		padding: 0.75rem;
		border-radius: 6px;
		margin-bottom: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.error-message button {
		background: none;
		border: none;
		color: #c62828;
		cursor: pointer;
		font-size: 1.2rem;
	}

	.form-group {
		margin-bottom: 1.5rem;
	}

	label {
		display: block;
		margin-bottom: 0.5rem;
		color: #555;
		font-weight: 500;
	}

	input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 6px;
		font-size: 1rem;
		transition: border-color 0.2s;
	}

	input:focus {
		outline: none;
		border-color: #667eea;
	}

	input:disabled {
		background: #f5f5f5;
		cursor: not-allowed;
	}

	.password-input {
		position: relative;
	}

	.toggle-password {
		position: absolute;
		right: 0.5rem;
		top: 50%;
		transform: translateY(-50%);
		background: none;
		border: none;
		cursor: pointer;
		font-size: 1.2rem;
		padding: 0.25rem;
	}

	.submit-btn {
		width: 100%;
		padding: 0.75rem;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		color: white;
		border: none;
		border-radius: 6px;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: transform 0.2s, opacity 0.2s;
	}

	.submit-btn:hover:not(:disabled) {
		transform: translateY(-2px);
	}

	.submit-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.footer-links {
		margin-top: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		text-align: center;
	}

	.footer-links a {
		color: #667eea;
		text-decoration: none;
		font-size: 0.875rem;
	}

	.footer-links a:hover {
		text-decoration: underline;
	}
</style>
