<script lang="ts">
  import { ChevronDown } from "lucide-svelte";
  import Label from "./Label.svelte";
  import { calculateOptimalPopupPosition } from "$lib/common/calculations";

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
    }
    active = !active;
  }

  function optionClick(option: Option) {
    selectedOption = option;
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
  
  div.optons.above {
    top: -100%;
  }

  div.options.below {
    top: 100%;
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

  button.option:hover {
    background-color: $backgroundSecondary;
  }

  // TODO: style the drop-down
</style>

<Label name={name}>{placeholder}</Label>
<button
  bind:this={selectWrapper}
  class="select"
  class:editable={editable}
  on:click={selectClick}
>
  <select
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
  >
  <!--
    {#each options as option}
      <option value={option.value}>{option.name}</option>
    {/each}
  -->
  </select>
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
  <div
    class="options"
    class:hidden={!active}
    class:above={optionsAbove}
    class:below={!optionsAbove}
  >
    {#each options as option}
      <button
        class="option" 
        on:click={() => optionClick(option)}
      >
        {option.name}
      </button>
    {/each}
  </div>
</button>