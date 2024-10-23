<script lang="ts">
  import Checkbox from "../interactive/Checkbox.svelte";

  export let value: boolean = false;
  export let description: string;
  export let name: string;

  export let editable: boolean = true;

  export let onChange: (value: boolean) => any = () => {};

  let click: () => any;
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

<div
  on:click={click}
  on:keypress={click}
  role="checkbox"
  tabindex="-1"
  aria-checked={value}
>
  <Checkbox
    bind:value
    name={name}
    onChange={onChange}
    enabled={editable}
    bind:click={click}
  />
  <label for={name}>
    {description}
  </label>
</div>