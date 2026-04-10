<script lang="ts" generics="T">
  import { focusIndicator } from "$lib/client/decoration";

  interface Props {
    name: string;
    value: T;
    selected: T | null;
    enabled?: boolean;
    onChange?: (value: T | null) => any;
  }

  let {
    name,
    value,
    selected = $bindable(),
    enabled = true,
    onChange = () => {},
  }: Props = $props();

  let checked = $derived(value == selected);

  const toggle = (e: MouseEvent | KeyboardEvent) => {
    if (checked) selected = null;
    else selected = value;
    onChange(selected);
    e.stopPropagation();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  button {
    all: unset;

    width: calc(2 * text.$lineHeightParagraph);
    height: text.$lineHeightParagraph;

    position: relative;

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;

    border-radius: calc(0.5 * text.$lineHeightParagraph);
    background-color: colors.$backgroundTertiary;
    overflow: hidden;
  }

  button.check {
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt} !important;
  }

  button::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;

    background-color: colors.$backgroundAccent;
    border-radius: calc(0.5 * text.$lineHeightParagraph);

    transition: transform animations.$animationSpeed;

    transform: scale(0);
  }

  button.check::after {
    transform: none;
  }

  .handle {
    height: calc(100% - dimensions.$gapSmaller);
    aspect-ratio: 1/1;
    border-radius: 50%;
    background-color: colors.$backgroundPrimary;

    left: dimensions.$gapSmaller;

    position: absolute;

    transform: scale(0.75);

    transition: transform animations.$animationSpeed animations.$cubic;

    z-index: 2;
  }

  .handle.check {
    transform: translateX(100%);
  }

  button.disabled {
    cursor: unset;
  }

  button :global(*) {
    pointer-events: none;
  }

  input {
    display: none;
  }
</style>

<!-- Components that use this toggle all implement for={name} -->
<!-- svelte-ignore a11y_consider_explicit_label -->
<button
  type="button"
  class:disabled={!enabled}
  class:check={checked}
  onclick={toggle}
  use:focusIndicator
>
  <div
    class="handle"
    class:check={checked}
  >
  </div>
  <input type="radio" id={`${name}-${value}`} name={name} value={value} bind:group={selected}>
</button>