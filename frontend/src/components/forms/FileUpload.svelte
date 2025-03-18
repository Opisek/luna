<script lang="ts">
  import Label from "./Label.svelte";

  import { alwaysValidFile, valid } from "$lib/client/validation";
  import { focusIndicator } from "$lib/client/decoration";
  import IconButton from "../interactive/IconButton.svelte";
  import { Download, FileEdit, Upload, X } from "lucide-svelte";
  import { downloadFileToClient } from "../../lib/client/net";

  let wrapper: HTMLDivElement | null = $state(null);
  let fileInput: HTMLInputElement | null = $state(null);

  interface Props {
    files: FileList | null;
    fileId?: string;
    placeholder: string;
    name: string;
    editable?: boolean;
    validation?: FileValidation;
    validity?: any;
  }

  let {
    files = $bindable(null),
    fileId = $bindable(""),
    placeholder,
    name,
    editable = true,
    validation = alwaysValidFile,
    validity = $bindable(valid)
  }: Props = $props();

  function select() {
    if (!editable) return;
    fileInput?.click();
  }

  function clear() {
    if (!editable) return;
    if (fileInput) fileInput.value = "";
    files = null;
  }

  function download() {
    if (fileId === "") downloadFileToClient(files)
    else downloadFileToClient(fileId);
  }

  // If the value is set programmatically, update the validity.
  // For example when opening a new form
  let lastValue: FileList | null = $state(null); // TODO: check if still needed in svelte 5
  $effect(() => {
    (async (value) => {
      if (!value || value === lastValue) return; // prevents some infinite loop that i don't understand, might be a svelte bug
      lastValue = value;
      if (wrapper != null && (document.activeElement === wrapper || wrapper.contains(document.activeElement))) return;
      validity = value && fileId == "" ? await validation(value) : valid;
    })(files);
  });

  // This determines whether input has errored due to empty value.
  // This is still considered an error, but we don't want to display it.
  let empty = $derived(files === null);

  // Update validity when the file changes
  async function internalOnChange() {
    if (files) {
      fileId = "";
    }
    validity = files ? await validation(files) : valid;
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/text.scss";

  div.wrapper {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapLarge;
    position: relative;
    border-radius: calc(dimensions.$borderRadius + 0.1em);
    overflow: hidden;
  }

  input {
    all: unset;
    padding: dimensions.$gapSmall;
    border-radius: dimensions.$borderRadius;
  }

  //input::file-selector-button {
  //  position: absolute;
  //  top: 50%;
  //  transform: translateY(-50%);
  //  right: calc(dimensions.$gapSmaller - 0.5 * #{dimensions.$gapSmall});
  //  outline: none;
  //  border: none;
  //  color: transparent;
  //  overflow: hidden;
  //  height: calc(100% - #{dimensions.$gapSmall});
  //  width: auto;
  //  aspect-ratio: 1 / 1;
  //  cursor: pointer;
  //  z-index: 5;
  //  background: transparent;
  //}

  input::file-selector-button {
    display: none;
    visibility: none;
    pointer-events: none;
  }

  input.empty {
    color: colors.$foregroundDim;
  }

  div.editable > input {
    background: colors.$backgroundSecondary;
    cursor: pointer;
  }
  div.noneditable {
    --barFocusIndicatorColor: transparent;
  }

  div.buttons {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    right: calc(1.5 * dimensions.$gapSmaller);
    color: colors.$foregroundDim;
    display: flex;
    flex-direction: row;
  }
  div.editable > div.buttons {
    background-color: colors.$backgroundSecondary;
  }

  span.label {
    font-size: text.$fontSizeSmall;
    margin-bottom: -(dimensions.$gapMiddle);
    padding-left: calc(dimensions.$gapSmall * (text.$fontSize / text.$fontSizeSmall));
  }

  span.errorMessage {
    color: colors.$backgroundFailure;
    font-size: text.$fontSizeSmall;
  }
</style>

<!-- TODO: use the Label component instead -->
<span class="label">
    <Label name={name} ownPositioning={false}>{placeholder}</Label>
{#if !validity?.valid && !empty}
    <span class="errorMessage">
    {validity.message}
    </span>
{/if}
</span>

<div
  class="wrapper"
  class:editable={editable} 
  class:noneditable={!editable} 
  tabindex="-1"
  use:focusIndicator
  class:error={!validity.valid && !empty}
  bind:this={wrapper}
>
<input
    type="file"
    accept=".ical,.ics,.ifb,.icalendar"
    onchange={internalOnChange}
    name={name}
    disabled={!editable}
    class:editable={editable}
    class:empty={files === null}
    tabindex={editable ? 0 : -1}
    bind:this={fileInput}
    bind:files
/>
  {#if editable}
    {#if empty}
      <div class="buttons">
        <IconButton click={select}>
            <!-- Upload, FileUp, MonitorUp, CloudUpload, HardDriveUpload -->
            <Upload size={16}/>
        </IconButton>
      </div>
    {:else}
      <div class="buttons">
        <IconButton click={download}>
          <Download size={16}/>
        </IconButton>
        <IconButton click={clear}>
            <X size={16}/>
        </IconButton>
      </div>
    {/if}
  {:else}
    <div class="buttons">
      <IconButton click={download}>
        <Download size={16}/>
      </IconButton>
    </div>
  {/if}
</div>

<!-- TODO: snippets and svelte:element in conjuction with {...otherProps} to reduce amount of rewritten code -->