<script lang="ts">
  import Label from "./Label.svelte";

  export let value: string;
  export let placeholder: string;
  export let name: string;

  export let editable: boolean = true;

  export let options: Option[];
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  select {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
    background: transparent;
    cursor: revert;
  }

  select.editable {
    background: $backgroundSecondary;
    cursor: pointer;
    user-select: none;
  }

  // TODO: style the drop-down
</style>

<Label name={name}>{placeholder}</Label>
<select
  bind:value={value}
  name={name}
  placeholder={placeholder}
  disabled={!editable}
  class:editable={editable}
>
  {#each options as option}
    <option value={option.value}>{option.name}</option>
  {/each}
</select>