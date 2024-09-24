<script lang="ts">
  import DateTimeInput from "../forms/DateTimeInput.svelte";
import TextInput from "../forms/TextInput.svelte";
    import Vertical from "../layout/Vertical.svelte";
  import Modal from "./Modal.svelte";

  export let event: EventModel;
  export let showModal: () => boolean;

  let title: string;
  $: title = (event && event.id) ? "Edit event" : "Create event";
</script>

<Modal title={title} bind:showModal={showModal}>
  {#if event}
    <Vertical>
      <TextInput bind:value={event.name} name="name" placeholder="Name" />
      <TextInput bind:value={event.desc} name="desc" placeholder="Description" multiline={true}/>
      <DateTimeInput bind:value={event.date.start} name="date_start" placeholder="Start" />
      <DateTimeInput bind:value={event.date.end} name="date_end" placeholder="End" />
    </Vertical>
  {/if}
</Modal>