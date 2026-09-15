<script lang="ts" generics="T">
  import { ChevronDown } from "lucide-svelte";

  import Label from "./Label.svelte";
  import Popup from "../popups/Popup.svelte";

  import { focusIndicator } from "$lib/client/decoration";
  import type { Option } from "../../types/options";
  import { AsyncNoOp, NoOp } from "$lib/client/placeholders";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";

  let active = $state(false);

  interface Props {
    value: T | null | undefined;
    placeholder: string;
    name: string;
    editable?: boolean;
    options: Option<T>[];
    showLabel?: boolean;
    click?: (value: T) => void;
  }

  let {
    value = $bindable(),
    placeholder,
    name,
    editable = true,
    options,
    showLabel = true,
    click = NoOp,
  }: Props = $props();
  let uniqueId = $props.id();
  let selectionId = $derived(generateUniqueElementId(["select"], uniqueId));

  let selectedOption: Option<T> | null = $derived(options.filter(x => x.value === value)[0] || null);

  let selectWrapper: HTMLElement | undefined = $state();

  let showPopup = $state(AsyncNoOp);
  let hidePopup = $state(NoOp);

  function selectClick() {
    if (!editable) return;

    if (!active) showPopup();
    else hidePopup();
  }

  function optionClick(option: Option<T>) {
    value = option.value;
    hidePopup();
    click(option.value);
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/dimensions.scss";

  button.select {
    all: unset;
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadius;
    background: transparent;
    display: flex;
    align-items: center;
    gap: dimensions.$gapSmall;
    justify-content: space-between;
    position: relative;
    transition: padding animations.$animationSpeedFast linear, border-radius animations.$animationSpeedFast linear, width animations.$animationSpeedFast linear;
    overflow: hidden;
  }

  button.editable {
    color: colors.$foregroundSecondary;
    background: colors.$backgroundSecondary;
    cursor: pointer;
    user-select: none;
  }

  select {
    display: none;
  }

  span.arrow {
    height: 100%;
    display: flex;
    align-items: center;
    transition: animations.$cubic animations.$animationSpeed;
  }

  span.arrow.active {
    transform: rotate(-180deg);
  }
  
  button.option {
    all: unset;
    transition: linear animations.$animationSpeedFast;
    width: 100%;
    padding: dimensions.$gapSmall;
    cursor: pointer;
  }

  button.option.selected {
    color: colors.$foregroundAccent;
    background-color: colors.$backgroundAccent;
  }

  button.option:hover, button.option:focus {
    color: colors.$foregroundTertiary;
    background-color: colors.$backgroundTertiary;
  }

  div.wrapper {
    width: 100%;
    padding-right: 2 * dimensions.$gapSmall;
    position: relative;
  }
  button {
    width: 100% !important;
  }

  .placeholder {
    color: color-mix(in srgb, colors.$foregroundSecondary 50%, transparent);
  }
</style>

{#if showLabel}
  <Label describes={selectionId}>{placeholder}</Label>
{/if}
<div class="wrapper" class:editable={editable}>
  <select
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
  ></select>
  <button
    id={selectionId}
    bind:this={selectWrapper}
    class="select"
    class:editable={editable}
    onclick={selectClick}
    type="button"
    use:focusIndicator={{ type: "bar" }}
    aria-labelledby={extendUniqueElementId(["label", selectionId])}
  >
    {#if selectedOption !== null}
      {selectedOption.name}
    {:else}
      <span class="placeholder">
        {"Select " + placeholder}
      </span>
    {/if}
    {#if editable}
      <span
        class="arrow"
        class:active={active} 
      >
        <ChevronDown size={16}/>
      </span>
    {/if}
  </button>
  <Popup anchor={selectWrapper} matchWidth={true} tooltip={false} triangle={false} bind:showPopup bind:hidePopup --padding="0" bind:visible={active}>
    {#each options as option (option.value)}
      <button
        class="option" 
        onclick={() => optionClick(option)}
        type="button"
        class:selected={option.value == value}
      >
        {option.name}
      </button>
    {/each}
  </Popup>
</div>