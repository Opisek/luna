<script lang="ts">
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";
  import Label from "./Label.svelte";

  export let value: string;
  export let placeholder: string;
  export let name: string;

  export let editable: boolean = true;

  export let multiline: boolean = false;

  export let password: boolean = false;
  let passwordVisible: boolean = false;
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

  div.visibility {
    text-align: right;
    position: relative;
    top: calc(-1.25em - 2 * $gapSmall - 0.5 * $fontSize);
    margin-bottom: calc(-1.25em - $gapSmall - 0.5 * $fontSize);
    right: calc(-100% + 1.25em + $gapSmall);
    width: fit-content;
    display: flex;
    justify-content: flex-end;
    color: $foregroundFaded;
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
{:else if password && !passwordVisible}
  <input
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
    type="password"
  />
{:else}
  <input
    bind:value={value}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
    type="text"
  />
{/if}
{#if password}
<div class="visibility">
  <VisibilityToggle bind:visible={passwordVisible} momentary={true} />
</div>
{/if}