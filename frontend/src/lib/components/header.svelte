<script lang="ts">
	import { CalendarDate } from '@internationalized/date';
	import DatePicker from './date-picker.svelte';
	import ToggleFont from './toggle-font.svelte';
	import { getMonthNotes } from '$lib/api/queries.svelte';
	import Search from './search.svelte';

	let { date = $bindable() }: { date?: CalendarDate } = $props();
	let navigationDate = $state<undefined | CalendarDate>(undefined);

	let notesQuery = getMonthNotes(() => ({
		year: navigationDate?.year ?? 2026,
		month: navigationDate?.month ?? 1,
		enabled: !!navigationDate
	}));

	let daysWithNotes = $derived.by(() => {
		if (!navigationDate) return [];

		return (
			notesQuery.data?.map(
				(day) => new CalendarDate(navigationDate!.year, navigationDate!.month, day)
			) ?? []
		);
	});

	const handleNavigation = (navDate: CalendarDate) => {
		navigationDate = navDate;
	};
</script>

<header>
	<div class="placeholder"></div>
	<DatePicker {date} markDates={daysWithNotes} onNavigation={handleNavigation} />
	<div class="actions">
		<Search />
		<ToggleFont />
	</div>
</header>

<style>
	header {
		max-width: var(--screen-lg);
		margin: auto;
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
		height: 100%;
		padding: 0 1rem;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.placeholder {
		display: none;
	}

	@media (min-width: 640px) {
		.placeholder {
			display: block;
		}
	}
</style>
