<script lang="ts">
  import { ChevronLeft, ChevronRight } from "lucide-svelte";

  import IconButton from "../interactive/IconButton.svelte";
  import Popup from "./Popup.svelte";

  import { NoOp } from '$lib/client/placeholders';
  import { focusIndicator } from "$lib/client/decoration";
  import { getMonthName } from "$lib/common/humanization";

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
    if (selectingMonth) selectedYear --;
    else selectedYear -= 10;
  }

  function next() {
    if (selectingMonth) selectedYear ++;
    else selectedYear += 10;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.grid {
    display: grid;
    gap: $gapSmall;
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
    border-radius: $borderRadiusSmall;
    color: $foregroundSecondary;
    background-color: $backgroundSecondary;
    padding: $gapSmall;
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
  {#if selectingMonth}
  <div class="grid month">
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
  {:else}
  <div class="grid year">
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
  {/if}
</Popup>