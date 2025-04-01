<script lang="ts">
  import Toggle from "../interactive/Toggle.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    value?: boolean;
    description: string;
    info?: string;
    name: string;
    editable?: boolean;
    onChange?: (value: boolean) => any;
  }

  let {
    value = $bindable(false),
    description,
    info,
    name,
    editable = true,
    onChange = NoOp,
  }: Props = $props();

  let click: (e: MouseEvent | KeyboardEvent) => void = $state(() => {});
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    align-items: start;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: dimensions.$gapSmall;
    align-items: center;
    justify-content: start;
    cursor: pointer;
    width: max-content;
  }
  
  label {
    cursor: pointer;
    width: max-content;
    user-select: none;
    color: colors.$foregroundDim;
  }
</style>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
  onclick={click}
  role="checkbox"
  tabindex="-1"
  aria-checked={value}
>
  <Toggle
    bind:value
    name={name}
    onChange={onChange}
    enabled={editable}
    bind:toggle={click}
  />
  <label for={name}>
    {description}
  </label>
  {#if info}
    <Tooltip tight={true}>
      {info}
    </Tooltip>
  {/if}
</div>