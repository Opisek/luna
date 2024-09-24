<script lang="ts">
  export let onClick: () => void;
</script>

<style lang="scss">
  button {
    all: unset;
    cursor: pointer;
  }
</style>

<button on:click={onClick}>
  <slot/>
</button>