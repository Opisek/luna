<script lang="ts">
  import { browser } from "$app/environment";
  import Label from "./Label.svelte";

  interface Props {
    value: number;

    min?: number;
    max?: number;
    step?: number | number[];

    transform?: (number: number) => number;
    detentTransform?: (number: number) => string;

    title: string;
    name: string;
  }

  let {
    value = $bindable(),
    min = 0,
    max = 100,
    step = 1,

    transform = (number: number) => number,
    detentTransform = (number: number) => `${(number - min) / (max-min) * 100}%`,

    title,
    name,
  }: Props = $props();

  let rawValue = $state(value);
  let steppedValue = $derived.by(() => {
    if (typeof step === 'number' && step > 0) {
      return Math.round((rawValue - min) / step) * step + min;
    } else if (Array.isArray(step) && step.length > 0) {
      let distances = step.map((s) => Math.abs(rawValue - s));
      let lowestDistance = Number.MAX_VALUE;
      let closest = step[0];
      for (let i = 0; i < step.length; i++) {
        if (distances[i] < lowestDistance) {
          lowestDistance = distances[i];
          closest = step[i];
        }
      }
      return closest;
    } else {
      return rawValue;
    }
  });
  let transformedValue = $derived(transform(steppedValue));

  let slider = $state<HTMLDivElement | null>(null);
  let handle = $state<HTMLDivElement | null>(null);

  let animationLength = 0.25;
  let animationMultiplier = $state(1);
  let mouseDownSince = $state(0);
  let mouseDownFor = $state(0);

  function mouseDown(e: MouseEvent) {
    if (!slider) return;
    if (!browser) return;

    window.addEventListener("mousemove", adjustValue);
    window.addEventListener("mouseup", mouseUp, { once: true });

    mouseDownSince = new Date().getTime();
    animationMultiplier = 1;

    const oldValue = rawValue;
    adjustValue(e);
    const percentualChange = Math.abs((rawValue - oldValue) / (max - min));
    if (percentualChange < 0.2) animationMultiplier = 0.2;
    else rawValue = transformedValue;
    setTimeout(() => {
      if (mouseDownSince == 0) return;
      if (rawValue != transformedValue) return;
      adjustValue(e);
    }, 150);
  }
  function mouseUp(e: MouseEvent) {
    if (!slider) return;
    window.removeEventListener("mousemove", adjustValue);
    adjustValue(e);
    mouseDownFor = 0;
    mouseDownSince = 0 ;
    rawValue = transformedValue;
    value = transformedValue;
  }
  function adjustValue(e: MouseEvent) {
    if (!slider || !handle) return;
    mouseDownFor = (new Date().getTime() - mouseDownSince) / 1000;
    e.stopPropagation();
    e.preventDefault();
    const sliderRect = slider.getBoundingClientRect();
    const handleRect = handle.getBoundingClientRect();
    const offsetX = e.clientX - sliderRect.left - handleRect.width / 2;
    const percent = offsetX / (sliderRect.width - handleRect.width);
    const newValue = min + (max - min) * percent;
    rawValue = Math.max(min, Math.min(max, newValue));
  }

  let backgroundPercentage = $state(0);
  let backgroundPertentageInterval = $state<ReturnType<typeof setInterval> | null>(null);
  function recalculateBackgroundPercentage() {
    if (!slider || !handle) return 0;
    rawValue;
    const sliderRect = slider.getBoundingClientRect();
    const handleRect = handle.getBoundingClientRect();
    const handleCenter = handleRect.left + handleRect.width / 2;
    backgroundPercentage =  ((handleCenter - sliderRect.left) / sliderRect.width) * 100
  }

  function transitionStart() {
    if (backgroundPertentageInterval != null) return;
    backgroundPertentageInterval = setInterval(() => {
      recalculateBackgroundPercentage();
    }, 10);
  }
  function transitionEnd() {
    if (backgroundPertentageInterval == null || mouseDownSince != 0) return;
    clearTimeout(backgroundPertentageInterval);
    backgroundPertentageInterval = null;
  }

  $effect(() => {
    // @ts-ignore
    rawValue, value, recalculateBackgroundPercentage();
  });
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.container {
    --barSize: 0.4em;
    --handleSize: 1em;
    --detentSize: 1em;
    --detentWidth: 0.2em;
    --backgroundBase: #{colors.$backgroundSecondary};
    --backgroundFilled: #{colors.$backgroundAccent};
    --baseAnimationSpeeed: #{animations.$animationSpeed};

    height: var(--detentSize);
    margin-top: 1.25em;
    margin-bottom: .25em;
    --indent: #{dimensions.$gapSmall};
    width: calc(100% - 2 * var(--indent));
    margin-left: var(--indent);
    position: relative;

    cursor: pointer;
  }

  input {
    width: 100%;
  }

  div.slider {
    width: 100%;
    height: var(--barSize);
    border-radius: calc(var(--barSize) / 2);
    position: absolute;
    pointer-events: none;
    top: calc(var(--detentSize) / 2 - var(--barSize) / 2);
  }

  div.handle {
    position: absolute;
    height: var(--handleSize);
    width: var(--handleSize);
    background-color: colors.$backgroundAccent;
    border-radius: 50%;;
    margin-top: calc(var(--handleSize) / -2);
    top: 50%;
    pointer-events: none;
  }

  div.detents {
    position: absolute;
    width: calc(100% - var(--handleSize));
    left: calc(var(--handleSize) / 2);
    height: var(--detentSize);
  }

  div.detent {
    position: absolute;
    background-color: colors.$backgroundSecondary;
    width: var(--detentWidth);
    border-radius: calc(var(--detentWidth) / 2);
    height: 100%;
    margin-left: calc(var(--detentWidth) / -2);
    pointer-events: none;
  }

  div.detent span {
    position: absolute;
    bottom: calc(50% + var(--detentSize) / 2 + 0.25em);
    font-size: text.$fontSizeSmall;
    left: 0;
    right: 0;
    margin: auto;
    text-align: center;
    margin-left: -2em;
    margin-right: -2em;
  }

  div.detent:first-child span {
    margin-left: -1em !important;
  }

  div.detent:last-child span {
    margin-right: -1em !important;
  }
