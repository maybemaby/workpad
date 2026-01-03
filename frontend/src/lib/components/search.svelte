<script lang="ts">
	import { Dialog, Command } from 'bits-ui';
	import DialogContent from './dialog/content.svelte';
	import { PressedKeys } from 'runed';

	const keys = new PressedKeys();

	let dialogOpen = $state(false);

	keys.onKeys(['meta', 'k'], () => {
		dialogOpen = !dialogOpen;
	});

	const items = [
		{ name: 'Project A', href: '/projects/project-a' },
		{ name: 'Project B', href: '/projects/project-b' },
		{ name: 'Project C', href: '/projects/project-c' }
	];
</script>

<Dialog.Root bind:open={dialogOpen}>
	<Dialog.Trigger class="btn muted icon">
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
			class="lucide lucide-search"
			><path d="m21 21-4.34-4.34"></path><circle cx="11" cy="11" r="8"></circle></svg
		>
		<span class="sr-only">Search</span>
	</Dialog.Trigger>

	<DialogContent>
		<Dialog.Title class="sr-only">Search Projects</Dialog.Title>
		<Dialog.Description class="sr-only">Search projects by name.</Dialog.Description>
		<Command.Root class="search__root">
			<Command.Input class="search__input" placeholder="Search projects..."></Command.Input>
			<Command.List class="search__list">
				<Command.Viewport>
					<Command.Group>
						<Command.GroupHeading>Projects</Command.GroupHeading>
						<Command.GroupItems>
							{#each items as item (item.href)}
								<Command.LinkItem href={item.href} onSelect={() => (dialogOpen = false)}>
									{item.name}
								</Command.LinkItem>
							{/each}
						</Command.GroupItems>
					</Command.Group>
				</Command.Viewport>
			</Command.List>
		</Command.Root>
	</DialogContent>
</Dialog.Root>

<style>
	:global {
		.search__root {
			--surface: var(--background);
			width: 100%;
		}

		.search__input {
			padding: 0.75rem 1rem;
			font-size: 1rem;
			width: 100%;
			box-sizing: border-box;
			border: none;
			outline: none;
			border-bottom: 1px solid var(--border);
			border-radius: 0.5rem 0.5rem 0 0;
			background: var(--surface);
		}

		.search__list {
			height: 400px;
			padding: 0.5rem 1rem;
			overflow-y: auto;
			background: var(--surface);
			border-radius: 0 0 0.5rem 0.5rem;
		}

		.search__list [data-command-group-heading] {
			margin-bottom: 0.5rem;
			font-weight: bold;
			font-size: 0.875rem;
			color: var(--muted-foreground);
		}

		.search__list [data-command-group-items] {
			display: flex;
			flex-direction: column;
		}

		.search__list [data-command-item] {
			padding: 0.5rem;
			cursor: pointer;
		}

		.search__list [data-command-item][data-selected] {
			background: hsl(var(--primary-bg));
			color: var(--primary-fg);
		}
	}
</style>
