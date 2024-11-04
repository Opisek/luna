<script lang="ts">
  import { ChevronLeft, ChevronRight } from "lucide-svelte";
  import { getMonthName } from "../../lib/common/humanization";
  import IconButton from "../interactive/IconButton.svelte";
  import Popup from "./Popup.svelte";
  import { focusIndicator } from "../../lib/client/decoration";

  let popupVisible: boolean = false;

  export const show = () => {
    if (popupVisible) return;
    selectedYear = year;
    selectingMonth = true;
    setTimeout(internalShow, 0);
  }
  let internalShow: () => void;
  let internalClose: () => void;

  export let month: number;
  export let year: number;

  let selectedYear: number;
  let decadeStart: number;

  $: decadeStart = Math.floor(selectedYear / 10) * 10;

  let selectingMonth: boolean;

  function clickMonth(e: MouseEvent, i: number) {
    //addRipple(e);
    month = i;
    year = selectedYear;
    internalClose();
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

  let clickedMonth = -1;
  let clickedYear = -1;
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
    padding: $paddingSmaller;
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

<Popup bind:show={internalShow} bind:close={internalClose} bind:visible={popupVisible}>
  <div class="topRow">
    <IconButton click={prev}>
      <ChevronLeft/>
    </IconButton>
    <button
      class="display"
      type="button"
      on:click={() => {selectingMonth = !selectingMonth}}
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
        class:click={clickedMonth === i}
        type="button"
        on:click={(e) => clickMonth(e, i)}
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
        class:click={clickedYear === i}
        type="button"
        on:click={(e) => clickYear(e, i)}
        use:focusIndicator
      >
        {decadeStart + i}
      </button>
    {/each}
  </div>
  {/if}
</Popup>