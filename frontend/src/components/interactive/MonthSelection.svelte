<script lang="ts">
  import { getMonthName } from "../../lib/common/humanization";
  import MonthPopup from "../popups/MonthPopup.svelte";
  import IconButton from "./IconButton.svelte";
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

  export let month: number;
  export let year: number;

  let showPopup: () => any;

  function previousMonth() {
    month--;
    if (month === -1) {
      month = 11;
      year--;
    }
  }

  function nextMonth() {
    if (month === 11) year++;
    month = (month + 1) % 12;
  }
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: row;
    gap: $gapSmall;
    align-items: center;
    position: relative;
  }
</style>

<div>
  <IconButton click={previousMonth}>
    <LeftIcon/>
  </IconButton>
  <IconButton click={nextMonth}>
    <RightIcon/>
  </IconButton>
  <span on:click={showPopup}>
    {`${getMonthName(month)} ${year}`}
  </span>
  <MonthPopup bind:show={showPopup}/>
</div>