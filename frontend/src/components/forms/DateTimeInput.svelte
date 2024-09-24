<script lang="ts">
  import Label from "./Label.svelte";

  export let value: Date;
  export let placeholder: string;
  export let name: string;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    flex-direction: row;
    margin-top: -$gapSmall;
    gap: $gapSmall;
    padding: $gapSmall;
  }
</style>

<Label name={name}>{placeholder}</Label>
<div>
  <span>{value.toLocaleDateString()}</span>
  <span>{value.toLocaleTimeString([], {hour: "2-digit", minute: "2-digit"})}</span>
</div>