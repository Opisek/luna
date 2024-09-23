<script lang="ts">
  import IconButton from "../interactive/IconButton.svelte";

  export const showModal = () => (visible = true);

  let visible = false;
  let dialog: HTMLDialogElement;

  $: if (dialog && visible) dialog.showModal();
</script>

<style lang="scss">

</style>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
  bind:this={dialog}
  on:close={() => (visible = false)}
  on:click|self={() => dialog.close()}
>
  <!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click|stopPropagation>
    <!-- TODO: add close button and title -->
		<slot />
	</div>
</dialog>