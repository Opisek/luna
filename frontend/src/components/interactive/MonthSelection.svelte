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
    granularity?: "month" | "week" | "day";
    onSelect?: (date: Date) => void;
  }

  let {
    date = $bindable(new Date()),
    granularity = "month",
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

  function previousWeek() {
    date = new Date(date.getFullYear(), date.getMonth(), date.getDate() - 7);
    onSelect(date);
  }

  function nextWeek() {
    date = new Date(date.getFullYear(), date.getMonth(), date.getDate() + 7);
    onSelect(date);
  }

  function previousDay() {
    date = new Date(date.getFullYear(), date.getMonth(), date.getDate() - 1);
    onSelect(date);
  }

  function nextDay() {
    date = new Date(date.getFullYear(), date.getMonth(), date.getDate() + 1);
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
  {#if granularity === "month"}
    {@render buttons(previousMonth, nextMonth)}
  {:else if granularity === "week"}
    {@render buttons(previousWeek, nextWeek)}
  {:else if granularity === "day"}
    {@render buttons(previousDay, nextDay)}
  {/if}
  <button onclick={showPopup} type="button" use:focusIndicator={{ type: "underline", ignoreParent: true }}>
    {`${getMonthName(date.getMonth())} ${date.getFullYear()}`}
  </button>
  <MonthPopup bind:showPopup bind:date={date} onSelect={onSelect}/>
</div>

{#snippet buttons(prev: () => void, next: () => void)}
  <IconButton click={prev}>
    <LeftIcon/>
  </IconButton>
  <IconButton click={next}>
    <RightIcon/>
  </IconButton>
{/snippet}