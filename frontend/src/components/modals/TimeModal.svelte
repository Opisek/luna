<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";

  import { NoOp } from '$lib/client/placeholders';

  interface Props {
    date: Date;
    dateCopy?: Date;
    onChange?: (date: Date) => void;
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    date = $bindable(),
    onChange = NoOp,
    showModal = $bindable(),
    hideModal = $bindable()
  }: Props = $props();

  let dateCopy: Date = $state(date);

  let pickingHour: boolean = $state(true);
  let goBackToHour: boolean = $state(false);
  let amPm: string = $state("am");

  let hourInput: HTMLInputElement;
  let minuteInput: HTMLInputElement;

  $effect(() => {
    if (amPm === "am") {
      dateCopy.setHours(dateCopy.getHours() % 12);
      if (dateCopy.getHours() === 0) {
        dateCopy.setHours(12);
      }
      dateCopy = dateCopy;
    } else {
      dateCopy.setHours(dateCopy.getHours() % 12 + 12);
      if (dateCopy.getHours() === 12) {
        dateCopy.setHours(24);
      }
      dateCopy = dateCopy;
    }
  });

  let showModalInternal: () => any = $state(NoOp);
  let hideModalInternal: () => any = $state(NoOp);

  showModal = () => {
    dateCopy = new Date(date);
    if (dateCopy.getHours() > 12 || dateCopy.getHours() === 0) {
      amPm = "pm";
    } else {
      amPm = "am";
    }
    pickingHour = true;
    goBackToHour = false;
    setTimeout(showModalInternal, 0);
  };

  hideModal = () => {
    hideModalInternal();
  };

  function dateSelected() {
    date = dateCopy;
    hideModalInternal();
    onChange(date);
  }
</script>

<style lang="scss">
  @use "sass:math";

  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  div.time {
    font-size: $fontSizeLarge;
    font-family: $fontFamilyTime;
    display: flex;
    justify-content: center;
    gap: $gapSmaller;
  }

  span.time {
    width: 1.25em;
    text-align: center;
    position: relative;
  }

  input {
    all: unset;
    width: 100%;
  }

  span.selecting::after {
    visibility: visible;
    position: absolute;
    display: inline-block;
    bottom: 0;
    left: 0;
    background-color: $foregroundPrimary;
    width: 100%;
    height: $borderWidth;
    border-radius: calc($borderWidth / 2);
    content: "";
  }

  div.clock {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    margin-top: 5em;
    margin-bottom: 7em + $paddingSmaller;
    font-family: $fontFamilyTime;
  }

  button {
    all: unset;
    background-color: $backgroundSecondary;
    color: $foregroundSecondary;
    //border-radius: $borderRadiusSmall;
    border-radius: 50%;
    padding: $paddingSmaller;
    cursor: pointer;
  }

  @for $i from 0 through 11 {
    $angle: 30deg * $i;

    .radial-#{$i}\/12 {
      position: absolute;
      top: 0;
      left: calc(50% - 1em);
      transform: translate(5em * math.sin($angle), -5em * math.cos($angle));
      width: 1.25em;
      text-align: center;
      height: 1.25em;
    }
  }
</style>

<Modal title="Pick Time" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal} focusElement={hourInput}>
  <div class="time">
    <span class="time" class:selecting={pickingHour}>
      <input
        bind:this={hourInput}
        type="numeric"
        min="0"
        max="23"
        value={dateCopy.getHours().toString().padStart(2, "0")}
        onchange={() => {
          let hours = parseInt(hourInput.value);
          if (hours < 0 || hours > 23) {
            hours = 0;
          }
          hourInput.value = hours.toString().padStart(2, "0");

          dateCopy.setHours(hours);
          dateCopy = dateCopy;

          if (hours > 12 || hours === 0) {
            amPm = "pm";
          } else {
            amPm = "am";
          }

          pickingHour = false;
          goBackToHour = false;
          minuteInput.focus();
        }}
        oninput={() => {
          if (Number.parseInt(hourInput.value) >= 3) {
            minuteInput.focus();
          }
        }}
        onfocusin={() => {
          pickingHour = true;
          if (!goBackToHour || hourInput.value.length >= 1 && hourInput.value[0] === "0") hourInput.value = "";
          else hourInput.value = hourInput.value.substring(0, hourInput.value.length);
          goBackToHour = false;
        }}
        onfocusout={() => {
          if (hourInput.value === "") {
            hourInput.value = dateCopy.getHours().toString().padStart(2, "0");
          }
        }}
      />
    </span>
    <span>
      :
    </span>
    <span class="time" class:selecting={!pickingHour}>
      <input
        bind:this={minuteInput}
        type="numeric"
        min="0"
        max="23"
        value={dateCopy.getMinutes().toString().padStart(2, "0")}
        onchange={() => {
          let minutes = parseInt(minuteInput.value);
          if (minutes < 0 || minutes > 59) {
            minutes = 0;
          }
          minuteInput.value = minutes.toString().padStart(2, "0");

          dateCopy.setMinutes(minutes);
          dateCopy = dateCopy;

          dateSelected();
        }}
        oninput={() => {
          if (Number.parseInt(minuteInput.value) >= 6) {
            minuteInput.blur();
          }
        }}
        onkeydown={(e) => {
          if ((e.key === "Backspace" || e.key === "Delete") && minuteInput.value === "") {
            goBackToHour = true;
            hourInput.focus();
          }
        }}
        onfocusin={() => {
          pickingHour = false;
          minuteInput.value = "";
        }}
        onfocusout={() => {
          if (minuteInput.value === "") {
            minuteInput.value = dateCopy.getMinutes().toString().padStart(2, "0");
          }
        }}
      />
    </span>
  </div>
  <div class="clock">
      {#each Array(12) as _, i}
        {#if pickingHour}
          <button class="button hour radial-{i}/12" tabindex="-11" onclick={() => {
            dateCopy.setHours(((i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12)) % 24);
            dateCopy = dateCopy;
            pickingHour = false;
            minuteInput.focus();
          }}>
          {((i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12)) % 24}
          </button>
        {:else}
          <button class="button hour radial-{i}/12" tabindex="-1" onclick={() => {
            dateCopy.setMinutes(i * 5);
            dateCopy = dateCopy;
            dateSelected();
          }}>
          {i * 5}
          </button>
        {/if}
      {/each}
  </div>
  <SelectButtons bind:value={amPm} name="AM/PM" placeholder="AM/PM" editable={true} options={[{name: "AM", value: "am"}, {name: "PM", value: "pm"}]} label={false}/>
  {#snippet buttons()}
      <Button onClick={dateSelected} color="success">Confirm</Button>
      <Button onClick={hideModalInternal} color="failure">Cancel</Button>
  {/snippet}
</Modal>