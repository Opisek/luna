import adapter from '@sveltejs/adapter-node';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: preprocess({}),
	kit: {
		adapter: adapter({
			precompress: true
		}),
		csrf: {
			checkOrigin: true
		}
	}
};

export default config;
