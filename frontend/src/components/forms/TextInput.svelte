<script lang="ts">
    import Label from "./Label.svelte";

  export let value: string;
  export let placeholder: string;
  export let name: string;

  export let multiline: boolean = false;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  input, textarea {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
    background: $backgroundSecondary;
  }

  textarea {
    height: 5em; // TODO: automatic resizing
    white-space: pre-wrap;
    word-wrap: break-word;
  }
</style>

<Label name={name}>{placeholder}</Label>
{#if multiline}
  <textarea bind:value={value} placeholder={placeholder} />
{:else}
  <input name={name} type="text" bind:value={value} placeholder={placeholder} class="input" />
{/if}