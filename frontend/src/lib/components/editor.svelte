<script lang="ts">
	import { Editor, type Content } from '@tiptap/core';
	import { StarterKit } from '@tiptap/starter-kit';
	import { TaskList, TaskItem } from '@tiptap/extension-list';
	import { onMount, untrack } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';
	import { Debounced } from 'runed';
	import { ProjectTag, findMentionParents, type MentionNodes } from '$lib/editor/project';

	let {
		editable = true,
		content,
		onUpdate
	}: {
		editable?: boolean;
		content?: Content;
		onUpdate?: (data: { html: string; mentionNodes: MentionNodes[] }) => void;
	} = $props();
	let el = $state<HTMLElement>();
	let editor = $derived<Editor | null>(content ? null : null);
	let editedOnce = $state(false);

	let lastUpdate = new SvelteDate();
	let debouncedUpdate = new Debounced(() => lastUpdate.getTime(), 800);

	$effect(() => {
		(() => debouncedUpdate.current)();

		untrack(() => {
			if (editor && editedOnce) {
				// const structure = editor.getJSON();

				// Find all parent nodes that contain mentions
				// const mentionNodes = findMentionParents(structure);

				const html = editor.getHTML();

				const listMentions = Array.from(
					document.querySelector('.tiptap')?.querySelectorAll('li[data-checked]:has(.mention)') ||
						[]
				);

				const pMentions = Array.from(
					document.querySelector('.tiptap')?.querySelectorAll('p:has(.mention)') || []
				);

				// Dedupe mentions found in list items
				const filteredPMentions = pMentions.filter((p) => {
					return typeof listMentions.find((li) => li.contains(p)) === 'undefined';
				});

				const mentions = [...listMentions, ...filteredPMentions];

				const mentionNodes = mentions.map((el) => {
					let projects: string[] = [];

					for (let mentionEl of el.querySelectorAll('.mention')) {
						const projectName = mentionEl.getAttribute('data-mention-id');
						if (projectName) {
							projects.push(projectName);
						}
					}

					return {
						node: el.outerHTML,
						mentioned: projects
					};
				});

				console.log('Mention Nodes:', mentionNodes);

				onUpdate?.({ html, mentionNodes });
			}
		});
	});

	onMount(() => {
		editor = new Editor({
			editable,
			element: el,
			extensions: [
				StarterKit.configure({
					paragraph: {
						HTMLAttributes: {
							class: 'para-node'
						}
					}
				}),
				TaskList,
				TaskItem,
				ProjectTag
			],
			content:
				content ||
				`
      <h2>Today's Tasks</h2>
      <ul data-type="taskList">
        <li data-type="taskItem" data-checked="false">Task 1</li>
      </ul>

      <br />
      <h2>Passdown</h2>
			<p></p>
      `,
			onUpdate(props) {
				editedOnce = true;
				editor = props.editor;
				lastUpdate.setTime(Date.now());
			}
		});

		return () => {
			editor?.destroy();
		};
	});
</script>

<div class="editor-container">
	<div bind:this={el}></div>
</div>

<style>
	.editor-container {
		position: relative;
	}
</style>