</style>

<Label name={name}>{title}</Label>
<div class="container"
  onmousedown={mouseDown}
  tabindex="0"
  role="slider"
  aria-valuemin={min}
  aria-valuenow={transformedValue}
  aria-valuemax={max}
>
  <input type="hidden" name={name} bind:value={rawValue} />
  {#if step != 0}
    {@const detents = Array.isArray(step) ? step : Array.from({ length: Math.floor((max - min) / step) + 1 }, (_, i) => min + i * step)}
    {@const showAllDetentLabels = detents.length <= 7}
    <div class="detents">
      {#each detents as detent, i}
        <div class="detent"
          style="left: calc(({(detent - min) / (max - min) * 100}%));"
        >
          {#if i == 0 || i == detents.length - 1 || showAllDetentLabels || i % 2 == 0}
            <span>
              {detentTransform(detent)}
            </span>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
  <div class="slider"
    bind:this={slider}
    style="background: linear-gradient(to right, var(--backgroundFilled) {backgroundPercentage}%, var(--backgroundBase) {backgroundPercentage}%);"
  >
    <div class="handle"
      bind:this={handle}
      ontransitionstart={transitionStart}
      ontransitionend={transitionEnd}
      ontransitioncancel={transitionEnd}
      style={`
        transition: left max(${mouseDownSince == 0 ? animationLength : animationMultiplier * (animationLength - mouseDownFor)}s, 0s) ease-in-out;
        left: calc((100% - var(--handleSize)) * ${(rawValue - min) / (max - min)});
      `}
    >
    </div>
  </div>
</div>