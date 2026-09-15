<script lang="ts" generics="T">
  import { ChevronDown } from "lucide-svelte";

  import Label from "./Label.svelte";
  import Popup from "../popups/Popup.svelte";

  import { focusIndicator } from "$lib/client/decoration";
  import type { Option } from "../../types/options";
  import { AsyncNoOp, NoOp } from "$lib/client/placeholders";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";
  import { passIfEnter } from "$lib/common/inputs";

  let active = $state(false);

  interface Props {
    values: T[];
    placeholder: string;
    name: string;
    editable?: boolean;
    options: Option<T>[];
    showLabel?: boolean;
    click?: (value: T) => void;
  }

  let {
    values = $bindable(),
    placeholder,
    name,
    editable = true,
    options,
    showLabel = true,
    click = NoOp,
  }: Props = $props();
  let uniqueId = $props.id();
  let selectionId = $derived(generateUniqueElementId(["select"], uniqueId));
  let labelId = $derived(extendUniqueElementId(["label"], selectionId));
  let listId = $derived(extendUniqueElementId(["listbox"], selectionId));
  let optionEntryIdFn = (option: Option<T>) => extendUniqueElementId(["option"].concat(`${option.value}`.toLocaleLowerCase().split(" ")), selectionId);

  let valuesSet = $derived(new Set(values));
  let selectedOptions: Option<T>[] = $derived(options.filter(x => valuesSet.has(x.value)));
  let focusedOption: Option<T> | null = $state(null);

  let selectWrapper: HTMLElement | undefined = $state();

  let showPopup = $state(AsyncNoOp);
  let hidePopup = $state(NoOp);

  function selectClick() {
    if (!editable) return;

    if (!active) {
      showPopup();
      let focusOption = focusedOption ?? selectedOptions.find(() => true) ?? (options.length != 0 ? options[1] : null)
      if (focusOption !== null) {
        setTimeout(() => {
          document.getElementById(optionEntryIdFn(focusOption))?.focus();
        }, 0);
      }
    } else hidePopup();
  }

  function optionClick(option: Option<T>) {
    if (valuesSet.has(option.value)) values = values.filter(x => x != option.value);
    else values.push(option.value);
    click(option.value);
  }

  function optionPress(event: KeyboardEvent, option: Option<T>) {
    passIfEnter(event, () => optionClick(option))
  }

  function optionFocus(option: Option<T>) {
    focusedOption = option;
  }

  $effect(() => {
    if (active) addEventListener("keydown", arrowKeyListener);
    else removeEventListener("keydown", arrowKeyListener);
  })

  function arrowKeyListener(event: KeyboardEvent) {
    if (focusedOption === null) return;

    const currentlyFocusedElement = document.getElementById(optionEntryIdFn(focusedOption));
    if (currentlyFocusedElement === null) return;

    const index = options.findIndex(x => x.value === focusedOption?.value);

    let newlyFocusedOption = null;
    if (event.key == "ArrowUp" && index > 0) {
      newlyFocusedOption = options[index - 1];
    }
    else if (event.key == "ArrowDown" && index < options.length-1) {
      newlyFocusedOption = options[index + 1];
    }
    if (newlyFocusedOption === null) return;

    focusedOption = newlyFocusedOption;
    let newlyFocusedElement = document.getElementById(optionEntryIdFn(newlyFocusedOption));
    if (newlyFocusedElement === null) return;

    newlyFocusedElement.focus();
    event.preventDefault();
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
  
  .option {
    all: unset;
    transition: linear animations.$animationSpeedFast;
    width: 100%;
    padding: dimensions.$gapSmall;
    cursor: pointer;
  }

  .option.selected {
    color: colors.$foregroundAccent;
    background-color: colors.$backgroundAccent;
  }

  .option:hover, .option:focus {
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

  .list {
    display: contents;
  }
</style>

{#if showLabel}
  <Label describes={selectionId}>{placeholder}</Label>
{/if}
<div class="wrapper" class:editable={editable}>
  <select
    bind:value={values}
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
    use:focusIndicator={{ type: "bar" }}
    aria-labelledby={labelId}
    aria-expanded={active}
    aria-controls={listId}
    aria-activedescendant={focusedOption !== null ? optionEntryIdFn(focusedOption) : undefined}
    role="combobox"
  >
    {#if selectedOptions.length != 0}
      {selectedOptions.map(x => x.name).join(", ")}
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
  <Popup
    anchor={selectWrapper}
    matchWidth={true}
    tooltip={false}
    triangle={false}
    bind:showPopup
    bind:hidePopup
    bind:visible={active}
    --padding="0"
    labelled={labelId}
  >
    <div
      id={listId}
      class="list"
      role="listbox"
      aria-multiselectable="true"
      aria-labelledby={labelId}
    >
      {#each options as option (option.value)}
        {@const optionEntryId = optionEntryIdFn(option)}
        <div
          id={optionEntryId}
          class="option" 
          class:selected={valuesSet.has(option.value)}
          aria-selected={valuesSet.has(option.value)}
          onclick={() => optionClick(option)}
          onkeypress={(e) => optionPress(e, option)}
          onfocusin={() => optionFocus(option)}
          role="option"
          tabindex=0
        >
          {option.name}
        </div>
      {/each}
    </div>
  </Popup>
</div>