<script lang="ts">
  import { X } from "lucide-svelte";
  import Title from "../layout/Title.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import CloseButton from "../interactive/CloseButton.svelte";

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
		animation: zoom 0.3s $cubic forwards;
	}
  dialog[open]::backdrop {
    animation: fade 10s $cubic forwards;
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
	</div>
</dialog>