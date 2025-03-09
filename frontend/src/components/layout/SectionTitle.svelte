<script lang="ts">
  interface Props {
    title: string;
  }

  let { title }: Props = $props();
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  span {
    margin: 0;
    padding: $gapSmall;
    color: $foregroundDim;
  }
</style>

<span>
  {title}
</span>