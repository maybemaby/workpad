<script lang="ts">
	import Editor from '$lib/components/editor.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let editable = $derived.by(() => {
		return (
			data.year === data.currentDate.year &&
			data.month === data.currentDate.month &&
			data.day === data.currentDate.day
		);
	});

	let content = $derived(data.note?.html_content ?? 'No Content');
</script>

{#key data}
	<Editor {editable} {content} />
{/key}
