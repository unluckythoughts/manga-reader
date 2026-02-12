<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Navigation from '$lib/components/Navigation.svelte';
	import { authStore, isAuthenticated } from '$lib/stores';

	let { children } = $props();
	
	// Check if we're on an auth page
	let isAuthPage = $derived($page.url.pathname.startsWith('/auth'));
	
	// Initialize auth store on mount
	onMount(() => {
		// Auth store will automatically check for stored tokens
		// No need to explicitly initialize
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
