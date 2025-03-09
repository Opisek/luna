<script lang="ts">
  interface Props {
    title: string;
  }

  let { title }: Props = $props();
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  span {
    margin: 0;
    padding: dimensions.$gapSmall;
    color: colors.$foregroundDim;
  }
</style>

<span>
  {title}
</span>