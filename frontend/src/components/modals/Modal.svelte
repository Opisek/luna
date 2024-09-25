<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  export let title: string;

  let visible = false;
  let dialog: HTMLDialogElement;

  $: if (dialog && visible) dialog.showModal();

  export const showModal = () => (visible = true);
  export const hideModal = () => dialog.close();
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  dialog {
    border: 0;
    padding: $gap $gapLarge $gapLarge $gapLarge;
    border-radius: $borderRadius;
    max-width: 50vw;
    min-width: 30em;
  }
  dialog::backdrop {
    backdrop-filter: blur($blur);
  }

  dialog[open] {
		animation: zoom $animationSpeed $cubic forwards;
	}
  
  div {
    display: flex;
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: $gap;
  }
</style>

<!-- svelte-ignore a11y-click-events-have-key-events a11y-no-noninteractive-element-interactions -->
<dialog
  bind:this={dialog}
  on:close={() => (visible = false)}
  on:click|self={() => dialog.close()}
  class:closed={visible}
>
  <!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click|stopPropagation>
    <Horizontal>
      <Title>
        {title}
      </Title>
      <CloseButton onClick={hideModal} />
    </Horizontal>
		<slot />
    <Horizontal position="right">
      <slot name="buttons">
        <Button onClick={hideModal} color="accent">Close</Button>
      </slot>
    </Horizontal>
	</div>
</dialog>