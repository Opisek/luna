<script lang="ts">
  import Toggle from "../interactive/Toggle.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";

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

  let uniqueId = $props.id();
  let checkboxId = $derived(generateUniqueElementId(["checkbox", name], uniqueId));
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
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
  }
</style>

<div tabindex="-1">
  <Toggle
    bind:value
    name={name}
    id={checkboxId}
    onChange={onChange}
    enabled={editable}
    hasDetails={info !== undefined}
  />
  <label
    for={checkboxId}
    id={extendUniqueElementId(["label"], checkboxId)}
  >
    {description}
  </label>
  {#if info}
    <Tooltip
      tight={true}
      describes={checkboxId}
    >
      {info}
    </Tooltip>
  {/if}
</div>