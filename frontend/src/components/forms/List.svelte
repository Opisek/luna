<script lang="ts" generics="T">
  import type { Snippet } from 'svelte';
  import Label from './Label.svelte';
  import Tooltip from '../interactive/Tooltip.svelte';
  import { t } from '@sveltia/i18n';
  import { extendUniqueElementId, generateUniqueElementId } from '$lib/common/dom';

  interface Props {
    label: string;
    info?: string;
    items: T[];
    template: Snippet<[T]>;
    id: (item: T) => string;
  }

  let {
    label,
    info = "",
    items,
    template,
    id,
  }: Props = $props();

  let uniqueId = $props.id();
  let listId = $derived(generateUniqueElementId(["list"], uniqueId));
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
  }
</style>

<Label describes={listId} info={info}>
  {label}
</Label>
<div
  aria-live="polite"
  aria-relevant="all"
  id={listId}
  aria-labelledby={extendUniqueElementId(["label"], listId)}
  role="list"
>
  {#each items as item (id(item))}
    {@render template(item)}
  {:else}
    {t("error.empty")}
  {/each}
</div>