import adapter from '@sveltejs/adapter-node';
import { sveltePreprocess } from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: sveltePreprocess({}),
	kit: {
		adapter: adapter({
			precompress: true
		}),
		csrf: {
			checkOrigin: false // TODO: reenable in production or implement an alterative
		}
	}
};

export default config;
