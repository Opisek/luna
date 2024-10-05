<script lang="ts">
  import Popup from "./Popup.svelte";

  export let show: () => void;
</script>

<style lang="scss">

</style>

<Popup bind:show={show}>
  Select month
</Popup>