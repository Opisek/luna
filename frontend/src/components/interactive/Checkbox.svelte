<script lang="ts">
  import { CheckIcon } from "lucide-svelte";

  import { addRipple, focusIndicator } from "$lib/client/decoration";

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
    toggle = $bindable()
  }: Props = $props();

  toggle = (e: MouseEvent | KeyboardEvent) => {
    value = !value;
    onChange(value);
    if (e instanceof MouseEvent) addRipple(e, false);
    e.stopPropagation();
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  button {
    all: unset;

    width: $lineHeightParagraph;
    height: $lineHeightParagraph;

    position: relative;

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;

    border-radius: $borderRadius;
    background-color: $backgroundSecondary;
    overflow: hidden;
  }

  button.check {
    color: $foregroundAccent;
    background-color: $backgroundAccent;
    --barFocusIndicatorColor: #{$barFocusIndicatorColorAlt};
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
    background-color: red;
  }
</style>

<button
  type="button"
  class:check={value}
  class:disabled={!enabled}
  onclick={toggle}
  use:focusIndicator
>
  {#if value}
    <CheckIcon size={16}/>
  {/if}
  <input type="hidden" name={name} value={value}>
</button>