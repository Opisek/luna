<script lang="ts">
  import Loader from "../decoration/Loader.svelte";

  interface Props {
    src: string;
    alt: string;
    aspectRatio?: string;
    small?: boolean;
    large?: boolean;
  }

  let {
    src,
    alt,
    aspectRatio = "1/1",
    small = false,
    large = false,
  }: Props = $props();

  let loaded = $state(false);
  let error = $state(false);
  let previousSrc = $state(src);

  $effect(() => {
    if (src !== previousSrc) {
      loaded = false;
      error = false;
      previousSrc = src;
    }
  });
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  div {
    display: flex;
    justify-content: center;
    align-items: center;
    max-width: fit-content;
    height: dimensions.$imageHeight;
    background-color: colors.$backgroundSecondary;
    border-radius: dimensions.$borderRadius;
  }

  div.small {
    height: dimensions.$imageHeightSmall;
  }

  div.large {
    height: dimensions.$imageHeightLarge;
  }

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: dimensions.$borderRadius;
    text-align: center;
    display: flex;
    align-items: center;
  }

  img.loading {
    display: none;
  }
</style>

<div
  style="aspect-ratio: {aspectRatio} !important;"
  class:small={small}
  class:large={large}
>
  {#if !loaded && !error}
    <Loader/>
  {/if}
  <img
    src={src}
    alt={alt}
    onerror={() => error = true}
    onload={() => loaded = true}
    class:loading={!loaded && !error}
  />
</div>