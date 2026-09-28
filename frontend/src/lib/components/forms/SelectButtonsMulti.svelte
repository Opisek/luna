<script lang="ts" generics="T">
  import Label from "./Label.svelte";

  import { addRipple, focusIndicator } from "$lib/scripts/client/decoration";
  import { NoOp } from "$lib/scripts/client/placeholders";
  import type { Option } from "$lib/types/options";

  interface Props {
    values: T[];
    name: string;
    placeholder?: string;
    info?: string;
    label?: boolean;
    editable?: boolean;
    compact?: boolean;
    options: Option<T>[];
    onClick?: (selected: T) => any;
  }

  let {
    values = $bindable(),
    name,
    placeholder = "",
    info = "",
    label = true,
    editable = true,
    compact = false,
    options,
    onClick = NoOp,
  }: Props = $props();

  let valuesSet = $derived(new Set(values));
  let selected: Option<T>[] = $derived(options.filter(option => valuesSet.has(option.value)));
</script>

<style lang="scss">
  @use "$lib/styles/animations.scss";
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";

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
  <Label name={name} info={info}>{placeholder}</Label>
{/if}
{#if editable}
  <div
    class="buttons"
    class:compact={compact} 
  >
    {#each options as option, i}
      <button
        type="button"
        class:selected={valuesSet.has(option.value)}
        class:first={i === 0}
        class:last={i === options.length - 1}
        onclick={() => {
          if (valuesSet.has(option.value)) values = values.filter(x => x != option.value);
          else values.push(option.value);
          onClick(option.value);
        }}
        onmousedown={addRipple}
        use:focusIndicator
      >
        {option.name}
      </button>
    {/each}
  </div>
{:else}
  <div class="display">
    {#if selected.length == 0}
      None selected
    {:else}
      {selected.map(x => x.name).join(", ")}
    {/if}
  </div>
{/if}