<script lang="ts">
  import LeftIcon from "lucide-svelte/icons/chevron-left";
  import RightIcon from "lucide-svelte/icons/chevron-right";

  import IconButton from "./IconButton.svelte";
  import MonthPopup from "../popups/MonthPopup.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMonthName } from "$lib/common/humanization";

  interface Props {
    date: Date;
    onSelect?: (date: Date) => void;
  }

  let {
    date = $bindable(new Date()),
    onSelect = NoOp,
  }: Props = $props();

  let showPopup: () => any = $state(NoOp);

  function previousMonth() {
    date = new Date(date.getFullYear(), date.getMonth() - 1, date.getDate());
    onSelect(date);
  }

  function nextMonth() {
    date = new Date(date.getFullYear(), date.getMonth() + 1, date.getDate());
    onSelect(date);
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
    {`${getMonthName(date.getMonth())} ${date.getFullYear()}`}
  </button>
  <MonthPopup bind:showPopup bind:date={date} onSelect={onSelect}/>
</div>