<script lang="ts">
  import Checkbox from "../interactive/Checkbox.svelte";

  import { NoOp } from "$lib/client/placeholders";

  interface Props {
    value?: boolean;
    description: string;
    name: string;
    editable?: boolean;
    onChange?: (value: boolean) => any;
  }

  let {
    value = $bindable(false),
    description,
    name,
    editable = true,
    onChange = NoOp,
  }: Props = $props();

  let click: (e: MouseEvent | KeyboardEvent) => void = $state(() => {});
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    align-items: start;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: $gapSmall;
    align-items: center;
    justify-content: start;
    cursor: pointer;
    width: max-content;
  }
  
  label {
    cursor: pointer;
    width: max-content;
  }
</style>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<div
  onclick={click}
  role="checkbox"
  tabindex="-1"
  aria-checked={value}
>
  <Checkbox
    bind:value
    name={name}
    onChange={onChange}
    enabled={editable}
    bind:toggle={click}
  />
  <label for={name}>
    {description}
  </label>
</div>