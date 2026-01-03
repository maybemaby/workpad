<script lang="ts">
	import { Dialog } from 'bits-ui';

	let { children, ...props }: Dialog.ContentProps = $props();
</script>

<Dialog.Portal>
	<Dialog.Overlay class="dialog-overlay">
		<Dialog.Content {...props} class="dialog-content">
			{@render children?.()}
		</Dialog.Content>
	</Dialog.Overlay>
</Dialog.Portal>

<style>
	:global {
		.dialog-content {
			max-width: 600px;
			width: 90vw;
			background: var(--surface);
			border-radius: 8px;
			z-index: 500;
			position: fixed;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
		}

		.dialog-overlay {
			background: hsl(0, 0%, 0%, 0.5);
			position: fixed;
			inset: 0;
			z-index: 400;
			animation: fadeIn 100ms ease-out;
		}

		.dialog-overlay[data-state='closed'] {
			animation: fadeOut 100ms ease-out;
		}
	}
</style>
