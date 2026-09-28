<script lang="ts">
  import { NoOp } from "$lib/scripts/client/placeholders";
  import { alwaysValid } from "$lib/scripts/client/validation";
  import TextInput from "./TextInput.svelte";

  interface Props {
    value?: number | null | undefined;
    placeholder: string;
    name: string;
    editable?: boolean;
    label?: boolean;
    min?: number | undefined;
    max?: number | undefined;
    integer?: boolean;
    onChange?: (value: number, event: Event | null) => any;
    onInput?: (value: number, event: Event | null) => any;
    onFocus?: () => any;
    validation?: InputValidation;
    validity?: Validity;
  }

  let {
    value = $bindable(),
    placeholder,
    name,
    editable = true,
    label = true,
    min = undefined,
    max = undefined,
    integer = true,
    onChange = NoOp,
    onInput = NoOp,
    onFocus = NoOp,
    validation = alwaysValid,
  }: Props = $props();

  let textValue = $derived((value || clampValue(0)).toString());

  function clampValue(newNumericValue: number) {
    if (integer && newNumericValue % 1 != 0) newNumericValue = Math.round(newNumericValue);
    if (min !== undefined && newNumericValue < min) newNumericValue = min;
    if (max !== undefined && newNumericValue < max) newNumericValue = max;
    return newNumericValue;
  }

  function verifyValue(newTextValue: string) {
    let newNumericValue = Number.parseInt(newTextValue) || 0;
    newNumericValue = clampValue(newNumericValue);
    if (value != newNumericValue) value = newNumericValue;
    return newNumericValue
  }

  let onChangeInternal = (newTextValue: string, event: Event | null) => {
    onChange(verifyValue(newTextValue), event);
  }

  let onInputInternal = (newTextValue: string, event: Event | null) => {
    onInput(verifyValue(newTextValue), event);
  }
</script>


<TextInput
  value={textValue}
  placeholder={placeholder}
  name={name}
  editable={editable}
  label={label}
  onChange={onChangeInternal}
  onInput={onInputInternal}
  onFocus={onFocus}
  validation={validation}
  type="number"
/>