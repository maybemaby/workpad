<script lang="ts">
	import { resolve } from '$app/paths';
	import { Popover, Calendar } from 'bits-ui';
	import { today, getLocalTimeZone, CalendarDate } from '@internationalized/date';

	let {
		date = $bindable(),
		markDates = [],
		onNavigation
	}: {
		date?: CalendarDate;
		markDates?: CalendarDate[];
		onNavigation?: (date: CalendarDate) => void;
	} = $props();

	const currentDate = today(getLocalTimeZone());
	let selectedDate = $derived(date || today(getLocalTimeZone()));
	let placeholderDate = $state<CalendarDate | undefined>(undefined);

	let noteDays = $derived.by(() => {
		return markDates
			.filter((d) => d.month === placeholderDate?.month && d.year === placeholderDate?.year)
			.map((d) => d.day);
	});

	$effect(() => {
		if (placeholderDate) onNavigation?.(placeholderDate);
	});
</script>

<Popover.Root>
	<Popover.Trigger
		class={{
			'date-picker-trigger': true,
			matchesToday: selectedDate.toString() === currentDate.toString()
		}}
	>
		<svg
			xmlns="http://www.w3.org/2000/svg"
			width="18"
			height="18"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
			class="lucide lucide-calendar-fold"
			><path
				d="M3 20a2 2 0 0 0 2 2h10a2.4 2.4 0 0 0 1.706-.706l3.588-3.588A2.4 2.4 0 0 0 21 16V6a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2z"
			></path><path d="M15 22v-5a1 1 0 0 1 1-1h5"></path><path d="M8 2v4"></path><path d="M16 2v4"
			></path><path d="M3 10h18"></path></svg
		>
		{selectedDate.month}/{selectedDate.day}/{selectedDate.year}
	</Popover.Trigger>

	<Popover.Portal>
		<Popover.Content sideOffset={8} class="date-picker-content popover-content">
			<Calendar.Root
				weekdayFormat="short"
				type="single"
				bind:value={selectedDate}
				bind:placeholder={placeholderDate}
			>
				{#snippet children({ months, weekdays })}
					<Calendar.Header class="calendar-header">
						<Calendar.MonthSelect aria-label="Select month" />
						<Calendar.YearSelect aria-label="Select year" />
					</Calendar.Header>
					<div class="date-picker-grid__container">
						{#each months as month, i (i)}
							<Calendar.Grid class="date-picker-grid__grid">
								<Calendar.GridHead>
									<Calendar.GridRow class="date-picker-grid__head_row">
										{#each weekdays as day, i (i)}
											<Calendar.HeadCell class="date-picker-grid__head_cell">
												<div>{day.slice(0, 2)}</div>
											</Calendar.HeadCell>
										{/each}
									</Calendar.GridRow>
								</Calendar.GridHead>
								<Calendar.GridBody>
									{#each month.weeks as weekDates, i (i)}
										<Calendar.GridRow class="date-picker-grid__body_row">
											{#each weekDates as date, i (i)}
												<Calendar.Cell
													{date}
													month={month.value}
													class="date-picker-grid__body_cell"
												>
													<a
														href={resolve(
															`/dates/${date.year}-${String(date.month).padStart(2, '0')}-${String(date.day).padStart(2, '0')}`
														)}
														data-sveltekit-preload-data="tap"
													>
														<Calendar.Day class="date-picker-grid__day">
															{#if selectedDate.day == date.day && selectedDate.month == date.month && selectedDate.year == date.year}
																<div class="date-picker-grid__day_indicator"></div>
																<!-- Marked dates -->
															{:else if noteDays.includes(date.day)}
																<div class="date-picker-grid__day_indicator noted"></div>
															{/if}
															{date.day}
														</Calendar.Day>
													</a>
												</Calendar.Cell>
											{/each}
										</Calendar.GridRow>
									{/each}
								</Calendar.GridBody>
							</Calendar.Grid>
						{/each}
					</div>
				{/snippet}
			</Calendar.Root>
			<hr />
			<a class="btn" href="/">Go to Today</a>
		</Popover.Content>
	</Popover.Portal>
</Popover.Root>

<style>
	.calendar-next {
		transform: rotate(180deg);
	}

	hr {
		margin: 1rem 0;
		border: none;
		border-top: 1px solid var(--border);
	}

	:global {
		.date-picker-trigger {
			font-family: var(--font-mono);
			cursor: pointer;
			background: hsl(var(--neutral) / 0.1);
			padding: 0.5rem 1rem;
			border-radius: 0.5rem;
			border: 1px solid var(--border);
			font-size: 1.05rem;
			display: inline-flex;
			align-items: center;
			gap: 0.5rem;
			transition: background 0.2s ease;
		}

		.date-picker-trigger.matchesToday {
			color: hsl(var(--primary-bg));
		}

		.date-picker-trigger:hover,
		.date-picker-trigger:focus {
			background: hsl(var(--neutral) / 0.3);
		}

		.date-picker-content {
			padding: 1rem;
			border-radius: 0.5rem;
			border: 1px solid var(--border);
			background: var(--background);
		}

		.calendar-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 0.5rem;
			gap: 0.5rem;
			min-width: 25ch;

			button {
				background: transparent;
				border: none;
				cursor: pointer;
				font-size: 1.25rem;
			}
		}

		.calendar-header [data-calendar-month-select],
		.calendar-header [data-calendar-year-select] {
			background-color: transparent;
			border: none;
			font-size: 1.1rem;
		}

		.date-picker-grid__container {
			display: flex;
			flex-direction: column;
			padding-top: 1rem;
			gap: 1rem;
		}

		.date-picker-grid__grid {
			--selected: hsl(var(--primary-bg));
			--selected-hover: hsl(var(--primary-bg) / 0.8);
			width: 100%;
			border-collapse: collapse;
			user-select: none;
			-moz-user-select: none;
			-webkit-user-select: none;
		}

		.date-picker-grid__head_row {
			margin-bottom: 0.25rem;
			display: flex;
			justify-content: space-between;
			width: 100%;
		}

		.date-picker-grid__head_cell {
			color: var(--muted-foreground);
			font-size: normal;
			width: 2.5rem;
			border-radius: 0.375rem;
			font-size: 0.75rem;
		}

		.date-picker-grid__body_row {
			display: flex;
			width: 100%;
		}

		.date-picker-grid__body_cell {
			width: 2.5rem;
			height: 2.5rem;
			position: relative;
			text-align: center;
			font-size: 0.875rem;
		}

		.date-picker-grid__day {
			border-radius: 0.5625rem; /* 9px */
			color: var(--foreground);
			position: relative;
			display: inline-flex;
			width: 2.5rem;
			height: 2.5rem;
			align-items: center;
			justify-content: center;
			white-space: nowrap;
			border: 1px solid transparent;
			background: transparent;
			padding: 0;
			font-size: 0.875rem;
			font-weight: normal;
		}
		.date-picker-grid__day:hover {
			background-color: var(--selected-hover);
			color: var(--primary-fg);
		}

		.date-picker-grid__day[data-selected] {
			background-color: var(--selected);
			color: var(--background);
			font-weight: 500;
		}
		.date-picker-grid__day[data-disabled] {
			color: color-mix(in srgb, var(--foreground) 30%, transparent);
			pointer-events: none;
		}
		.date-picker-grid__day[data-unavailable] {
			color: var(--muted-foreground);
			text-decoration: line-through;
		}
		.date-picker-grid__day[data-outside-month] {
			pointer-events: none;
		}
		/* Dot indicator */
		.date-picker-grid__day_indicator {
			position: absolute;
			top: 0.3125rem; /* 5px */
			width: 0.25rem;
			height: 0.25rem;
			border-radius: 9999px;
			background-color: var(--foreground);
			display: none;
		}

		.date-picker-grid__day_indicator.noted {
			background-color: hsl(var(--primary-bg));
			display: block;
		}

		.date-picker-grid__day[data-selected] .date-picker-grid__day_indicator {
			display: block;
			background-color: var(--background);
		}
		.date-picker-grid__day[data-today] .date-picker-grid__day_indicator {
			display: block;
		}
	}
</style>
