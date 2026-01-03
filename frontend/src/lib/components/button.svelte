<script lang="ts">
	import { type HTMLButtonAttributes } from 'svelte/elements';

	let {
		class: classes,
		variant = 'primary',
		children,
		size = 'default',
		...rest
	}: HTMLButtonAttributes & {
		variant?: 'primary' | 'muted';
		size?: 'icon' | 'default';
	} = $props();
</script>

<button
	class={[{ btn: true, muted: variant === 'muted', icon: size === 'icon' }, classes]}
	{...rest}
>
	{@render children?.()}
</button>

<style>
</style>
