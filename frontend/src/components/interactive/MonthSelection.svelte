<script lang="ts">
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

  import IconButton from "./IconButton.svelte";
  import MonthPopup from "../popups/MonthPopup.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMonthName } from "$lib/common/humanization";

  interface Props {
    month: number;
    year: number;
  }

  let { month = $bindable(), year = $bindable() }: Props = $props();

  let showPopup: () => any = $state(NoOp);

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
    width: max-content;
    user-select: none;
  }

  button {
    all: unset;
    cursor: pointer;
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
  <button onclick={showPopup} type="button" use:focusIndicator={{ type: "underline", ignoreParent: true }}>
    {`${getMonthName(month)} ${year}`}
  </button>
  <MonthPopup bind:showPopup bind:year={year} bind:month={month}/>
</div>