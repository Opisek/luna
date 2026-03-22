<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
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

  let awaitingConfirm = $state(false);
  function confirm() {
    awaitingConfirm = true;
    confirmCallback().catch(err => {
      queueNotification(ColorKeys.Danger, err)
    }).finally(() => {
      awaitingConfirm = false;
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
      {#if awaitingConfirm}
        <Loader/>
      {:else}
        <Check/>
      {/if}
    </IconButton>
    <IconButton onClick={cancel} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>