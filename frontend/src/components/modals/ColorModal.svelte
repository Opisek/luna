<script lang="ts">
  import { browser } from "$app/environment";
  import { HSLtoRGB, isValidColor, parseRGB, recommendedColors, RGBtoHSL, serializeRGB } from "../../lib/common/colors";
  import TextInput from "../forms/TextInput.svelte";
  import Button from "../interactive/Button.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import ColorCircle from "../misc/ColorCircle.svelte";
  import Modal from "./Modal.svelte";

  export let color: string | null;
  let currentColor: string;
  let currentHSL: [number, number, number] = [0, 100, 50];

  let pickerActive = false;
  let hueActive = false;

  function setHSLFromColor() {
    let newRGB: [number, number, number];
    if (!currentColor || currentColor === "") {
      newRGB = [255, 0, 0];
    } else {
      newRGB = parseRGB(currentColor);
    }
    currentHSL = RGBtoHSL(newRGB);
  }

  function setColorFromHSL() {
    currentColor = serializeRGB(HSLtoRGB(currentHSL));
  }

  export const showModal = () => {
    pickerActive = false;
    hueActive = false;
    currentColor = (color || "").toUpperCase();
    setHSLFromColor();
    setTimeout(showModalInternal, 0);
  };

  const hideModal = () => {
    pickerActive = false;
    hueActive = false;
    hideModalInternal();
  }

  let showModalInternal: () => any;
  let hideModalInternal: () => any;

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

  function setColor(color: string | null) {
    if (color === null) {
      currentColor = "";
      currentHSL = [0, 100, 50];
    } else {
      color = color.toUpperCase();
      currentColor = color;
      currentHSL = RGBtoHSL(parseRGB(color));
    }
  }

  function codeFocus() {
    if (!currentColor || currentColor.length == 0) {
      currentColor = "#";
    }
  }

  function calculateCoordsRelativeToElement(e: MouseEvent): [number, number] {
    const target = e.target as HTMLElement;

    const rect = target.getBoundingClientRect();

    let x = e.clientX - rect.left;
    let y = e.clientY - rect.top;

    x = x / rect.width * 100;
    y = y / rect.height * 100;

    x = Math.max(0, Math.min(100, x));
    y = Math.max(0, Math.min(100, y));

    return [Math.round(x), Math.round(y)];
  }

  function pickerDown(e: MouseEvent) {
    e.stopPropagation();
    pickerActive = true;
    pickerMove(e);
  }
  function pickerMove(e: MouseEvent) {
    e.stopPropagation();
    if (!pickerActive) return;
    const [x,y] = calculateCoordsRelativeToElement(e);

    const s = 100 - y;
    const minL = 0.5 * s;
    const l = (100 - x) * (100 - minL) / 100 + minL;

    currentHSL[1] = s;
    currentHSL[2] = l;

    setColorFromHSL();
  }
  function pickerUp(e: MouseEvent) {
    e.stopPropagation();
    pickerActive = false
  }

  function hueDown(e: MouseEvent) {
    e.stopPropagation();
    hueActive = true;
    hueMove(e);
  }
  function hueMove(e: MouseEvent) {
    e.stopPropagation();
    if (!hueActive) return;
    const [x,_] = calculateCoordsRelativeToElement(e);

    currentHSL[0] = x / 100 * 360;
    setColorFromHSL();
  }
  function hueUp(e: MouseEvent) {
    e.stopPropagation();
    hueActive = false
  }

  if (browser) {
    window.addEventListener("mouseup", (e) => {
      if (pickerActive) pickerUp(e);
      if (hueActive) hueUp(e);
    });
  }
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div.grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr auto;
    grid-template-areas: "current hSL" "code Hsl";
    gap: $gapSmall;
  }

  div.suggestions {
    display: flex;
    flex-wrap: wrap;
    gap: $gapSmall;
    justify-content: center;
  }

  div.picker {
    grid-area: hSL;
    width: 100%;
    height: 0;
    mask: linear-gradient(270deg, white, transparent);
    padding-bottom: 100%;
    border-radius: $borderRadius;
  }

  div.hue {
    grid-area: Hsl;
    width: 100%;
    flex-grow: 1;
    padding: $gapSmall;
    border-radius: $borderRadius;
  }
</style>

<Modal title="Pick Color" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal}>
  <div class="grid">
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      class="picker"
      style="background: linear-gradient(0deg, hsl({currentHSL[0]} 0 0), hsl({currentHSL[0]} 100% 50%))"
      on:mousedown={pickerDown}
      on:mousemove={pickerMove}
      on:mouseup={pickerUp}
    ></div>
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      class="hue"
      style="background: linear-gradient(90deg in hsl longer hue, hsl(0 {currentHSL[1]}% {currentHSL[2]}%), hsl(360 {currentHSL[1]}% {currentHSL[2]}%))"
      on:mousedown={hueDown}
      on:mousemove={hueMove}
      on:mouseup={hueUp}
    ></div>
    <ColorCircle color={currentColor} size="fill" shape="squircle"/>
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
    <IconButton click={() => setColor(null)}>
      <ColorCircle color={null} size="medium"/>
    </IconButton>
    {#each recommendedColors as color}
      <IconButton click={() => setColor(color)}>
        <ColorCircle color={color} size="medium"/>
      </IconButton>
    {/each}
  </div>
  <svelte:fragment slot="buttons">
    <Button onClick={confirm} color="success">Confirm</Button>
    <Button onClick={cancel} color="failure">Cancel</Button>
  </svelte:fragment>
</Modal>