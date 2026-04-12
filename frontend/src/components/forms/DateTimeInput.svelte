<script lang="ts">
  import DateModal from "../modals/DateModal.svelte";
  import Label from "./Label.svelte";
  import TimeModal from "../modals/TimeModal.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { date, time } from "@sveltia/i18n";

  interface Props {
    value: Date;
    allDay: boolean;
    placeholder: string;
    name: string;
    editable: boolean;
    wrap?: boolean;
    onChange?: (value: Date) => void;
  }

  let {
    value = $bindable(),
    allDay,
    placeholder,
    name,
    editable,
    wrap = false,
    onChange = NoOp
  }: Props = $props();

  let dateButton: HTMLButtonElement | null = $state(null);
  let timeButton: HTMLButtonElement | null = $state(null);

  let showDateModal: (initial: Date) => Promise<Date> = $state(Promise.reject);
  let showTimeModal: (initial: Date) => Promise<Date> = $state(Promise.reject);

  async function dateClick(e: MouseEvent | KeyboardEvent) {
    if (!editable) return;
    await showDateModal(value).then((result) => {
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
    await showTimeModal(value).then((result) => {
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
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

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
      {date(value)}
    </button>
    {#if !allDay}
      <button
        bind:this={timeButton}
        onclick={timeClick}
        type="button"
        tabindex={editable ? 0 : -1}
        use:focusIndicator
      >
        {time(value, { hour: "2-digit", minute: "2-digit" })}
      </button>
    {/if}
  </div>
{/snippet}

<DateModal bind:showModal={showDateModal}/>
<TimeModal bind:showModal={showTimeModal}/>