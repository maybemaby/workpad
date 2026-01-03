import { apiClient } from '$lib/api/client';
import { debounce } from '$lib/utils/debounce';
import { computePosition, offset } from '@floating-ui/dom';
import { mergeAttributes, type JSONContent } from '@tiptap/core';
import Mention, { type MentionNodeAttrs } from '@tiptap/extension-mention';

class ToolTip {
	el: HTMLElement;
	selectedIndex: number = 0;

	constructor(public command: (props: MentionNodeAttrs) => void) {
		this.el = document.createElement('div');
		this.el.className = 'mention-tooltip';
		this.el.style.position = 'absolute';
		this.el.innerHTML = 'Mention Tooltip';
	}

	mount() {
		document.body.appendChild(this.el);
	}

	remove() {
		this.el.remove();
	}

	async updatePosition(rootEl: Element) {
		await computePosition(rootEl, this.el, {
			placement: 'top-end',
			middleware: [
				offset({
					alignmentAxis: -120
				})
			]
		}).then(({ x, y }) => {
			this.el.style.left = `${x}px`;
			this.el.style.top = `${y}px`;
		});
	}

	async setOptions(options: string[]) {
		this.el.innerHTML = options.map((option) => `<button>${option}</button>`).join('');

		const children = this.el.children;

		for (let i = 0; i < children.length; i++) {
			// TODO: Remove previous listeners to avoid multiple triggers
			children[i].addEventListener('click', () => {
				this.command?.({
					id: options[i],
					label: options[i]
				});
			});
		}

		this.toggleHighlight(this.selectedIndex);
	}

	setSelectedIndex(index: number) {
		const children = this.el.children;

		if (index < 0 || index >= children.length) {
			return;
		}

		this.toggleHighlight(this.selectedIndex);

		this.selectedIndex = index;

		this.toggleHighlight(index);
	}

	toggleHighlight(index: number) {
		const children = this.el.children;

		children[index].classList.toggle('highlighted');
	}

	selectCurrent() {
		const children = this.el.children;
		if (children[this.selectedIndex]) {
			(children[this.selectedIndex] as HTMLButtonElement).click();
		}
	}
}

// Cache to store the most recent project results
const projectCache: Map<string, string[]> = new Map();

// Debounced API call to fetch projects and update cache
const debouncedProjectQuery = debounce(async (query: string) => {
	const res = await apiClient.GET('/api/projects', {
		params: {
			query: {
				prefix: query
			}
		}
	});

	if (!res.error && res.data) {
		projectCache.set(
			query,
			res.data.map((project) => project.name)
		);
	}
}, 300);

export const ProjectTag = Mention.configure({
	HTMLAttributes: {
		class: 'mention'
	},
	suggestion: {
		items: async ({ query }) => {
			// Trigger the debounced API call
			debouncedProjectQuery(query);

			// Return cached results immediately
			const hasCache = projectCache.has(query);

			if (hasCache) {
				return projectCache.get(query)!;
			}

			const options = await apiClient
				.GET('/api/projects', {
					params: {
						query: {
							prefix: query
						}
					}
				})
				.then((res) => {
					if (res.error) {					
						return [];
					}
					return res.data.map((project) => project.name);
				});

			// Include the current query as an option if it's non-empty
			// This allows users to create new projects on the fly
			if (query.length > 0) {
				return Array.from(new Set([...options, query]));
			}

			return options;
		},
		render() {
			let toolTip: ToolTip;
			return {
				onStart: async (props) => {
					toolTip = new ToolTip(props.command);

					if (props.decorationNode) {
						toolTip.mount();
						await toolTip.setOptions(props.items);
						await toolTip.updatePosition(props.decorationNode);
					}
				},
				onUpdate: async (props) => {
					if (props.decorationNode) {
						await toolTip.setOptions(props.items);
						await toolTip.updatePosition(props.decorationNode);
						toolTip.setSelectedIndex(0);
					}
				},
				onKeyDown(props) {
					// Boolean return to indicate whether the key event was handled
					if (props.event.key === 'Escape') {
						toolTip.remove();
						return true;
					}

					if (props.event.key === 'ArrowDown') {
						toolTip.setSelectedIndex(toolTip.selectedIndex + 1);
						return true;
					}

					if (props.event.key === 'ArrowUp') {
						toolTip.setSelectedIndex(toolTip.selectedIndex - 1);
						return true;
					}

					if (props.event.key === 'Enter') {
						toolTip.selectCurrent();
						return true;
					}

					return false;
				},
				onExit: () => {
					toolTip.remove();
				}
			};
		}
	},
	renderText({ options, node }) {
		return `${options.suggestion.char}${node.attrs.label ?? node.attrs.id}`;
	},
	renderHTML({ node, options }) {
		// console.log('Rendering mention HTML for node:', options.HTMLAttributes);
		return [
			'span',
			mergeAttributes(options.HTMLAttributes, {
				'data-mention-id': node.attrs.id
			}),
			`@${node.attrs.label ?? node.attrs.id}`
		];
	}
});

export type MentionJSON = {
	type: 'mention';
	attrs: {
		id: string;
		label: string;
	};
};

export type MentionNodes = {
	node: string;
	mentioned: string[];
};

const MENTION_PARENT_TYPES = ['paragraph', 'taskItem'];

/**
 * Traverse the document JSON to find all mention nodes and return their nearest parent nodes.
 * Only searches for parents of types 'paragraph' and 'taskItem'.
 *
 * @param doc - The document JSON structure from editor.getJSON()
 * @returns Array of parent nodes that contain mentions, with their children
 */
export function findMentionParents(doc: JSONContent): MentionNodes[] {
	const parentNodes: Map<JSONContent, JSONContent[]> = new Map();

	function traverse(node: JSONContent): void {
		if (!node || !node.content || node.content.length === 0) {
			return;
		}

		// Check if this node is a mentionable parent type
		if (node.type && MENTION_PARENT_TYPES.includes(node.type)) {
			// Check if any direct child is a mention node
			const hasMention = node.content.some((child) => child.type === 'mention');

			if (hasMention) {
				parentNodes.set(node, node.content);
			}
		}

		// Recursively traverse all children
		for (const child of node.content) {
			traverse(child);
		}
	}

	// Start traversal from the root node
	traverse(doc);

	// Convert map to array of objects
	return Array.from(parentNodes.entries()).map(([parent, children]) => ({
		node: JSON.stringify(parent),
		mentioned: children
			.filter((child) => child.type === 'mention' && child.attrs && child.attrs.id)
			.map((mentionNode) => mentionNode.attrs?.id)
	}));
}
