<script lang="ts" generics="T">
  import Label from "./Label.svelte";

  import { addRipple, focusIndicator } from "$lib/client/decoration";
  import { EmptyOption } from "../../lib/client/placeholders";

  interface Props {
    value: T | null;
    name: string;
    placeholder?: string;
    label?: boolean;
    editable?: boolean;
    compact?: boolean;
    options: Option<T>[];
  }

  let {
    value = $bindable(),
    name,
    placeholder = "",
    label = true,
    editable = true,
    compact = false,
    options
  }: Props = $props();

  let selected: Option<T> = $derived(options.filter(option => option.value === value)[0] || options[0] || EmptyOption);
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div.display {
    margin: dimensions.$gapSmall;
  }

  div.buttons {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    width: 100%; 
    gap: dimensions.$gapSmaller;
    user-select: none;
  }

  div.compact {
    width: max-content;
    gap: 0;
  }
  div.compact > button {
    min-width: dimensions.$buttonMinWidthCompact;
    padding: dimensions.$gapSmall;
  }

  button {
    all: unset;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    padding: dimensions.$gapSmall;
    cursor: pointer;
    flex: 1;
    text-align: center;
    position: relative;
    overflow: hidden;
  }

  button.first {
    border-top-left-radius: dimensions.$borderRadius;
    border-bottom-left-radius: dimensions.$borderRadius;
  }

  button.last {
    border-top-right-radius: dimensions.$borderRadius;
    border-bottom-right-radius: dimensions.$borderRadius;
  }

  button.selected {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt};
  }
</style>

{#if label && placeholder}
  <Label name={name}>{placeholder}</Label>
{/if}
{#if editable}
  <div
    class="buttons"
    class:compact={compact} 
  >
    {#each options as option, i}
      <button
        type="button"
        class:selected={option.value === value}
        class:first={i === 0}
        class:last={i === options.length - 1}
        onclick={() => {value = option.value}}
        onmousedown={addRipple}
        use:focusIndicator
      >
        {option.name}
      </button>
    {/each}
  </div>
{:else}
  <div class="display">
    {selected.name}
  </div>
{/if}