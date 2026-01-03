<script lang="ts">
	import { Dialog, type DialogTriggerProps } from 'bits-ui';

	let { children, ...props }: DialogTriggerProps = $props();
</script>

<Dialog.Trigger {...props} class="dialog-trigger btn muted icon">
	{@render children?.()}
</Dialog.Trigger>

<style>
	:global(.dialog-trigger) {
		padding: 0.5rem;
	}
</style>
