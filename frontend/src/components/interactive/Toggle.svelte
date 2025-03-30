<script lang="ts">
  import { focusIndicator } from "$lib/client/decoration";

  interface Props {
    value: boolean;
    name: string;
    enabled?: boolean;
    onChange?: (value: boolean) => any;
    toggle: (e: MouseEvent | KeyboardEvent) => void;
  }

  let {
    value = $bindable(),
    name,
    enabled = true,
    onChange = () => {},
    toggle = $bindable(),
  }: Props = $props();

  toggle = (e: MouseEvent | KeyboardEvent) => {
    value = !value;
    onChange(value);
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

    width: 2 * text.$lineHeightParagraph;
    height: text.$lineHeightParagraph;

    position: relative;

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;

    border-radius: 0.5 * text.$lineHeightParagraph;
    background-color: colors.$backgroundSecondaryActive;
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
    border-radius: 0.5 * text.$lineHeightParagraph;

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
    all: unset;
    position: absolute;
    left: 0;
    top: 0;
  }
</style>

<!-- Components that use this toggle all implement for={name} -->
<!-- svelte-ignore a11y_consider_explicit_label -->
<button
  type="button"
  class:disabled={!enabled}
  class:check={value}
  onclick={toggle}
  use:focusIndicator
>
  <div
    class="handle"
    class:check={value}
  >
  </div>
  <input type="hidden" name={name} value={value}>
</button>