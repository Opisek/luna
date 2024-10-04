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

  export let label: boolean = true;

  export let onChange: (value: string) => any = () => {};
  export let onFocus: () => any = () => {};

  // TODO: automatic height 
  // let textArea: HTMLTextAreaElement;
  // let textAreaRows: number = 4;
  // // https://stackoverflow.com/questions/2035910/how-to-get-the-number-of-lines-in-a-textarea
  // function updateTextAreaHeight() {
  //   if (!textArea) return;
  //   textAreaRows = Math.ceil(textArea.scrollHeight / 18.5)
  // }
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";

  input, textarea {
    all: unset;
  }

  input, div.textarea-wrapper {
    padding: $gapSmall;
    border-radius: $borderRadius;
  }

  .editable {
    background: $backgroundSecondary;
  }

  textarea {
    min-height: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    width: 100%;
    padding: 0;
    margin: 0;
  }

  div.visibility {
    text-align: right;
    position: relative;
    // TODO: don't use hard-coded values
    top: calc(-1.31em - 2 * $gapSmall - 0.5 * $fontSize);
    margin-bottom: calc(-1.5em - $gapSmall - 0.5 * $fontSize);
    right: calc(-100% + 1.25em + $gapSmall);
    width: fit-content;
    display: flex;
    justify-content: flex-end;
    color: $foregroundFaded;
  }
</style>

{#if label}
<Label name={name}>{placeholder}</Label>
{/if}
{#if multiline}
  <div
    class="textarea-wrapper"
    class:editable={editable} 
  >
    <textarea
      bind:value={value}
      on:change={() => onChange(value)}
      on:focusout={() => onChange(value)}
      on:focusin={() => onFocus()}
      name={name}
      placeholder={placeholder}
      disabled={!editable}
      rows=6
    />
  </div>
{:else if password && !passwordVisible}
  <input
    bind:value={value}
    on:change={() => onChange(value)}
    on:focusout={() => onChange(value)}
    on:focusin={() => onFocus()}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
    type="password"
  />
{:else}
  <input
    bind:value={value}
    on:change={() => onChange(value)}
    on:focusout={() => onChange(value)}
    on:focusin={() => onFocus()}
    name={name}
    placeholder={placeholder}
    disabled={!editable}
    class:editable={editable}
    type="text"
  />
{/if}
{#if password && editable}
<div class="visibility">
  <VisibilityToggle bind:visible={passwordVisible} momentary={true} />
</div>
{/if}