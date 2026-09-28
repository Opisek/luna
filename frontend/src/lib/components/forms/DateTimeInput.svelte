<script lang="ts">
  import DateModal from "../modals/DateModal.svelte";
  import Label from "./Label.svelte";
  import TimeModal from "../modals/TimeModal.svelte";

  import { NoOp } from "$lib/scripts/client/placeholders";
  import { focusIndicator } from "$lib/scripts/client/decoration";
  import { t } from "@sveltia/i18n";
  import { getSettings } from "$lib/scripts/client/data/settings.svelte";
  import { UserSettingKeys } from "$lib/types/settings";

  interface Props {
    value: Date | null | undefined;
    allDay?: boolean;
    placeholder: string;
    name: string;
    editable: boolean;
    wrap?: boolean;
    onChange?: (value: Date) => void;
  }

  let {
    value = $bindable(),
    allDay = false,
    placeholder,
    name,
    editable,
    wrap = false,
    onChange = NoOp
  }: Props = $props();

  let dateLocale = $derived(getSettings().userSettings[UserSettingKeys.DateLocale]);

  let dateButton: HTMLButtonElement | null = $state(null);
  let timeButton: HTMLButtonElement | null = $state(null);

  let showDateModal: (initial: Date) => Promise<Date> = $state(Promise.reject);
  let showTimeModal: (initial: Date) => Promise<Date> = $state(Promise.reject);

  async function dateClick(e: MouseEvent | KeyboardEvent) {
    if (!editable) return;
    await showDateModal(value || new Date()).then((result) => {
      value = result;
      onChange(result);
    }).catch(NoOp).finally(() => {
      if (dateButton && e.detail !== 0) {
        dateButton.blur();
      }
    });
  }

  async function timeClick(e: MouseEvent | KeyboardEvent) {
    if (!editable) return;
    await showTimeModal(value || new Date()).then((result) => {
      value = result;
      onChange(result);
    }).catch(NoOp).finally(() => {
      if (timeButton && e.detail !== 0) {
        timeButton.blur();
      }
    });
  }
</script>

<style lang="scss">
  @use "$lib/styles/animations.scss";
  @use "$lib/styles/colors.scss";
  @use "$lib/styles/dimensions.scss";
  @use "$lib/styles/text.scss";

  div.row {
    font-family: text.$fontFamilyTime;
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapSmall;
    margin: dimensions.$gapSmall;
  }

  div.row.editable {
    margin: 0;
  }

  button {
    all: unset;
    border-radius: dimensions.$borderRadius;
    cursor: text;
    transition: padding animations.$animationSpeedFast linear, border-radius animations.$animationSpeedFast linear;
    padding: dimensions.$gapSmall;
    margin: -(dimensions.$gapSmall);
    position: relative;
    overflow: hidden;
  }

  div.row.editable button {
    color: colors.$foregroundSecondary;
    background: colors.$backgroundSecondary;
    cursor: pointer;
    margin: 0;
  }

  div.wrapper {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapMiddle;
  }
</style>

{#if wrap}
  <div class="wrapper">
    {@render inputSnippet()}
  </div>
{:else}
  {@render inputSnippet()}
{/if}

{#snippet inputSnippet()}
  <Label name={name}>{placeholder}</Label>
  <div class="row" class:editable={editable}>
    <button
      bind:this={dateButton}
      onclick={dateClick}
      type="button"
      tabindex={editable ? 0 : -1}
      use:focusIndicator
    >
      {t("date.formatted.date", { values: { date: value || new Date() } })}
    </button>
    {#if !allDay}
      <button
        bind:this={timeButton}
        onclick={timeClick}
        type="button"
        tabindex={editable ? 0 : -1}
        use:focusIndicator
      >
        {t("date.formatted.time", { values: { date: value || new Date() } })}
      </button>
    {/if}
  </div>
{/snippet}

<DateModal bind:showModal={showDateModal}/>
<TimeModal bind:showModal={showTimeModal}/>