<script lang="ts">
  import { X } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";

  export let onClick: () => void;
</script>

<IconButton bind:click={onClick}>
  <X/>
</IconButton>