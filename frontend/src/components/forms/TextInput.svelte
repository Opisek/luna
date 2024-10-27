<script lang="ts">
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
    transition: padding $animationSpeedFast linear;
    border-radius: calc($borderRadius + 0.1em);
  }

  div.wrapper.editable {
    background-color: $backgroundAccent;
  }

  input, textarea {
    all: unset;
    padding: $gapSmall;
    border-radius: $borderRadius;
    transition: padding $animationSpeedFast linear, border-radius $animationSpeedFast linear;
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

  div.wrapper.editable:focus-within {
    padding-left: $borderActiveWidth;
    border-top-left-radius: $borderRadius;
    border-bottom-left-radius: $borderRadius;
  }
  div.wrapper.editable:focus-within > input, 
  div.wrapper.editable:focus-within > textarea {
    padding-left: calc($gapSmall - $borderActiveWidth);
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
</style>

{#if label}
  <Label name={name}>{placeholder}</Label>
{/if}
<div
  class="wrapper"
  class:editable={editable} 
  tabindex="-1"
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