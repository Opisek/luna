import adapter from 'svelte-adapter-bun';
import { sveltePreprocess } from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: sveltePreprocess({}),
	kit: {
		adapter: adapter({
			precompress: true
		}),
		csrf: {
			checkOrigin: false // own solution implemented in src/routes/api/[...endpoint]/+server.ts
		}
	}
};

export default config;
