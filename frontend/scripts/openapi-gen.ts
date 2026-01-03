import { $ } from 'bun';
const specEndpoint = 'http://localhost:8000/docs/openapi.json';
const specOutput = 'openapi.json';
const clientOutput = 'src/lib/api/spec.ts';

const res = await fetch(specEndpoint);

if (!res.ok) {
	console.error(
		`Failed to fetch OpenAPI spec from ${specEndpoint}: ${res.status} ${res.statusText}`
	);
}

const spec = await res.json();

const file = Bun.write(specOutput, JSON.stringify(spec, null, 2));
console.log(`Wrote OpenAPI spec to ${specOutput} (${file.bytesWritten} bytes)`);

await $`pnpm openapi-typescript ${specOutput} --output ${clientOutput}`;
