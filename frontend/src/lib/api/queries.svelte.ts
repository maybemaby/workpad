import { createMutation, createQuery } from '@tanstack/svelte-query';
import { apiClient } from './client';
import type { components } from './spec';
import type { Getter } from 'runed';

export const createGetProjectsQuery = (query?: string) => {
	return createQuery(() => ({
		queryKey: ['projects', query],
		queryFn: async () => {
			const res = await apiClient.GET('/api/projects', {
				params: {
					query: {
						prefix: query
					}
				}
			});
			if (res.error) {
				throw new Error('Failed to fetch projects');
			}

			return res.data;
		}
	}));
};

export const createProjectsMutation = () => {
	return createMutation(() => ({
		mutationFn: async (data: components['schemas']['ProjectsCreateMultipleProjectsRequest']) => {
			return apiClient.POST('/api/projects/batch', {
				body: data
			});
		}
	}));
};

export const createAddNoteMutation = () => {
	return createMutation(() => ({
		mutationFn: async (data: components['schemas']['NotesCreateNoteRequest']) => {
			return apiClient.POST('/api/notes', {
				body: data
			});
		}
	}));
};

export const getNoteByDateQuery = (date: string) => {
	return createQuery(() => ({
		queryKey: ['note-by-date', date],
		queryFn: async () => {
			const res = await apiClient.GET('/api/notes/by-date', {
				params: {
					query: {
						date
					}
				}
			});

			if (res.error && res.response.status === 404) {
				return null;
			} else if (res.error) {
				throw new Error('Failed to fetch note by date');
			}

			return res.data;
		},
		refetchOnWindowFocus: false
	}));
};

interface GetMonthNotesParams {
	enabled: boolean;
	year: number;
	month: number;
}

export const getMonthNotes = (props: Getter<GetMonthNotesParams>) => {
	return createQuery(() => ({
		queryKey: ['month-notes', props().year, props().month],
		queryFn: async () => {
			const res = await apiClient.GET('/api/notes/for-month', {
				params: {
					query: {
						year: props().year,
						month: props().month
					}
				}
			});

			if (res.error) {
				throw new Error('Failed to fetch month notes');
			}

			return res.data;
		},
		enabled: props().enabled
	}));
};

export const createUpdateExcerptMutation = () => {
	return createMutation(() => ({
		mutationFn: async (data: components['schemas']['NotesUpdateNoteExcerptRequest']) => {
			const res = await apiClient.PUT('/api/notes/excerpts', {
				body: data
			});

			if (res.error) {
				throw new Error('Failed to update note excerpt');
			}

			return res.data;
		}
	}));
};

export const getExcerptsQuery = (name: string) => {
	return createQuery(() => ({
		queryKey: ['note-excerpts', name],
		queryFn: async () => {
			const res = await apiClient.GET('/api/notes/excerpts/{project}', {
				params: {
					path: {
						project: name
					}
				}
			});

			if (res.error) {
				throw new Error('Failed to fetch note excerpts');
			}

			return res.data;
		}
	}));
};
