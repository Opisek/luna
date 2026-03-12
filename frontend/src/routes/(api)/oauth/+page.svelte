<script lang="ts">
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';
	import type { PageProps } from './$types';

  let { data }: PageProps = $props();

  onMount(() => {
    if (!browser) return;

    if (data.request) {
      localStorage.setItem(`oauth/${data.request}/response`, JSON.stringify(data))
    }

    window.close();
  });
</script>