<script lang="ts">
  import { browser } from "$app/environment";
  import { page } from "$app/state";

  import Box from "../../components/layout/Box.svelte";
  import Button from "../../components/interactive/Button.svelte";
  import Horizontal from "../../components/layout/Horizontal.svelte";
  import SimplePage from "../../components/layout/SimplePage.svelte";
  import Title from "../../components/layout/Title.svelte";

  import { t } from "@sveltia/i18n";
</script>

<style lang="scss">

</style>

<SimplePage>
  <Box>
    <Title>
      {t("error.code", { values: { code: page.status }})}
    </Title>

    {#if page.error}
      {page.error.message}
    {:else}
      {t("error.message.generic")}
    {/if}

    <Horizontal position="right">
      <Button href="/">{t("nav.home")}</Button>
      {#if browser && window.history.length > 1}
        <Button onClick={() => history.back()}>{t("nav.back")}</Button>
      {/if}
    </Horizontal>
  </Box>
</SimplePage>