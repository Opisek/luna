<script lang="ts">
  import Label from "./Label.svelte";
  import VisibilityToggle from "../interactive/VisibilityToggle.svelte";

  import { alwaysValid, valid } from "$lib/client/validation";
  import { focusIndicator } from "$lib/client/decoration";
  import { NoOp } from "../../lib/client/placeholders";
  import IconButton from "../interactive/IconButton.svelte";
  import { Copy } from "lucide-svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";

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
    mono?: boolean;
    displayCopyButton?: boolean;
    label?: boolean;
    onChange?: (value: string, event: Event | null) => any;
    onInput?: (value: string, event: Event | null) => any;
    onFocus?: () => any;
    validation?: InputValidation;
    formatting?: (value: string, event: Event | null) => string;
    validity?: Validity;
  }

  let {
    value = $bindable(),
    placeholder,
    name,
    editable = true,
    multiline = false,
    password = false,
    mono = false,
    displayCopyButton = false,
    label = true,
    onChange = NoOp,
    onInput = NoOp,
    onFocus = NoOp,
    validation = alwaysValid,
    formatting = (value, _) => value,
    validity = $bindable(valid)
  }: Props = $props();

  let element: HTMLInputElement | HTMLTextAreaElement | null = $state(null);

  // If the value is set programmatically, update the validity.
  // For example when opening a new form
  let lastValue: string | null = $state(null); // TODO: check if still needed in svelte 5
  $effect(() => {
    (async (value) => {
      if (!value || value === lastValue) return; // prevents some infinite loop that i don't understand, might be a svelte bug
      lastValue = value;
      if (wrapper != null && (document.activeElement === wrapper || wrapper.contains(document.activeElement))) return;
      validity = value ? await validation(value) : valid;
    })(value);
  });

  // This determines whether input has errored due to empty value.
  // This is still considered an error, but we don't want to display it.
  let empty = $state(value === "");

  // Once the user has finished typing, update the validity.
  async function internalOnChange(event: Event | null) {
    if (!value) return;
    value = formatting(value, event);
    validity = await validation(value);
    empty = value === "";
    onChange(value, event);
  }

  // Immediately tell the user if the input becomes valid,
  // but not if it becomes invalid, as they are not done typing yet.
  async function internalOnInput(event: Event | null) {
    if (!value) return;
    value = formatting(value, event);
    const res = await validation(value);
    if (res.valid) validity = res;
    onInput(value, event);
  }

  // If the validation function changes, like for the repeat password field,
  // rerun the validation function.
  $effect(() => {
    ((_) => {
      if (validation === lastValidationFunction) return;
      lastValidationFunction = validation;
      internalOnChange(null);
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

  function copy() {
    navigator.clipboard.writeText(value || "").then(() => {
      queueNotification(ColorKeys.Success, "Copied to clipboard");
    });
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.wrapper {
    display: flex;
    flex-direction: row;
    gap: dimensions.$gapTiny;
    align-items: center;
    justify-content: center;
    position: relative;

    border-radius: calc(dimensions.$borderRadius + 0.1em);
    padding: 0 dimensions.$gapSmall;
    width: 100%;
    overflow: hidden;

    color: color-mix(in srgb, colors.$foregroundSecondary 50%, transparent);
  }

  input, textarea {
    all: unset;
    flex-grow: 1;
    margin: dimensions.$gapSmall 0;
    padding: 0;
  }

  div.wrapper.editable {
    background: colors.$backgroundSecondary;
  }
  div.wrapper.editable > input, div.wrapper.editable > textarea {
    color: colors.$foregroundSecondary;
  }
  div.noneditable {
    --barFocusIndicatorColor: transparent;
  }

  textarea {
    min-height: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
  }

  div.wrapper.mono > input, div.wrapper.mono > textarea {
    font-family: text.$fontFamilyTime;
  }

  span.label {
    font-size: text.$fontSizeSmall;
    margin-bottom: -(dimensions.$gapMiddle);
    display: flex;
    justify-content: space-between;
  }

  span.errorMessage {
    color: colors.$backgroundFailure;
    font-size: text.$fontSizeSmall;
  }
</style>

{#if label || (!validity?.valid && !empty)}
  <!-- TODO: use the Label component instead -->
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
  class:mono={mono}
  tabindex="-1"
  use:focusIndicator
  class:error={!validity.valid && !empty}
  bind:this={wrapper}
>
  {#if multiline}
    <textarea
      bind:this={element}
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
      bind:this={element}
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
      bind:this={element}
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
  {#if password &&editable}
    <VisibilityToggle bind:visible={passwordVisible} momentary={true} />
  {/if}
  {#if displayCopyButton}
    {@render copyButton()}
  {/if}
</div>

{#snippet copyButton()}
  <IconButton>
    <Copy size={16} onclick={copy}/>
  </IconButton>
{/snippet}

<!-- TODO: snippets and svelte:element in conjuction with {...otherProps} to reduce amount of rewritten code -->