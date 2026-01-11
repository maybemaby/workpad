<script lang="ts">
	import { getLocalTimeZone, now, parseDate } from '@internationalized/date';
	import type { Page } from '@sveltejs/kit';
	import { Portal } from 'bits-ui';

	let {
		id,
		path,
		onSelected,
		open = $bindable(false)
	}: {
		id: Page['route']['id'];
		path: string;
		onSelected?: (dateString: string) => void;
		open?: boolean;
	} = $props();

	let switchableRoute = $derived(id === '/' || id === '/dates/[date]');
	let isDateRoute = $derived(id === '/dates/[date]');
	let defaultDate = $derived.by(() => {
		if (isDateRoute) {
			return parseDate(path.split('/')[2]);
		}
		return now(getLocalTimeZone());
	});

	let selectedDate = $state(defaultDate);

	function handleKeyDown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'j' && switchableRoute) {
			e.preventDefault();

			if (open) {
				selectedDate = defaultDate;
			}

			open = !open;
		}

		if (e.key === 'Escape' && open) {
			selectedDate = defaultDate;
			open = false;
		}

		if (e.key === 'ArrowRight' && (e.metaKey || e.ctrlKey) && open) {
			e.preventDefault();
			selectedDate = selectedDate.add({ days: 1 });
		} else if (e.key === 'ArrowLeft' && (e.metaKey || e.ctrlKey) && open) {
			e.preventDefault();
			selectedDate = selectedDate.subtract({ days: 1 });
		} else if (e.key === 'ArrowUp' && (e.metaKey || e.ctrlKey) && open) {
			e.preventDefault();
			selectedDate = selectedDate.subtract({ months: 1 });
		} else if (e.key === 'ArrowDown' && (e.metaKey || e.ctrlKey) && open) {
			e.preventDefault();
			selectedDate = selectedDate.add({ months: 1 });
		}

		if (e.key === 'Enter' && open) {
			e.preventDefault();
			const dateString = `${selectedDate.year}-${String(selectedDate.month).padStart(2, '0')}-${String(
				selectedDate.day
			).padStart(2, '0')}`;
			onSelected?.(dateString);
		}
	}
</script>

<svelte:window onkeydown={handleKeyDown} />
{#if open}
	<Portal to="body">
		<div class="container" role="dialog" aria-modal="true" aria-label="Quick Switcher" data-state-open={open}>
			<div class="label">Go to:</div>
			<div>
				{selectedDate.year}-{String(selectedDate.month).padStart(2, '0')}-{String(
					selectedDate.day
				).padStart(2, '0')}
			</div>
		</div>
	</Portal>
{/if}

<style>
	@keyframes entrance {
		from {
			opacity: 0;
			transform: translateY(1rem) scale(0.85);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
	}

	.container {
		position: absolute;
		bottom: 4rem;
		right: 4rem;
		background-color: var(--surface);
		animation: entrance 0.2s ease-out;
		padding: 1rem;
		border-radius: 0.5rem;
		border: 1px solid var(--border);
		font-weight: 600;
		color: hsl(var(--primary-bg));
	}

	.label {
		font-size: 0.875rem;
		color: var(--muted-foreground);
		margin-bottom: 0.25rem;
	}
</style>
