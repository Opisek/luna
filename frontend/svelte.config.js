import path from 'path';
import adapter from 'svelte-adapter-bun';
import { sveltePreprocess } from 'svelte-preprocess';
import { fileURLToPath } from 'url';

const dirname = path.dirname(fileURLToPath(import.meta.url));

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: sveltePreprocess({
		// https://stackoverflow.com/questions/76021948/vscode-cant-find-scss-import-in-svelte-component-if-path-contains-lib-alias
		scss: {
			importer: [
				(url) => {
					if (url.startsWith('$lib')) {
						return {
							file: url.replace(/^\$lib/, path.join(dirname, 'src', 'lib')),
						};
					}
					return url;
				},
			]
		}
	}),
	kit: {
		adapter: adapter({
			precompress: true
		}),
		csrf: {
			checkOrigin: false // own solution implemented in src/routes/api/[...endpoint]/+server.ts
		}
	},
	compilerOptions: {
		experimental: {
			async: true
		}
	}
};

export default config;
