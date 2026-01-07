import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import QuickSwitcher from './quick-switcher.svelte';
import { getLocalTimeZone, now } from '@internationalized/date';

describe('QuickSwitcher', () => {
	it('should render hidden by default', () => {
		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		expect(() => screen.getByRole('dialog')).toThrow();
	});

	it('should render on hotkey', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');

		expect(screen.getByRole('dialog')).toBeInTheDocument();
	});

	it('should close on hotkey', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByRole('dialog')).toBeInTheDocument();

		await user.keyboard('{Meta>}j');
		expect(() => screen.getByRole('dialog')).toThrow();
	});

	it('should close on Escape', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByRole('dialog')).toBeInTheDocument();

		await user.keyboard('{Escape}');
		expect(() => screen.getByRole('dialog')).toThrow();
	});

	it.each([
		[{ id: '/dates/[date]' as const, path: '/dates/2024-06-15' }, '2024-06-15'],
		[{ id: '/' as const, path: '/' }, now(getLocalTimeZone()).toString().slice(0, 10)]
	])('should show correct default date for path %s', async (path, expectedDate) => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: path.id,
				path: path.path
			}
		});

		await user.keyboard('{Meta>}j');

		expect(screen.getByText(expectedDate)).toBeInTheDocument();
	});

	it('should change date on left/right arrow', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByText('2024-06-15')).toBeInTheDocument();

		await user.keyboard('{ArrowRight}');
		expect(screen.getByText('2024-06-16')).toBeInTheDocument();

		await user.keyboard('{ArrowLeft}{ArrowLeft}');
		expect(screen.getByText('2024-06-14')).toBeInTheDocument();
	});

	it('should change month on up/down arrow', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByText('2024-06-15')).toBeInTheDocument();

		await user.keyboard('{ArrowDown}');
		expect(screen.getByText('2024-07-15')).toBeInTheDocument();

		await user.keyboard('{ArrowUp}{ArrowUp}');
		expect(screen.getByText('2024-05-15')).toBeInTheDocument();
	});

	it('should reset date on close', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15'
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByText('2024-06-15')).toBeInTheDocument();

		await user.keyboard('{ArrowRight}{ArrowRight}');
		expect(screen.getByText('2024-06-17')).toBeInTheDocument();

		await user.keyboard('{Escape}');

		await user.keyboard('{Meta>}j');
		expect(screen.getByText('2024-06-15')).toBeInTheDocument();
	});

	it('should not open on hotkey if route is not switchable', async () => {
		const user = userEvent.setup();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				// @ts-expect-error testing behavior
				id: '/other/[id]',
				path: '/other/123'
			}
		});

		await user.keyboard('{Meta>}j');

		expect(() => screen.getByRole('dialog')).toThrow();
	});

	it('should call callback on Enter', async () => {
		const user = userEvent.setup();
		const onSelect = vi.fn();

		render(QuickSwitcher, {
			target: document.body,
			props: {
				id: '/dates/[date]',
				path: '/dates/2024-06-15',
				onSelected(dateString) {
					onSelect(dateString);
				}
			}
		});

		await user.keyboard('{Meta>}j');
		expect(screen.getByText('2024-06-15')).toBeInTheDocument();

		await user.keyboard('{Enter}');
		expect(onSelect).toHaveBeenCalledWith('2024-06-15');
	});
});
