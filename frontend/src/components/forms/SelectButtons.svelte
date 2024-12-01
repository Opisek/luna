<script lang="ts">
  import Label from "./Label.svelte";

  import { addRipple, focusIndicator } from "$lib/client/decoration";
  import { EmptyOption } from "../../lib/client/placeholders";

  interface Props {
    value: string;
    name: string;
    placeholder?: string;
    label?: boolean;
    editable?: boolean;
    compact?: boolean;
    options: Option[];
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

  // TODO: redo this
  let selected: Option = $derived(options.filter(option => option.value === value)[0] || options[0] || EmptyOption);
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.display {
    margin: $gapSmall;
  }

  div.buttons {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    width: 100%; 
    gap: $gapSmaller;
    user-select: none;
  }

  div.compact {
    width: max-content;
    gap: 0;
  }
  div.compact > button {
    min-width: max-content;
    padding: $gapSmaller;
  }

  button {
    all: unset;
    background-color: $backgroundSecondary;
    color: $foregroundSecondary;
    padding: $gapSmall;
    cursor: pointer;
    flex: 1;
    text-align: center;
    position: relative;
    overflow: hidden;
  }

  button.first {
    border-top-left-radius: $borderRadius;
    border-bottom-left-radius: $borderRadius;
  }

  button.last {
    border-top-right-radius: $borderRadius;
    border-bottom-right-radius: $borderRadius;
  }

  button.selected {
    background-color: $backgroundAccent;
    color: $foregroundAccent;
    --barFocusIndicatorColor: #{$barFocusIndicatorColorAlt};
  }
</style>

{#if label}
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