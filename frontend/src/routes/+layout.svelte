<script lang="ts">
	import { page } from '$app/state';
	import { CalendarDate } from '@internationalized/date';
	import '@fontsource-variable/playfair-display';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import favicon from '$lib/assets/favicon.svg';
	import Header from '$lib/components/header.svelte';
	import '../app.css';
	import type { LayoutProps } from './$types';
	import QuickSwitcher from '$lib/components/quick-switcher.svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	let { children }: LayoutProps = $props();

	const queryClient = new QueryClient();

	let routeId = $derived(page.route.id);

	let pageDate = $derived.by(() => {
		// Check if we're on a /dates/[date] route
		if (routeId?.startsWith('/dates/')) {
			const dateParam = page.params.date;
			if (dateParam) {
				const [year, month, day] = dateParam.split('-').map(Number);
				return new CalendarDate(year, month, day);
			}
		}

		// Default to today if not on a date page
		return undefined;
	});

	function handleSwitcherSelected(dateString: string) {
		// Navigate to the selected date
		goto(resolve(`/dates/${dateString}`));
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<div class="app-container">
		<div class="app-header__outer">
			<Header date={pageDate}></Header>
		</div>
		<main>
			{@render children()}
		</main>
		<QuickSwitcher id={routeId} path={page.url.pathname} onSelected={handleSwitcherSelected} />
	</div>
</QueryClientProvider>

<style>
	.app-header__outer {
		flex-shrink: 0;
		height: var(--header-height);
	}

	main {
		flex-grow: 1;
		min-height: var(--main-height);
		max-width: var(--screen-lg);
		margin: auto;
		padding: 1rem;
	}
</style>
