<script lang="ts">
  import { ChevronDown } from "lucide-svelte";
  import Label from "./Label.svelte";
  import { calculateOptimalPopupPosition } from "$lib/common/calculations";
  import { browser } from "$app/environment";
  import { addRipple } from "../../lib/client/decoration";

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

  let clickedWithMouse = false;
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
    overflow: hidden;
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

  div.focus {
    background-color: $backgroundAccent;
    height: 100%;
    width: $borderActiveWidth;
    position: absolute;
    left: 0;
    top: 0;
    transform: translateX(-100%);
    transition: transform $animationSpeedFast linear;
  }

  button.select:focus-within:not(.click) > div.focus {
    transform: none;
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
    class:click={clickedWithMouse}
    on:click={selectClick}
    on:mousedown={(e) => {addRipple(e); clickedWithMouse = true}}
    on:focusout={() => {if (!active) clickedWithMouse = false}}
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
    <div class="focus"></div>
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