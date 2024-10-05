<script lang="ts">
  import { ChevronLeft, ChevronRight } from "lucide-svelte";
  import { getMonthName } from "../../lib/common/humanization";
  import IconButton from "../interactive/IconButton.svelte";
  import Popup from "./Popup.svelte";
  import { addRipple } from "../../lib/client/decoration";

  export const show = () => {
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
    justify-content: center;
    align-items: space-between;
  }

  button.display {
    all: unset;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    cursor: pointer;
    user-select: none;
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
  }

  button.year {
    width: 3em;
  }
</style>

<Popup bind:show={internalShow} bind:close={internalClose}>
  <div class="topRow">
    <IconButton click={prev}>
      <ChevronLeft/>
    </IconButton>
    {#if selectingMonth}
      <button class="display year" on:click={() => {selectingMonth = false}}>
        {selectedYear}
      </button>
    {:else}
      <button class="display decade" on:click={() => {selectingMonth = true}}>
        {decadeStart} - {decadeStart + 9}
      </button>
    {/if}
    <IconButton click={next}>
      <ChevronRight/>
    </IconButton>
  </div>
  {#if selectingMonth}
  <div class="grid month">
    {#each Array(12) as _, i}
      <button class="button month" on:click={(e) => clickMonth(e, i)}>
        {getMonthName(i).substring(0, 3)}
      </button>
    {/each}
  </div>
  {:else}
  <div class="grid year">
    {#each Array(10) as _, i}
      <button class="button year" on:click={(e) => clickYear(e, i)}>
        {decadeStart + i}
      </button>
    {/each}
  </div>
  {/if}
</Popup>