<script lang="ts" generics="T">
  import { ChevronDown } from "lucide-svelte";

  import { browser } from "$app/environment";

  import Label from "./Label.svelte";

  import { calculateOptimalPopupPosition } from "$lib/common/calculations";
  import { focusIndicator } from "$lib/client/decoration";
  import type { Option } from "../../types/options";

  let active = $state(false);
  let optionsAbove = $state(false);

  interface Props {
    value: T | null;
    placeholder: string;
    name: string;
    editable?: boolean;
    options: Option<T>[];
  }

  let {
    value = $bindable(null),
    placeholder,
    name,
    editable = true,
    options
  }: Props = $props();

  let selectedOption: Option<T> | null = $derived(options.filter(x => x.value === value)[0] || null);

  let selectWrapper: HTMLElement;
  let optionsWrapper: HTMLElement;

  function selectClick() {
    if (!active && !editable) return;
    active = !active;
    if (!active) return;

    const res = calculateOptimalPopupPosition(selectWrapper)
    optionsAbove = res.bottom;
    if (browser) {
      window.addEventListener("click", clickOutside);
      window.addEventListener("keydown", keyboardClick);
    }
    setTimeout(() => {
      const els = optionsWrapper.getElementsByClassName("selected");
      if (els.length > 0 && els[0]) (els[0] as HTMLElement).focus();
    }, 10);
  }

  $effect(() => {
    if (!active) {
      window.removeEventListener("click", clickOutside);
      window.removeEventListener("keydown", keyboardClick);
    }
  })

  function optionClick(option: Option<T>) {
    value = option.value;
    selectWrapper.focus();
  }

  function clickOutside(event: MouseEvent) {
    if (!selectWrapper || selectWrapper == event.target as Node || selectWrapper.contains(event.target as Node)) return;

    active = false;
    event.stopPropagation();
  }

  function keyboardClick(event: KeyboardEvent) {
    if (!optionsWrapper) return;

    switch(event.key) {
      case "ArrowDown":
        event.preventDefault();
        focusNext();
        break;
      case "ArrowUp":
        event.preventDefault();
        focusPrevious();
        break;
      case "Escape":
        event.preventDefault();
        active = false;
        break;
    }
  }

  function focusNext() {
    const currentFocus = document.activeElement;

    let currentIndex = -1;
    for (const [i, option] of Array.from(optionsWrapper.children).entries()) {
      if (option === currentFocus) {
        currentIndex = i;
        break;
      }
    }
    currentIndex++;

    if (currentIndex >= options.length) return;

    (optionsWrapper.children[currentIndex] as HTMLElement).focus();
  }

  function focusPrevious() {
    const currentFocus = document.activeElement;

    let currentIndex = options.length;
    for (const [i, option] of Array.from(optionsWrapper.children).entries()) {
      if (option === currentFocus) {
        currentIndex = i;
        break;
      }
    }
    currentIndex--;

    if (currentIndex < 0) return;

    (optionsWrapper.children[currentIndex] as HTMLElement).focus();
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

  div.options {
    position: absolute;
    background-color: colors.$backgroundSecondary;
    width: 100%;
    left: 0;
    box-shadow: decorations.$boxShadow;
    z-index: 10;
    border-radius: dimensions.$borderRadius;
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
    transition: linear animations.$animationSpeedFast;
    padding: dimensions.$gapSmall;
    cursor: pointer;
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