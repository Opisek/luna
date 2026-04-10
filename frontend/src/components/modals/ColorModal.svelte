<script lang="ts">
  import { browser } from "$app/environment";

  import ColorCircle from "../misc/ColorCircle.svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import Modal from "./Modal.svelte";
  import TextInput from "../forms/TextInput.svelte";

  import { HSLtoRGB, isValidColor, parseRGB, recommendedColorNames, recommendedColors, RGBtoHSL, serializeRGB } from "$lib/common/colors";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    showModal: (initial: string | null) => Promise<string>;
  }

  let {
    showModal = $bindable(),
  }: Props = $props();

  let picker: HTMLElement;
  let hue: HTMLElement;

  let initialColor: string = $state("#000000");
  let currentColor: string = $state("#000000");
  let currentHSL: [number, number, number] = $state([0, 100, 50]);

  let pickerActive = false;
  let hueActive = false;

  let showModalInternal: () => Promise<string> = $state(Promise.reject);
  let success: (color: string) =>  void = $state(NoOp);
  let failure: () => void = $state(NoOp);

  function mouseUp(e: MouseEvent) {
    if (pickerActive) pickerUp(e);
    if (hueActive) hueUp(e);
  }

  function mouseMove(e: MouseEvent) {
    pickerMove(e);
    hueMove(e);
  }

  showModal = (initial: string | null) => {
    pickerActive = false;
    hueActive = false;
    currentColor = (initial || "").toUpperCase();
    initialColor = currentColor;
    setHSLFromColor();

    if (browser) {
      window.addEventListener("mouseup", mouseUp);
      window.addEventListener("mousemove", mouseMove);
    }

    return showModalInternal();
  };

  function onModalHide() {
    pickerActive = false;
    hueActive = false;

    if (browser) {
      window.removeEventListener("mouseup", mouseUp);
      window.removeEventListener("mousemove", mouseMove);
    }
  }

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

  function validateColor() {
    if (!currentColor.startsWith("#")) {
      currentColor = "#" + currentColor;
    }
    if (!isValidColor(currentColor)) {
      if (currentColor == "#") currentColor = "";
      else currentColor = initialColor;
    } else {
      currentColor = currentColor.toUpperCase();
      currentHSL = RGBtoHSL(parseRGB(currentColor));
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
  function codeInput() {
    if (!currentColor || currentColor.length == 0) {
      currentColor = "#";
    } else if (!currentColor.startsWith("#")) {
      currentColor = "#" + currentColor;
    }

    let replacementColor = currentColor; 
    do {
      currentColor = replacementColor;
      replacementColor = replacementColor.replaceAll(/(.+)[^0-9A-Fa-f]/g, "$1");
    } while (replacementColor != currentColor);

    if (currentColor.length > 7) {
      currentColor = currentColor.slice(0, 7);
    }
  }

  function calculateCoordsRelativeToElement(e: MouseEvent, target: HTMLElement): [number, number] {
    if (!target) return [0, 0];

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
    if (pickerActive || hueActive) return;
    pickerActive = true;
    pickerMove(e);
  }
  function pickerMove(e: MouseEvent) {
    if (!pickerActive) return;
    const [x,y] = calculateCoordsRelativeToElement(e, picker);

    const s = 100 - y;
    const minL = 0.5 * s;
    const l = (100 - x) * (100 - minL) / 100 + minL;

    currentHSL[1] = s;
    currentHSL[2] = l;

    setColorFromHSL();
  }
  function pickerUp(e: MouseEvent) {
    pickerActive = false
  }

  function hueDown(e: MouseEvent) {
    if (pickerActive || hueActive) return;
    hueActive = true;
    hueMove(e);
  }
  function hueMove(e: MouseEvent) {
    if (!hueActive) return;
    const [x,_] = calculateCoordsRelativeToElement(e, hue);

    currentHSL[0] = x / 100 * 360;
    setColorFromHSL();
  }
  function hueUp(e: MouseEvent) {
    hueActive = false
  }
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div.grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-template-rows: 1fr auto;
    grid-template-areas: "current hSL" "code Hsl";
    gap: dimensions.$gapSmall;
  }

  div.suggestions {
    display: flex;
    flex-wrap: wrap;
    gap: dimensions.$gapSmall;
    justify-content: center;
  }

  div.picker {
    mask: linear-gradient(270deg, white, transparent);
    cursor: pointer;
    width: 100%;
    height: 100%;
  }
  div.pickerBackground {
    grid-area: hSL;
    background-color: #ffffff;
    width: 100%;
    height: 100%;
    border-radius: dimensions.$borderRadius;
    overflow: hidden;
  }

  div.hue {
    grid-area: Hsl;
    width: 100%;
    flex-grow: 1;
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadius;
    cursor: pointer;
  }
</style>

<Modal title="Pick Color"
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
  onModalHide={onModalHide}
  onModalSubmit={() => success(currentColor)}>
  <div class="grid">
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="pickerBackground">
      <div
        bind:this={picker}
        class="picker"
        style="background: linear-gradient(0deg, hsl({currentHSL[0]} 0 0), hsl({currentHSL[0]} 100% 50%))"
        onmousedown={pickerDown}
        onmousemove={pickerMove}
        onmouseup={pickerUp}
      ></div>
    </div>
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      bind:this={hue}
      class="hue"
      style="background: linear-gradient(90deg in hsl longer hue, hsl(0 {currentHSL[1]}% {currentHSL[2]}%), hsl(360 {currentHSL[1]}% {currentHSL[2]}%))"
      onmousedown={hueDown}
      onmousemove={hueMove}
      onmouseup={hueUp}
    ></div>
    <ColorCircle color={currentColor} size="fill" shape="squircle"/>
    <TextInput
      bind:value={currentColor}
      placeholder="Color"
      name="color"
      editable={true}
      label={false}
      onInput={codeInput}
      onChange={validateColor}
      onFocus={codeFocus}
    />
  </div>
  <div class="suggestions">
    <IconButton onClick={() => setColor(null)} alt="Default color">
      <ColorCircle color={null} size="medium"/>
    </IconButton>
    {#each recommendedColors as color, index}
      <IconButton onClick={() => setColor(color)} alt={recommendedColorNames[index]}>
        <ColorCircle color={color} size="medium"/>
      </IconButton>
    {/each}
  </div>
  {#snippet buttons()}
      <IconButton onClick={() => success(currentColor)} color={ColorKeys.Success} alt="Confirm"><Check/></IconButton>
      <IconButton onClick={failure} color={ColorKeys.Danger} alt="Cancel"><X/></IconButton>
  {/snippet}
</Modal>