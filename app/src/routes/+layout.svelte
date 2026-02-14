<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Navigation from '$lib/components/Navigation.svelte';
	import { auth } from '$lib/stores';

	let { children } = $props();
	
	// Check if we're on an auth page
	let isAuthPage = $derived($page.url.pathname.startsWith('/auth'));
	
	// Check if auth is disabled
	const authDisabled = import.meta.env.VITE_ENABLE_AUTH === 'false';
	
	// Initialize auth store on mount
	onMount(async () => {
		// Try to load current user session
		const user = await auth.init();
		
		// If auth is disabled and no user is logged in, auto-login with dummy credentials
		if (authDisabled && !user) {
			try {
				await auth.login('dummy@example.com', 'dummy');
				// Redirect to home page if currently on auth page
				if (isAuthPage) {
					goto('/');
				}
			} catch (error) {
				console.error('Auto-login failed:', error);
			}
		}
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="min-h-screen bg-gray-50">
	{#if !isAuthPage}
		<Navigation />
	{/if}
	
	<main>
		{@render children()}
	</main>
</div>
