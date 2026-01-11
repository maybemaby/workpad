<script lang="ts">
	import {
		createAddNoteMutation,
		createProjectsMutation,
		createUpdateExcerptMutation,
		getNoteByDateQuery
	} from '$lib/api/queries.svelte';
	import Editor from '$lib/components/editor.svelte';
	import type { MentionNodes } from '$lib/editor/project';
	import { getLocalTimeZone, today } from '@internationalized/date';

	const updateNote = createAddNoteMutation();
	const updateExcerpts = createUpdateExcerptMutation();
	const createProjects = createProjectsMutation();

	let currentDate = $state(today(getLocalTimeZone()));

	let query = getNoteByDateQuery(currentDate.toString());

	const onUpdate = async (data: { html: string; mentionNodes: MentionNodes[] }) => {
		const updateNotesPromise = updateNote.mutateAsync({ html_content: data.html });
		const createProjectsPromise = createProjects.mutateAsync({
			projects: data.mentionNodes.flatMap((node) => node.mentioned)
		});

		await Promise.all([updateNotesPromise, createProjectsPromise]);

		await updateExcerpts.mutateAsync({
			excerpts: data.mentionNodes.map((mention) => ({
				node: mention.node,
				projects: mention.mentioned
			})),
			date: `${currentDate.year}-${String(currentDate.month).padStart(2, '0')}-${String(
				currentDate.day
			).padStart(2, '0')}`
		});
	};

	let content = $derived(query.data?.html_content ?? undefined);
</script>

{#if !query.isPending}
	{#key content}
		<Editor editable {onUpdate} {content} />
	{/key}
{/if}
