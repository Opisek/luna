<script lang="ts">
  import { X } from "lucide-svelte";

  import IconButton from "./IconButton.svelte";

  interface Props {
    onClick: () => void;
  }

  let {
    onClick = $bindable()
  }: Props = $props();
</script>

<IconButton click={onClick}>
  <X/>
</IconButton>