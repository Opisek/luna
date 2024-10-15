<script lang="ts">
  import SelectButtons from "../forms/SelectButtons.svelte";
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  export let date: Date;
  export let dateCopy: Date = new Date();

  export let onChange: (date: Date) => void = () => {};

  let pickingHour: boolean;
  let amPm: string;

  let hourInput: HTMLInputElement;
  let minuteInput: HTMLInputElement;

  $: if (amPm === "am") {
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

  export const showModal = () => {
    dateCopy = new Date(date);
    //dateCopy.setHours(0);
    //dateCopy.setMinutes(0);
    if (dateCopy.getHours() > 12 || dateCopy.getHours() === 0) {
      amPm = "pm";
    } else {
      amPm = "am";
    }
    pickingHour = true;
    setTimeout(showModalInternal, 0);
  };

  let showModalInternal: () => any;
  let hideModalInternal: () => any;

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

<Modal title="Pick Time" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <div class="time">
    <span class="time" class:selecting={pickingHour}>
      <input
        bind:this={hourInput}
        type="numeric"
        min="0"
        max="23"
        value={dateCopy.getHours().toString().padStart(2, "0")}
        on:change={() => {
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
          minuteInput.focus();
        }}
        on:input={() => {
          if (Number.parseInt(hourInput.value) >= 2) {
            minuteInput.focus();
          }
        }}
        on:focusin={() => {
          pickingHour = true;
          hourInput.value = "";
        }}
        on:focusout={() => {
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
        on:change={() => {
          let minutes = parseInt(minuteInput.value);
          if (minutes < 0 || minutes > 59) {
            minutes = 0;
          }
          minuteInput.value = minutes.toString().padStart(2, "0");

          dateCopy.setMinutes(minutes);
          dateCopy = dateCopy;

          dateSelected();
        }}
        on:input={() => {
          if (minuteInput.value === "") {
            hourInput.focus();
          } else if (Number.parseInt(minuteInput.value) >= 6) {
            minuteInput.blur();
          }
        }}
        on:keydown={(e) => {
          if ((e.key === "Backspace" || e.key === "Delete") && minuteInput.value === "") {
            hourInput.focus();
          }
        }}
        on:focusin={() => {
          pickingHour = false;
          minuteInput.value = "";
        }}
        on:focusout={() => {
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
          <button class="button hour radial-{i}/12" on:click={() => {
            dateCopy.setHours(((i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12)) % 24);
            dateCopy = dateCopy;
            pickingHour = false;
          }}>
          {((i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12)) % 24}
          </button>
        {:else}
          <button class="button hour radial-{i}/12" on:click={() => {
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
  <svelte:fragment slot="buttons">
    <Button onClick={dateSelected} color="success">Confirm</Button>
    <Button onClick={hideModalInternal} color="failure">Cancel</Button>
  </svelte:fragment>
</Modal>