<script lang="ts">
  import Loader from "../decoration/Loader.svelte";
  import Button from "../interactive/Button.svelte";
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";

  export let title: string;

  let loading = false;
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
    width: 100%;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: $gap;
    box-shadow: $boxShadow;
  }
</style>

<form method="POST">
  <Title>{title}</Title>
  <slot/>
  <Horizontal position="right">
    <Button type="submit" color="success" onClick={ () => (loading = true) }>
      {#if loading}
        <Loader/>
      {:else}
        Save
      {/if}
    </Button>
  </Horizontal>
</form>