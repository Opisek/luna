<script lang="ts" generics="T">

  import { NoOp } from "$lib/client/placeholders";
  import RadioToggle from "../interactive/RadioToggle.svelte";
  import type { Option } from "../../types/options";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";

  interface Props {
    value: T | null;
    groupName: string;
    editable?: boolean;
    options: Option<T>[];
    onClick?: (selected: T) => any;
  }

  let {
    value = $bindable(),
    groupName,
    editable = true,
    options,
    onClick = NoOp,
  }: Props = $props();
  let uniqueId = $props.id();
  let radioGroupId = $derived(generateUniqueElementId(["radio"].concat(groupName.toLocaleLowerCase().split(" ")), uniqueId))
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div.group {
    display: contents;
  }

  div.radio {
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

<div class="group" role="radiogroup">
  {#each options as option (option.value)}
    {@const radioId = extendUniqueElementId(option.name.toLocaleLowerCase().split(" "), radioGroupId)}
    <div class="radio">
      <RadioToggle
        bind:selected={value}
        groupName={groupName}
        id={radioId}
        value={option.value}
        enabled={editable}
        onChange={(x) => { if (x !== null) onClick(x); }}
      />
      <label for={radioId} id={extendUniqueElementId(["label"], radioId)}>
        {option.name}
      </label>
    </div>
  {/each}
</div>