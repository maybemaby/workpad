import createClient from 'openapi-fetch';
import type { paths } from './spec';

export const apiClient = createClient<paths>({
	baseUrl: import.meta.env.DEV ? 'http://localhost:8000' : window.location.origin,
	credentials: 'include',
	mode: 'cors'
});
