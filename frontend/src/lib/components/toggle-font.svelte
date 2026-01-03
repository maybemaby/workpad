<script lang="ts">
	import { PersistedState } from 'runed';
	import Button from './button.svelte';
	import { onMount } from 'svelte';

	let font = new PersistedState<'serif' | 'sans'>('editorFont', 'serif');

	const onclick = () => {
		document.documentElement.style.setProperty(
			'--font-tiptap',
			font.current === 'serif' ? 'var(--font-sans)' : 'var(--font-serif)'
		);

		font.current = font.current === 'serif' ? 'sans' : 'serif';
	};

	onMount(() => {
		// Set initial font on mount
		if (!font.current) {
			font.current = 'serif';
		}
		document.documentElement.style.setProperty(
			'--font-tiptap',
			font.current === 'serif' ? 'var(--font-serif)' : 'var(--font-sans)'
		);
	});
</script>

<Button variant="muted" size="icon" class="font-toggle" {onclick}
	>Ab
	<span class="sr-only">Toggle font</span>
</Button>

<style>
	:global(.font-toggle) {
		font-family: var(--font-tiptap);
		font-size: 600;
		font-weight: 700;
	}
</style>
