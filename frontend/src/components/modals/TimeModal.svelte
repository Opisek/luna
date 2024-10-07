<script lang="ts">
  import SelectButtons from "../forms/SelectButtons.svelte";
import Modal from "./Modal.svelte";

  export let date: Date;
  export let dateCopy: Date = new Date();

  let pickingHour: boolean;
  let amPm: string;

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
    dateCopy.setHours(0);
    dateCopy.setMinutes(0);
    pickingHour = true;
    amPm = "am";
    setTimeout(showModalInternal, 0);
  };

  let showModalInternal: () => any;
  let hideModalInternal: () => any;

  function dateSelected() {
    date = dateCopy;
    hideModalInternal();
  }
</script>

<style lang="scss">
  @use "sass:math";

  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div.time {
    font-size: $fontSizeLarge;
    display: flex;
    justify-content: center;
    gap: $gapSmaller;
  }

  span.time {
    width: 1.25em;
    text-align: center;
    position: relative;
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
    border-radius: $borderWidth / 2;
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
      {dateCopy.getHours().toString().padStart(2, "0")}
    </span>
    <span>
      :
    </span>
    <span class="time" class:selecting={!pickingHour}>
      {dateCopy.getMinutes().toString().padStart(2, "0")}
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
  <SelectButtons bind:value={amPm} name="AM/PM" placeholder="AM/PM" editable={true} options={[{name: "AM", value: "am"}, {name: "PM", value: "PM"}]} label={false}/>
</Modal>