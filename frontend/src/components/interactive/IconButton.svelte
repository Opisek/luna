<script lang="ts">
  export let callback: () => void;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  button {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: $paddingTiny;
  }
</style>

<button on:click={callback}>
  <slot/>
</button>