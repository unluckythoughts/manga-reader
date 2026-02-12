<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';

	onMount(async () => {
		try {
			await auth.logout();
		} catch (error) {
			console.error('Logout failed:', error);
		} finally {
			// Always redirect to home after logout attempt
			window.location.href = '/';
		}
	});
</script>

<svelte:head>
	<title>Logging out...</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex items-center justify-center">
	<div class="text-center">
		<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
		<p class="mt-4 text-gray-600">Logging you out...</p>
	</div>
</div>
