import createClient from 'openapi-fetch';
import type { paths } from './spec';

export const apiClient = createClient<paths>({
	baseUrl: 'http://localhost:8000',
	credentials: 'include',
	mode: 'cors'
});
