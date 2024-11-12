<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Loader from "../decoration/Loader.svelte";
  import Title from "../layout/Title.svelte";

  interface Props {
    title: string;
    submittable?: boolean;
    children?: import('svelte').Snippet;
  }

  let { title, submittable = true, children }: Props = $props();

  let loading = $state(false);

  function onSubmit(e: SubmitEvent) {
    if (!submittable) {
      e.preventDefault();
      return; // TODO: add some user feedback when a form fails to submit
    }
    loading = true;
  }
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/decoration.scss";

  form {
    border-radius: $borderRadius;
    max-width: 50vw;
    min-width: 30em;
    padding: $gap $gapLarge $gapLarge $gapLarge;
    border-radius: $borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: $gap;
    box-shadow: $boxShadow;
  }
</style>

<form method="POST" onsubmit={onSubmit}>
  <Title>{title}</Title>
  {@render children?.()}
  <Horizontal position="right">
    <Button type="submit" color="success" enabled={submittable}>
      {#if loading}
        <Loader/>
      {:else}
        Submit
      {/if}
    </Button>
  </Horizontal>
</form>