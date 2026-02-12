<script lang="ts">
	import { page } from '$app/stores';
	import { authStore, isAuthenticated, logout } from '$lib/stores';

	let currentPath = $derived($page.url.pathname);

	function handleLogout() {
		logout();
		window.location.href = '/auth/login';
	}
</script>

<nav class="bg-white border-b shadow-sm sticky top-0 z-50">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex justify-between items-center h-16">
			<!-- Logo / Brand -->
			<div class="flex items-center gap-8">
				<a href="/" class="text-2xl font-bold text-blue-600 hover:text-blue-700 transition">
					📚 Book Reader
				</a>

				<!-- Main Navigation Links -->
				<div class="hidden md:flex items-center gap-4">
					<a
						href="/"
						class="px-4 py-2 rounded-md font-medium transition {currentPath === '/'
							? 'bg-blue-100 text-blue-700'
							: 'text-gray-700 hover:bg-gray-100'}"
					>
						Home
					</a>
					<a
						href="/search"
						class="px-4 py-2 rounded-md font-medium transition {currentPath === '/search'
							? 'bg-blue-100 text-blue-700'
							: 'text-gray-700 hover:bg-gray-100'}"
					>
						Search
					</a>
				</div>
			</div>

			<!-- Auth Section -->
			<div class="flex items-center gap-4">
				{#if $isAuthenticated && $authStore.user}
					<div class="flex items-center gap-4">
						<span class="text-gray-700 font-medium hidden sm:block">
							{$authStore.user.username}
						</span>
						<button
							onclick={handleLogout}
							class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition font-medium"
						>
							Logout
						</button>
					</div>
				{:else}
					<div class="flex items-center gap-2">
						<a
							href="/auth/login"
							class="px-4 py-2 text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 transition font-medium"
						>
							Login
						</a>
						<a
							href="/auth/register"
							class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition font-medium"
						>
							Register
						</a>
					</div>
				{/if}
			</div>
		</div>

		<!-- Mobile Navigation -->
		<div class="md:hidden pb-4 flex items-center gap-2">
			<a
				href="/"
				class="px-4 py-2 rounded-md text-sm font-medium transition {currentPath === '/'
					? 'bg-blue-100 text-blue-700'
					: 'text-gray-700 hover:bg-gray-100'}"
			>
				Home
			</a>
			<a
				href="/search"
				class="px-4 py-2 rounded-md text-sm font-medium transition {currentPath === '/search'
					? 'bg-blue-100 text-blue-700'
					: 'text-gray-700 hover:bg-gray-100'}"
			>
				Search
			</a>
		</div>
	</div>
</nav>
