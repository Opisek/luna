<script lang="ts">
  import { browser } from "$app/environment";
  import { page } from "$app/state";

  import Box from "../components/layout/Box.svelte";
  import Button from "../components/interactive/Button.svelte";
  import Horizontal from "../components/layout/Horizontal.svelte";
  import SimplePage from "../components/layout/SimplePage.svelte";
  import Title from "../components/layout/Title.svelte";
</script>

<style lang="scss">

</style>

<SimplePage>
  <Box>
    <Title>
      Error {page.status}
    </Title>

    {#if page.error}
      {page.error.message}
    {:else}
      An error occurred.
    {/if}

    <Horizontal position="right">
      <Button href="/">Home</Button>
      {#if browser && window.history.length > 1}
        <Button onClick={() => history.back()}>Go Back</Button>
      {/if}
    </Horizontal>
  </Box>
</SimplePage>