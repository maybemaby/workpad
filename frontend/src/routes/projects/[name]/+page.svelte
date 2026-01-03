<script lang="ts">
	import { getExcerptsQuery } from '$lib/api/queries.svelte';
	import type { PageProps } from './$types';

	let { params }: PageProps = $props();

	let name = $derived(params.name);

	const query = getExcerptsQuery(() => name);

	const projectName = $derived.by(() => {
		return query.data?.[0]?.project_name ?? name;
	});

	const groupedExcerpts = $derived.by(() => {
		return Object.groupBy(query.data ?? [], (excerpt) => excerpt.date.slice(0, 10));
	});
</script>

<h1>{projectName}</h1>
{#if query.data}
	<div class="tiptap project-excerpts">
		{#each Object.entries(groupedExcerpts) as [date, excerpts] (date)}
			<div class="group">
				<h2>{date}</h2>
				<div class="items">
					{#each excerpts as excerpt (excerpt.id)}
						{@html excerpt.excerpt}
					{/each}
				</div>
			</div>
		{/each}
	</div>
{/if}

<style>
	:global {
		.project-excerpts li {
			display: flex;
			align-items: center;
			gap: 0.5em;
		}
	}

	h1 {
		margin-bottom: 2.5rem;
		font-family: var(--font-tiptap);
	}

	.group {
		margin-bottom: 2rem;
		display: flex;
	}

	.items {
		flex-grow: 1;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		margin-left: 3rem;
	}

	.item {
		margin-bottom: 1.5rem;
		display: flex;
		gap: 1rem;
	}
</style>
