<script lang="ts">
  import Checkbox from "../interactive/Checkbox.svelte";

  export let value: boolean;
  export let description: string;
  export let name: string;

  export let editable: boolean = true;

  export let onChange: (value: boolean) => any = () => {};
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
  }
</style>

<div>
  <Checkbox
    bind:value
    name={name}
    onChange={onChange}
    enabled={editable}
  />
  <label for={name}>
    {description}
  </label>
</div>