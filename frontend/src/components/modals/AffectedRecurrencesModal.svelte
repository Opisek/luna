<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import RadioInput from "../forms/RadioInput.svelte";

  import { t } from "@sveltia/i18n";

  interface Props {
    showModal: (edit: boolean) => Promise<"all" | "thisandfuture" | "this">;
  }

  let success: (result: "all" | "thisandfuture" | "this") => void = $state(NoOp);
  let failure: (reason?: string | Error) => void = $state(NoOp);

  let {
    showModal = $bindable(),
  }: Props = $props();

  let showModalInternal: () => Promise<"all" | "thisandfuture" | "this"> = $state(Promise.reject);
  let editing = $state(false);
  let chosen: "all" | "thisandfuture" | "this" | null = $state("this");

  showModal = async (edit) => {
    editing = edit;
    chosen = "this";
    return showModalInternal();
  }
</script>

<Modal title={t(`event.recurrence.title.${editing ? "edit" : "delete"}`)} bind:showModal={showModalInternal} bind:success bind:failure>
  <RadioInput
    name="recurrence_affect"
    bind:value={chosen}
    options={[
      { name: t("event.recurrence.affect.this"), value: "this" },
      { name: t("event.recurrence.affect.thisandfuture"), value: "thisandfuture" },
      { name: t("event.recurrence.affect.all"), value: "all" }
    ]}
  />
  {#snippet buttons()}
    <IconButton onClick={() => success(chosen || "this")} color={ColorKeys.Success} type="submit" alt={t("button.confirm")} canRenderAsButton={true} enabled={chosen != null}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt={t("button.cancel")} canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>