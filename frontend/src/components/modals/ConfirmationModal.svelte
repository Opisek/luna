<script lang="ts">
  import Modal from "./Modal.svelte";

  import { queueNotification } from "$lib/client/notifications";
  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";

  interface Props {
    confirmCallback: () => Promise<void>;
    cancelCallback?: () => void;
    showModal: () => any;
    hideModal?: () => any;
    children?: import('svelte').Snippet;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
    confirmCallback,
    cancelCallback = () => {},
    children
  }: Props = $props();

  async function confirm() {
    await confirmCallback().catch(err => {
      queueNotification(ColorKeys.Danger, err)
    }).finally(() => {
      hideModal()
    });
  }

  function cancel() {
    hideModal()
    cancelCallback()
  }
</script>

<Modal title="Confirmation" bind:showModal={showModal} bind:hideModal={hideModal}>
  {@render children?.()}
  {#snippet buttons()}
    <IconButton onClick={confirm} color={ColorKeys.Success} type="submit" alt="Confirm" canRenderAsButton={true}>
      <Check/>
    </IconButton>
    <IconButton onClick={cancel} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>