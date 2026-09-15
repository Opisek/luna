<script lang="ts">
  import { CircleAlert, Info } from "lucide-svelte";

  import type { Snippet } from "svelte";
  import Popup from "../popups/Popup.svelte";
  import { AsyncNoOp, NoOp } from "../../lib/client/placeholders";
  import { passIfEnter } from "$lib/common/inputs";
  import { extendUniqueElementId } from "$lib/common/dom";
  import { t } from "@sveltia/i18n";

  interface Props {
    error?: boolean;
    children?: Snippet;
    icon?: Snippet;
    tight?: boolean;
    tiny?: boolean;
    inline?: boolean;
    inheritColor?: boolean;
    pointerCursor?: boolean;
    describes: string;
  }

  let {
    error = false,
    tight = false,
    tiny = false,
    inline = false,
    inheritColor = false,
    pointerCursor = false,
    describes,
    children,
    icon,
  }: Props = $props();

  let showPopover = $state(AsyncNoOp);
  let hidePopover = $state(NoOp);

  let visible = $state(false);
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";
  @use "../../styles/text.scss";

  div.button {
    position: relative;
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    cursor: help;
    display: flex;
    justify-content: center;
    outline: 0;
    padding: dimensions.$gapSmaller;
  }

  div.button.error {
    color: colors.$backgroundFailure;
  }

  div.button.tight {
    padding: 0
  }

  div.button.inline {
    display: inline-flex;
    vertical-align: bottom;
    padding: 0;
    margin-bottom: 0.25ch;
  }

  div.button.inheritColor {
    color: inherit;
  }

  div.button.pointerCursor {
    cursor: pointer;
  }
  
  div.icon {
    display: contents;
  }
</style>

<svelte:window
  onresize={hidePopover}
/>

<div
  class="button"
  class:error={error}
  class:tight={tight}
  class:inline={inline}
  class:inheritColor={inheritColor}
  class:pointerCursor={pointerCursor}
  role="button"
  tabindex="0"
  onmouseenter={() => showPopover().catch(NoOp)}
  onmouseleave={hidePopover}
  onclick={() => showPopover().catch(NoOp)}
  onfocus={() => showPopover().catch(NoOp)}
  onkeypress={(x) => passIfEnter(x, () => showPopover().catch(NoOp))}
  onblur={hidePopover}
  aria-describedby={extendUniqueElementId(["tooltip"], describes)}
  aria-pressed={visible}
  aria-label={t("button.tooltip")}
>
  {#if icon}
    {@render icon?.()}
  {:else if error}
    <CircleAlert size={tiny ? 14 : 16}/>
  {:else}
    <Info size={tiny ? 14 : 16}/>
  {/if}

  <Popup
    bind:showPopup={showPopover}
    bind:hidePopup={hidePopover}
    bind:visible={visible}
    describes={describes}
  >
    <div class="icon" role="presentation" aria-hidden="true">
      {@render children?.()}
    </div>
  </Popup>
</div>