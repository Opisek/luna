<script lang="ts">
  export let name: string;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  label {
    margin-bottom: -2 * $gapSmall;
    margin-left: $gapSmall * ($fontSize / $fontSizeSmall);
    color: $foregroundFaded;
    font-size: $fontSizeSmall;
  }
</style>

<label for={name}>
  <slot/>
</label>