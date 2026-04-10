<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import type { Snippet } from "svelte";
  import Popup from "../popups/Popup.svelte";
  import { NoOp } from "../../lib/client/placeholders";

  interface Props {
    error?: boolean;
    children?: Snippet;
    icon?: Snippet;
    tight?: boolean;
    tiny?: boolean;
    inline?: boolean;
    inheritColor?: boolean;
    pointerCursor?: boolean;
    role?: string;
  }

  let {
    error = false,
    tight = false,
    tiny = false,
    inline = false,
    inheritColor = false,
    pointerCursor = false,
    role = "tooltip",
    children,
    icon,
  }: Props = $props();

  let showPopover = $state(NoOp);
  let hidePopover = $state(NoOp);
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/text.scss";

  div {
    position: relative;
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    cursor: help;
    display: flex;
    justify-content: center;
    outline: 0;
    padding: dimensions.$gapSmaller;
  }

  div.error {
    color: colors.$backgroundFailure;
  }

  div.tight {
    padding: 0
  }

  div.inline {
    display: inline-flex;
    vertical-align: bottom;
    padding: 0;
    margin-bottom: 0.25ch;
  }

  div.inheritColor {
    color: inherit;
  }

  div.pointerCursor {
    cursor: pointer;
  }
</style>

<svelte:window
  onresize={hidePopover}
/>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<div
  class:error={error}
  class:tight={tight}
  class:inline={inline}
  class:inheritColor={inheritColor}
  class:pointerCursor={pointerCursor}
  role={role}
  tabindex="-1"
  onmouseenter={showPopover}
  onmouseleave={hidePopover}
  onfocus={hidePopover}
  onblur={hidePopover}
>
  {#if icon}
    {@render icon?.()}
  {:else if error}
    <CircleAlert size={tiny ? 14 : 16}/>
  {:else}
    <Info size={tiny ? 14 : 16}/>
  {/if}

  <Popup bind:showPopup={showPopover} bind:hidePopup={hidePopover}>
    {@render children?.()}
  </Popup>
</div>