<script lang="ts">
  import {  setContext } from "svelte";
  import { ChevronLeft, ChevronRight } from "lucide-svelte";

  import IconButton from "../interactive/IconButton.svelte";
  import Popup from "./Popup.svelte";

  import { AsyncNoOp, NoOp } from '$lib/client/placeholders';
  import { focusIndicator } from "$lib/client/decoration";
  import { getMonthName } from "$lib/common/humanization";
  import { svelteFlyInHorizontal, svelteFlyOutHorizontal } from "$lib/client/animations";
  import { getSettings } from "../../lib/client/data/settings.svelte";
  import { UserSettingKeys } from "../../types/settings";
  import { ColorKeys } from "../../types/colors";
  import { t } from "@sveltia/i18n";

  interface Props {
    date: Date;
    anchor?: HTMLElement | undefined;
    showPopup?: () => Promise<void>;
    hidePopup?: () => void;
    onSelect?: (date: Date) => void;
  }

  let {
    date = $bindable(new Date()),
    anchor = undefined,
    showPopup = $bindable(AsyncNoOp),
    hidePopup = $bindable(NoOp),
    onSelect = NoOp,
  }: Props = $props();

  const settings = getSettings();

  let selectingMonth: boolean = $state(true);

  let internalShow: () => Promise<void> = $state(AsyncNoOp);
  let internalClose: () => void = $state(NoOp);

  /* Popup */
  showPopup = () => {
    selectedYear = date.getFullYear();
    selectingMonth = true;
    return internalShow();
  }

  hidePopup = () => {
    internalClose();
    onSelect(date);
  }

  /* Animation */
  let viewIteration = $state(0);
  let flyDirection = $state("left");

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

<Popup bind:showPopup={internalShow} bind:hidePopup={internalClose} tooltip={false} anchor={anchor}>
  <div class="topRow">
    <IconButton onClick={prev} alt={t("button.month.previous")} color={ColorKeys.Accent}>
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
    <IconButton onClick={next} alt={t("button.month.next")} color={ColorKeys.Accent}>
      <ChevronRight/>
    </IconButton>
  </div>
  <div class="body">
    {#if selectingMonth}
      {#if settings.userSettings[UserSettingKeys.AnimateMonthSelectionSwipe]}
        {#each [ selectedYear ] as _ (viewIteration)}
          {@render monthGrid(true)}
        {/each}
      {:else}
        {@render monthGrid(false)}
      {/if}
    {:else}
      {#if settings.userSettings[UserSettingKeys.AnimateMonthSelectionSwipe]}
        {#each [ selectedYear ] as _ (viewIteration)}
          {@render yearGrid(true)}
        {/each}
      {:else}
        {@render yearGrid(false)}
      {/if}
    {/if}
  </div>
</Popup>


{#snippet monthGrid(animate: boolean)}
  <div
    class="grid month"
    class:animate={animate}
    in:svelteFlyInHorizontal={{duration: animate ? 500 * settings.userSettings[UserSettingKeys.AnimationDuration] : 0, flyDirection: () => flyDirection}}
    out:svelteFlyOutHorizontal={{duration: animate ? 500 * settings.userSettings[UserSettingKeys.AnimationDuration] : 0, flyDirection: () => flyDirection}}
  >
    {#each Array(12) as _, i}
      <button
        class="button month"
        type="button"
        onclick={(e) => clickMonth(e, i)}
        use:focusIndicator
      >
        {getMonthName(i, true)}
      </button>
    {/each}
  </div>
{/snippet}

{#snippet yearGrid(animate: boolean)}
  <div
    class="grid year"
    class:animate={animate}
    in:svelteFlyInHorizontal={{duration: animate ? 500 * settings.userSettings[UserSettingKeys.AnimationDuration] : 0, flyDirection: () => flyDirection}}
    out:svelteFlyOutHorizontal={{duration: animate ? 500 * settings.userSettings[UserSettingKeys.AnimationDuration] : 0, flyDirection: () => flyDirection}}
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
{/snippet}