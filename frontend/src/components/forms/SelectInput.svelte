<script lang="ts">
  import { ChevronDown } from "lucide-svelte";

  import { browser } from "$app/environment";

  import Label from "./Label.svelte";

  import { calculateOptimalPopupPosition } from "$lib/common/calculations";
  import { focusIndicator } from "$lib/client/decoration";

  let active = $state(false);
  let optionsAbove = $state(false);

  interface Props {
    value: string;
    placeholder: string;
    name: string;
    editable?: boolean;
    options: Option[];
  }

  let {
    value = $bindable(""),
    placeholder,
    name,
    editable = true,
    options
  }: Props = $props();

  let selectedOption: Option | null = $derived(options.filter(x => x.value === value)[0] || null);

  let selectWrapper: HTMLElement;
  let optionsWrapper: HTMLElement;

  function selectClick() {
    if (!active) {
      const res = calculateOptimalPopupPosition(selectWrapper)
      optionsAbove = res.bottom;
      if (browser) {
        window.addEventListener("click", clickOutside);
      }
      setTimeout(() => {
        (optionsWrapper.getElementsByClassName("selected")[0] as HTMLElement).focus();
      }, 0);
    } else {
      window.removeEventListener("click", clickOutside);
    }
    active = !active;
  }

  function optionClick(option: Option) {
    value = option.value;
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
    cursor: pointer;
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

  .placeholder {
    color: $foregroundFaded;
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
    onclick={selectClick}
    type="button"
    use:focusIndicator={{ type: "bar", ignoreParent: true }}
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
  <div
    class="options"
    class:hidden={!active}
    style="top: {optionsAbove ? -100 * options.length : 100}%"
    bind:this={optionsWrapper}
  >
    {#each options as option}
      <button
        class="option" 
        onclick={() => optionClick(option)}
        type="button"
        class:selected={option.value === value}
      >
        {option.name}
      </button>
    {/each}
  </div>
</div>