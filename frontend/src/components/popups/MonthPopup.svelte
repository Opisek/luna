<script lang="ts">
  import {  setContext } from "svelte";
  import { ChevronLeft, ChevronRight } from "lucide-svelte";

  import IconButton from "../interactive/IconButton.svelte";
  import Popup from "./Popup.svelte";

  import { NoOp } from '$lib/client/placeholders';
  import { focusIndicator } from "$lib/client/decoration";
  import { getMonthName } from "$lib/common/humanization";
  import { svelteFlyInHorizontal, svelteFlyOutHorizontal } from "$lib/client/animations";

  interface Props {
    date: Date;
    showPopup?: () => any;
    hidePopup?: () => any;
    onSelect?: (date: Date) => void;
  }

  let {
    date = $bindable(new Date()),
    showPopup = $bindable(NoOp),
    hidePopup = $bindable(NoOp),
    onSelect = NoOp,
  }: Props = $props();

  let popupVisible: boolean = $state(false);
  let selectingMonth: boolean = $state(true);

  let internalShow: () => void = $state(NoOp);
  let internalClose: () => void = $state(NoOp);

  /* Popup */
  showPopup = () => {
    if (popupVisible) return;
    selectedYear = date.getFullYear();
    selectingMonth = true;
    setTimeout(internalShow, 0);
  }

  hidePopup = () => {
    internalClose();
    onSelect(date);
  }

  /* Animation */
  let viewIteration = $state(0);
  let flyDirection = $state("left");
  setContext("flyDirection", () => flyDirection);

  /* Selection */
  let selectedMonth: number = $state(date.getMonth());
  let selectedYear: number = $state(date.getFullYear());
  let decadeStart: number = $derived(Math.floor(selectedYear / 10) * 10);

  function clickMonth(e: MouseEvent, i: number) {
    //addRipple(e);
    selectedMonth = i;
    date = new Date(selectedYear, selectedMonth, date.getDate());
    hidePopup();
  }

  function clickYear(e: MouseEvent, i: number) {
    e.stopPropagation();
    //addRipple(e);
    selectedYear = decadeStart + i;
    selectingMonth = true;
  }

  function prev() {
    flyDirection = "right";
    viewIteration++;
    if (selectingMonth) selectedYear --;
    else selectedYear -= 10;
  }

  function next() {
    flyDirection = "left";
    viewIteration++;
    if (selectingMonth) selectedYear ++;
    else selectedYear += 10;
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div.body {
    overflow: hidden;
    position: relative;
  }

  div.grid {
    display: grid;
    gap: dimensions.$gapSmall;
  }

  div.grid:not(:first-child) {
    position: absolute;
    top: 0;
    left: 0;
  }

  div.grid.month {
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: repeat(3, 1fr);
  }

  div.grid.year {
    grid-template-columns: repeat(5, 1fr);
    grid-template-rows: repeat(2, 1fr);
  }

  div.topRow {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  button.display {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    user-select: none;
    position: relative;
  } 

  button.button {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    border-radius: dimensions.$borderRadiusSmall;
    //color: colors.$foregroundTertiary;
    //background-color: colors.$backgroundTertiary;
    padding: dimensions.$gapSmall;
    cursor: pointer;
    user-select: none;
    position: relative;
    overflow: hidden;
  }

  button.month {
    width: 2em;
  }

  button.year {
    width: 3em;
  }
</style>

<Popup bind:showPopup={internalShow} bind:hidePopup={internalClose} bind:visible={popupVisible}>
  <div class="topRow">
    <IconButton click={prev}>
      <ChevronLeft/>
    </IconButton>
    <button
      class="display"
      type="button"
      onclick={() => {selectingMonth = !selectingMonth}}
      use:focusIndicator={{ type: "underline" }}
    >
      {#if selectingMonth}
        {selectedYear}
      {:else}
        {decadeStart} - {decadeStart + 9}
      {/if}
    </button>
    <IconButton click={next}>
      <ChevronRight/>
    </IconButton>
  </div>
  <div class="body">
    {#if selectingMonth}
      {#each [ selectedYear ] as _ (viewIteration)}
        <div
          class="grid month"
          in:svelteFlyInHorizontal={{duration: 500}}
          out:svelteFlyOutHorizontal={{duration: 500}}
        >
          {#each Array(12) as _, i}
            <button
              class="button month"
              type="button"
              onclick={(e) => clickMonth(e, i)}
              use:focusIndicator
            >
              {getMonthName(i).substring(0, 3)}
            </button>
          {/each}
        </div>
      {/each}
    {:else}
      {#each [ selectedYear ] as _ (viewIteration)}
        <div
          class="grid year"
          in:svelteFlyInHorizontal={{duration: 500}}
          out:svelteFlyOutHorizontal={{duration: 500}}
        >
          {#each Array(10) as _, i}
            <button
              class="button year"
              type="button"
              onclick={(e) => clickYear(e, i)}
              use:focusIndicator
            >
              {decadeStart + i}
            </button>
          {/each}
        </div>
      {/each}
    {/if}
  </div>
</Popup>