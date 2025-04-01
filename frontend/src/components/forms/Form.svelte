<script lang="ts">
  import Button from "../interactive/Button.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Loader from "../decoration/Loader.svelte";
  import Title from "../layout/Title.svelte";
  import type { Snippet } from "svelte";

  interface Props {
    title: string;
    submittable?: boolean;
    children?: Snippet;
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
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";

  form {
    border-radius: dimensions.$borderRadius;
    max-width: 50vw;
    min-width: 30em;
    padding: dimensions.$gapLarge dimensions.$gapLarger dimensions.$gapLarger dimensions.$gapLarger;
    border-radius: dimensions.$borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    box-shadow: decorations.$boxShadow;
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