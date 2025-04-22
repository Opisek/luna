<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidPassword, valid } from "../../lib/client/validation";
  import { deepCopy } from "../../lib/common/misc";

  interface Props {
    prompt: () => Promise<string>;
  }

  let {
    prompt: showModal = $bindable(),
  }: Props = $props();

  let password = $state("");
  let passwordValidity = $state(valid);

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  let resolvePromise = $state<(password: string) => Promise<any>>(Promise.reject);
  let rejectPromise = $state(NoOp);

  showModal = async () => {
    password = "";
    passwordValidity = valid;

    showModalInternal();

    return new Promise((resolve, reject) => {
      resolvePromise = (async (password) => {
        resolve(password);
        resolvePromise = Promise.reject;
      })

      rejectPromise = (() => {
        reject("No password provided");
        rejectPromise = NoOp;
      })
    })
  };

  function confirm() {
    resolvePromise(password);
    hideModalInternal();
  }

  function cancel() {
    rejectPromise();
    hideModalInternal();
  }
</script>

<Modal title="Authentication" bind:showModal={showModalInternal} bind:hideModal={hideModalInternal} onModalHide={rejectPromise}>
  To do this action, you must authenticate yourself with your password.

  <TextInput
    name="password"
    placeholder="Password"
    password={true}
    bind:value={password}
    validation={isValidPassword}
    bind:validity={passwordValidity}
  />

  {#snippet buttons()}
      <Button onClick={confirm} color={ColorKeys.Success} type="submit" enabled={password != "" && passwordValidity.valid}>
        Confirm
      </Button>
      <Button onClick={cancel} color={ColorKeys.Danger}>Cancel</Button>
  {/snippet}
</Modal>