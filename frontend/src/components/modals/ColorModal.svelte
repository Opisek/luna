<script lang="ts">
  import { isValidColor, recommendedColors } from "../../lib/common/colors";
  import TextInput from "../forms/TextInput.svelte";
  import Button from "../interactive/Button.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Modal from "./Modal.svelte";

  export let color: string | null;
  let currentColor: string;

  export const showModal = () => {
    currentColor = color || "";
    setTimeout(showModalInternal, 0);
  };

  let showModalInternal: () => any;
  let hideModal: () => any;

  function confirm() {
    color = currentColor;
    hideModal();
  }

  function cancel() {
    hideModal();
  }

  function validateColor() {
    if (!isValidColor(currentColor)) {
      currentColor = color || "";
      return;
    } else {
      // @ts-ignore currentColor cant't be null due to isValidColor check
      currentColor = currentColor.toUpperCase();
    }
  }

  function codeFocus() {
    if (!currentColor || currentColor.length == 0) {
      currentColor = "#";
    }
  }
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr auto auto;
    grid-template-areas: "current input" "code input" "code input";
    gap: $gapSmall;
  }

  div.suggestions {
    display: flex;
    flex-wrap: wrap;
    gap: $gapSmall;
    justify-content: center;
  }
</style>

<Modal title="Pick Color" bind:showModal={showModalInternal} bind:hideModal={hideModal}>
  <div class="grid">
    <ColorCircle color={currentColor} size="fill" shape="squircle"/>
    <div style="grid-area:input">
      TODO: color picker circle here
    </div>
    <TextInput
      bind:value={currentColor}
      placeholder="Color"
      name="color"
      editable={true}
      label={false}
      onChange={validateColor}
      onFocus={codeFocus}
    />
  </div>
  <div class="suggestions">
    {#each recommendedColors as color}
      <IconButton click={() => {currentColor = color}}>
        <ColorCircle color={color} size="medium"/>
      </IconButton>
    {/each}
  </div>
  <svelte:fragment slot="buttons">
    <Button onClick={confirm} color="success">Confirm</Button>
    <Button onClick={cancel} color="failure">Cancel</Button>
  </svelte:fragment>
</Modal>