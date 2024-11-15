<script lang="ts">
  import Label from "./Label.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { alwaysValid, valid } from "$lib/client/validation";
  import { focusIndicator } from "$lib/client/decoration";
  import { NoOp } from "../../lib/client/placeholders";

  let passwordVisible: boolean = $state(false);

  let wrapper: HTMLDivElement | null = $state(null);

  let lastValidationFunction = $state(alwaysValid); // TODO: check if still needed in svelte 5
  interface Props {
    value?: string | undefined;
    placeholder: string;
    name: string;
    editable?: boolean;
    multiline?: boolean;
    password?: boolean;
    label?: boolean;
    onChange?: (value: string) => any;
    onInput?: (value: string) => any;
    onFocus?: () => any;
    validation?: InputValidation;
    validity?: any;
  }

  let {
    value = $bindable(),
    placeholder,
    name,
    editable = true,
    multiline = false,
    password = false,
    label = true,
    onChange = NoOp,
    onInput = NoOp,
    onFocus = NoOp,
    validation = alwaysValid,
    validity = $bindable(value ? validation(value) : valid)
  }: Props = $props();


  // If the value is set programmatically, update the validity.
  // For example when opening a new form
  let lastValue: string | null = $state(null); // TODO: check if still needed in svelte 5
  $effect(() => {
    ((value) => {
      if (!value || value === lastValue) return; // prevents some infinite loop that i don't understand, might be a svelte bug
      lastValue = value;
      if (wrapper != null && (document.activeElement === wrapper || wrapper.contains(document.activeElement))) return;
      validity = value ? validation(value) : valid;
    })(value);
  });

  // This determines whether input has errored due to empty value.
  // This is still considered an error, but we don't want to display it.
  let empty = $state(value === "");

  // Once the user has finished typing, update the validity.
  function internalOnChange() {
    if (!value) return;
    validity = validation(value);
    empty = value === "";
    onChange(value);
  }

  // Immediately tell the user if the input becomes valid,
  // but not if it becomes invalid, as they are not done typing yet.
  function internalOnInput() {
    if (!value) return;
    const res = validation(value);
    if (res.valid) validity = res;
    onInput(value);
  }

  // If the validation function changes, like for the repeat password field,
  // rerun the validation function.
  $effect(() => {
    ((_) => {
      if (validation === lastValidationFunction) return;
      lastValidationFunction = validation;
      internalOnChange();
    })(validation);
  });

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
  div.noneditable {
    --barFocusIndicatorColor: transparent;
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

  span.label {
    margin-bottom: -2 * $gapSmall;
    padding: 0 calc($gapSmall * ($fontSize / $fontSizeSmall));
    display: flex;
    justify-content: space-between;
  }

  span.errorMessage {
    color: $backgroundFailure;
    font-size: $fontSizeSmall;
  }
</style>

{#if label || (!validity?.valid && !empty)}
  <span class="label">
    {#if label}
      <Label name={name} ownPositioning={false}>{placeholder}</Label>
    {/if}
    {#if !validity?.valid && !empty}
      <span class="errorMessage">
        {validity.message}
      </span>
    {/if}
  </span>
{/if}
<div
  class="wrapper"
  class:editable={editable} 
  class:noneditable={!editable} 
  tabindex="-1"
  use:focusIndicator
  class:error={!validity.valid && !empty}
  bind:this={wrapper}
>
  {#if multiline}
      <textarea
        bind:value={value}
        onchange={internalOnChange}
        oninput={internalOnInput}
        onfocusout={internalOnChange}
        onfocusin={onFocus}
        name={name}
        placeholder={placeholder}
        disabled={!editable}
        rows=6
        tabindex={editable ? 0 : -1}
></textarea>
  {:else if password && !passwordVisible}
    <input
      bind:value={value}
      onchange={internalOnChange}
      oninput={internalOnInput}
      onfocusout={internalOnChange}
      onfocusin={onFocus}
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
      onchange={internalOnChange}
      oninput={internalOnInput}
      onfocusout={internalOnChange}
      onfocusin={onFocus}
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

<!-- TODO: snippets and svelte:element in conjuction with {...otherProps} to reduce amount of rewritten code -->