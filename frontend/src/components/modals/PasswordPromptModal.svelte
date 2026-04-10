<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidPassword, valid } from "../../lib/client/validation";

  interface Props {
    prompt: () => Promise<string>;
  }

  let {
    prompt: showModal = $bindable(),
  }: Props = $props();

  let password = $state("");
  let passwordValidity = $state(valid);

  let showModalInternal: () => Promise<string> = $state(Promise.reject);
  let success: (password: string) => void = $state(NoOp);
  let failure: () => void = $state(NoOp);

  showModal = async () => {
    password = "";
    passwordValidity = valid;
    return showModalInternal();
  };

  function confirm() {
    success($state.snapshot(password));
    password = "";
  }

  function cancel() {
    failure();
    password = "";
  }
</script>

<Modal
  title="Authentication"
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
>
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