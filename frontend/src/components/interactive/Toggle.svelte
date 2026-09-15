<script lang="ts">
  import { focusIndicator } from "$lib/client/decoration";
  import { extendUniqueElementId } from "$lib/common/dom";

  interface Props {
    value: boolean;
    name: string;
    id: string
    hasDetails?: boolean;
    enabled?: boolean;
    onChange?: (value: boolean) => any;
  }

  let {
    value = $bindable(),
    name,
    id,
    hasDetails = false,
    enabled = true,
    onChange = () => {},
  }: Props = $props();

  function toggle(e: MouseEvent | KeyboardEvent) {
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
</style>

<!-- Components that use this toggle all implement for={name} -->
<button
  type="button"
  class:disabled={!enabled}
  class:check={value}
  id={id}
  role="checkbox"
  aria-checked={value}
  aria-labelledby={extendUniqueElementId(["label"], id)}
  aria-describedby={hasDetails ? extendUniqueElementId(["tooltip"], id) : undefined}
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