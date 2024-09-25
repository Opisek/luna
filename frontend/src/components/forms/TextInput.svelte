<script lang="ts">
  import Label from "./Label.svelte";

  export let value: string;
  export let placeholder: string;
  export let name: string;

  export let editable: boolean = true;

  export let multiline: boolean = false;
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  input, textarea {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
    cursor: revert;
  }

  input.editable, textarea.editable {
    background: $backgroundSecondary;
    cursor: text;
  }

  textarea {
    height: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  textarea.editable {
    height: 5em; // TODO: automatic resizing
  }
</style>

<Label name={name}>{placeholder}</Label>
{#if multiline}
  <textarea
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
  />
{:else}
  <input
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
  />
{/if}