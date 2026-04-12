<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidPassword, valid } from "../../lib/client/validation";
  import { t } from "@sveltia/i18n";

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
  title={t("confirmation.password.title")}
  bind:showModal={showModalInternal}
  bind:success
  bind:failure
>
  {t("confirmation.password.prompt")}

  <TextInput
    name="password"
    placeholder={t("login.password")}
    password={true}
    bind:value={password}
    validation={isValidPassword}
    bind:validity={passwordValidity}
  />

  {#snippet buttons()}
      <Button onClick={confirm} color={ColorKeys.Success} type="submit" enabled={password != "" && passwordValidity.valid}>{t("button.confirm")}</Button>
      <Button onClick={cancel} color={ColorKeys.Danger}>{t("button.cancel")}</Button>
  {/snippet}
</Modal>