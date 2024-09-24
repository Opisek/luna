<script lang="ts">
  import { X } from "lucide-svelte";

  import RawButton from "./RawButton.svelte";

  export let onClick: () => void;
</script>

<RawButton bind:onClick={onClick}>
  <X/>
</RawButton>