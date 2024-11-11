<script lang="ts">
  export let name: string;
  export let ownPositioning: boolean = true;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  label {
    color: $foregroundFaded;
    font-size: $fontSizeSmall;
    cursor: text;
  }

  .ownPositioning {
    margin-bottom: -2 * $gapSmall;
    margin-left: calc($gapSmall * ($fontSize / $fontSizeSmall));
  }
</style>

<label for={name} tabindex="-1" class:ownPositioning={ownPositioning}>
  <slot/>
</label>