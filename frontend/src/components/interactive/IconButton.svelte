<script lang="ts">
  export let callback: () => void;
  export let visible: boolean = true;
</script>

<style lang="scss">
  @import "../../styles/dimensions.scss";

  button {
    all: unset;
    border-radius: 50%;
    display: flex;
    align-items: center;
    padding: $paddingTiny;
    cursor: pointer;
  }

  button.hidden {
    visibility: hidden;
  }
</style>

<button on:click={callback} class:hidden={!visible}>
  <slot/>
</button>