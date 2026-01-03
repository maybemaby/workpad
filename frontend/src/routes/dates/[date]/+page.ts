import { error, redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import { apiClient } from '$lib/api/client';
import type { components } from '$lib/api/spec';
import { getLocalTimeZone, today } from '@internationalized/date';

const dateRegex = /^\d{4}-\d{1,2}-\d{1,2}$/;

export const load: PageLoad = async (event) => {
	const date = event.params.date;

	const res = date.match(dateRegex);

	const segments = date.split('-').map(Number);

	if (!res || segments.length !== 3) {
		error(400, 'Invalid date format');
	}
	const currentDate = today(getLocalTimeZone());
	const formattedDate = `${segments[0]}-${String(segments[1]).padStart(2, '0')}-${String(segments[2]).padStart(2, '0')}`;

	if (currentDate.toString() === formattedDate) {
		redirect(307, `/`);
	}

	const noteRes = await apiClient.GET('/notes/by-date', {
		params: {
			query: {
				date: formattedDate
			}
		}
	});

	let note: null | components['schemas']['NotesNote'] = null;

	if (noteRes.data) {
		note = noteRes.data;
	}

	return {
		year: segments[0],
		month: segments[1],
		day: segments[2],
		currentDate,
		note
	};
};
