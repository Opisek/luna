<script lang="ts">
  import Modal from "./Modal.svelte";

  import { NoOp } from "../../lib/client/placeholders";
  import { ColorKeys } from "../../types/colors";
  import { Check, X } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import RadioInput from "../forms/RadioInput.svelte";

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

<Modal title={`${editing ? "Editing" : "Deleting"} recurring event`} bind:showModal={showModalInternal} bind:success bind:failure>
  <RadioInput
    name="recurrence_affect"
    bind:value={chosen}
    options={[
      { name: "This event", value: "this" },
      { name: "This and future events", value: "thisandfuture" },
      { name: "All events", value: "all" }
    ]}
  />
  {#snippet buttons()}
    <IconButton onClick={() => success(chosen || "this")} color={ColorKeys.Success} type="submit" alt="Confirm" canRenderAsButton={true} enabled={chosen != null}>
      <Check/>
    </IconButton>
    <IconButton onClick={failure} color={ColorKeys.Danger} alt="Cancel" canRenderAsButton={true}><X/></IconButton>
  {/snippet}
</Modal>