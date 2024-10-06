<script lang="ts">
  import SelectButtons from "../forms/SelectButtons.svelte";
import Modal from "./Modal.svelte";

  export let date: Date;
  export let dateCopy: Date;

  let pickingHour: boolean;
  let amPm: string;

  export const showModal = () => {
    dateCopy = new Date(date);
    pickingHour = true;
    amPm = "am";
    showModalInternal();
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

  div {
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    margin-top: 5em;
    margin-bottom: 6.5em + $paddingSmaller;
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

  @each $div in 12 {
    @for $i from 0 through 11 {
      $angle: calc(360deg * $i / $div);

      .radial-#{$i}\/#{$div} {
        position: absolute;
        top: 0;
        left: calc(50% - 1em);
        transform: translate(5em * math.sin($angle), -5em * math.cos($angle));
        width: 1.25em;
        text-align: center;
        height: 1.25em;
      }
    }
  }
</style>

<Modal title="Pick Time" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <div>
      {#each Array(12) as _, i}
        {#if pickingHour}
          <button class="button hour radial-{i}/12" on:click={() => {
            dateCopy.setHours((i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12));
            pickingHour = false;
          }}>
          {(i == 0 ? 12 : i) + (amPm === "am" ? 0 : 12)}
          </button>
        {:else}
          <button class="button hour radial-{i}/12" on:click={() => {
            dateCopy.setMinutes(i * 5);
            dateSelected();
          }}>
          {i * 5}
          </button>
        {/if}
      {/each}
  </div>
  {#if pickingHour}
  <SelectButtons bind:value={amPm} name="AM/PM" placeholder="AM/PM" editable={true} options={[{name: "AM", value: "am"}, {name: "PM", value: "PM"}]} label={false}/>
  {/if}
</Modal>