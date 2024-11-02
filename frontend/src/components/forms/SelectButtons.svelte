<script lang="ts">
  import { addRipple, barFocusIndicator } from "$lib/client/decoration";
  import Label from "./Label.svelte";

  export let value: string;
  export let name: string;
  export let placeholder: string;
  export let label: boolean = true;

  export let editable: boolean = true;

  export let options: Option[];
  // TODO: redo this
  let selected: Option = options[0];
  $: selected = options.filter(option => option.value === value)[0];
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
    --barFocusIndicatorColor: #{$backgroundSecondary};
  }
</style>

{#if label}
<Label name={name}>{placeholder}</Label>
{/if}
{#if editable}
  <div class="buttons">
    {#each options as option, i}
      <button
        type="button"
        class:selected={option.value === value}
        class:first={i === 0}
        class:last={i === options.length - 1}
        on:click={() => {value = option.value}}
        on:mousedown={addRipple}
        use:barFocusIndicator
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