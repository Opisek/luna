<script lang="ts">
  import { ChevronDown } from "lucide-svelte";
  import Label from "./Label.svelte";
  import { calculateOptimalPopupPosition } from "$lib/common/calculations";
  import { browser } from "$app/environment";

  export let value: string;
  export let placeholder: string;
  export let name: string;

  export let editable: boolean = true;
  let active = false;
  let optionsAbove = false;

  export let options: Option[];
  let selectedOption: Option | null = options.length > 0 ? options[1] : null;
  $: if (selectedOption) value = selectedOption.value;


  let selectWrapper: HTMLElement;
  function selectClick() {
    if (!active) {
      const res = calculateOptimalPopupPosition(selectWrapper)
      optionsAbove = res.bottom;
      if (browser) {
        window.addEventListener("click", clickOutside);
      }
    } else {
      window.removeEventListener("click", clickOutside);
    }
    active = !active;
  }

  function optionClick(option: Option) {
    selectedOption = option;
    selectWrapper.focus();
  }

  function clickOutside(event: MouseEvent) {
    if (!selectWrapper || selectWrapper == event.target as Node || selectWrapper.contains(event.target as Node)) return;

    active = false;
    event.stopPropagation();
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/decoration.scss";
  @import "../../styles/dimensions.scss";

  button.select {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
    background: transparent;
    display: flex;
    align-items: center;
    gap: $gapSmall;
    justify-content: space-between;
    position: relative;
    transition: padding $animationSpeedFast linear, border-radius $animationSpeedFast linear, width $animationSpeedFast linear;
  }

  button.editable {
    background: $backgroundSecondary;
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
    transition: $cubic $animationSpeed;
  }

  span.arrow.active {
    transform: rotate(-180deg);
  }

  div.options {
    position: absolute;
    background-color: $backgroundPrimary;
    width: 100%;
    left: 0;
    box-shadow: $boxShadow;
    z-index: 10;
    border-radius: $borderRadius;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  div.options.hidden {
    display: none;
  }

  button.option {
    all: unset;
    width: 100%;
    transition: linear $animationSpeedFast;
    padding: $gapSmall;
  }

  button.option:hover, button.option:focus {
    background-color: $backgroundSecondary;
  }

  div.wrapper {
    width: 100%;
    padding-right: 2 * $gapSmall;
    position: relative;
  }
  button {
    width: 100% !important;
  }
  div.wrapper.editable {
    background-color: $backgroundAccent;
    border-radius: calc($borderRadius + 0.1em);
    transition: padding $animationSpeedFast linear;
  } 
  div.wrapper.editable:focus-within {
    padding-left: $borderActiveWidth;
    border-top-left-radius: $borderRadius;
    border-bottom-left-radius: $borderRadius;
  }
  div.wrapper.editable:focus-within > button {
    padding-left: calc($gapSmall - $borderActiveWidth);
    width: calc(100% + $borderActiveWidth) !important;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div class="wrapper" class:editable={editable}>
  <select
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
  ></select>
  <button
    bind:this={selectWrapper}
    class="select"
    class:editable={editable}
    on:click={selectClick}
    type="button"
  >
    {#if selectedOption !== null}
      {selectedOption.name}
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
  <div
    class="options"
    class:hidden={!active}
    style="top: {optionsAbove ? -100 * options.length : 100}%"
  >
    {#each options as option}
      <button
        class="option" 
        on:click={() => optionClick(option)}
        type="button"
      >
        {option.name}
      </button>
    {/each}
  </div>
</div>