<script lang="ts" generics="T">
  import type { Snippet } from 'svelte';
  import Label from './Label.svelte';
  import Tooltip from '../interactive/Tooltip.svelte';

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
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapSmall;
  }
</style>

<Label name={label.toLowerCase().replaceAll(" ", "-")}>
  {label}
  {#if info != ""}
    <Tooltip>
      {info}
    </Tooltip>
  {/if}
</Label>
<div>
  {#each items as item (id(item))}
    {@render template(item)}
  {/each}
</div>