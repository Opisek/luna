<style lang="scss">
  @import "../../styles/dimensions.scss";

  div {
    display: flex;
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: $gap;
  }
</style>

<div>
  <slot/>
</div>