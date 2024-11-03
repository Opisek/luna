<script lang="ts">
  import { focusIndicator } from "../../lib/client/decoration";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";
  import Label from "./Label.svelte";

  export let value: string = "";
  export let placeholder: string;
  export let name: string;

  export let editable: boolean = true;

  export let multiline: boolean = false;

  export let password: boolean = false;
  let passwordVisible: boolean = false;

  export let label: boolean = true;

  export let onChange: (value: string) => any = () => {};
  export let onInput: (value: string) => any = () => {};
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
  @import "../../styles/animations.scss";
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  div.wrapper {
    display: flex;
    flex-direction: column;
    gap: $gap;
    position: relative;
    border-radius: calc($borderRadius + 0.1em);
    overflow: hidden;
  }

  input, textarea {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
  }

  div.editable > input, div.editable > textarea {
    background: $backgroundSecondary;
  }

  textarea {
    min-height: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    //padding: 0;
    //margin: 0;
  }

  div.visibility {
    position: absolute;
    //height: 100%;
    top: 50%;
    transform: translateY(-50%);
    right: $gapSmaller;
    color: $foregroundFaded;
  }
</style>

{#if label}
  <Label name={name}>{placeholder}</Label>
{/if}
<div
  class="wrapper"
  class:editable={editable} 
  tabindex="-1"
  use:focusIndicator
>
  {#if multiline}
      <textarea
        bind:value={value}
        on:change={() => onChange(value)}
        on:input={() => onInput(value)}
        on:focusout={() => onChange(value)}
        on:focusin={() => onFocus()}
        name={name}
        placeholder={placeholder}
        disabled={!editable}
        rows=6
        tabindex={editable ? 0 : -1}
      />
  {:else if password && !passwordVisible}
    <input
      bind:value={value}
      on:change={() => onChange(value)}
      on:input={() => onInput(value)}
      on:focusout={() => onChange(value)}
      on:focusin={() => onFocus()}
      name={name}
      placeholder={placeholder}
      disabled={!editable}
      class:editable={editable}
      tabindex={editable ? 0 : -1}
      type="password"
    />
  {:else}
    <input
      bind:value={value}
      on:change={() => onChange(value)}
      on:input={() => onInput(value)}
      on:focusout={() => onChange(value)}
      on:focusin={() => onFocus()}
      name={name}
      placeholder={placeholder}
      disabled={!editable}
      class:editable={editable}
      tabindex={editable ? 0 : -1}
      type="text"
    />
  {/if}
  {#if password && editable}
  <div class="visibility">
    <VisibilityToggle bind:visible={passwordVisible} momentary={true} />
  </div>
  {/if}
</div>