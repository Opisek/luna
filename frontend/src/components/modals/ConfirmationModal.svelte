<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";

  interface Props {
    showModal: () => Promise<void>;
    children?: import('svelte').Snippet;
  }

  let success: (result: void) => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);

  let {
    showModal = $bindable(),
    children
  }: Props = $props();
</script>

<Modal title="Confirmation" bind:showModal bind:success bind:failure>
  {@render children?.()}
  {#snippet buttons()}
    <IconButton onClick={success} color={ColorKeys.Success} type="submit" alt="Confirm" canRenderAsButton={true}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>