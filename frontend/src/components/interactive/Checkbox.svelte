<script lang="ts">
  import { addRipple } from "$lib/client/decoration";
  import { CheckIcon } from "lucide-svelte";

  export let value: boolean;
  export let name: string;

  export let onChange: (value: boolean) => any = () => {};

  export let enabled: boolean = true;

  function click(e: MouseEvent) {
    value = !value;
    onChange(value);
    addRipple(e, false);
  }
</script>

<style lang="scss">
  @import "../../styles/colors.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/text.scss";

  button {
    all: unset;

    width: $lineHeightParagraph;
    height: $lineHeightParagraph;

    position: relative;

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;

    border-radius: $borderRadius;
    background-color: $backgroundSecondary;
    overflow: hidden;
  }

  button.check {
    color: $foregroundAccent;
    background-color: $backgroundAccent;
  }

  button.disabled {
    cursor: unset;
  }

  div :global(*) {
    pointer-events: none;
  }

  input {
    all: unset;
    position: absolute;
    left: 0;
    top: 0;
    background-color: red;
  }
</style>

<button class:check={value} class:disabled={!enabled} on:click={click}>
  {#if value}
    <CheckIcon size={20}/>
  {/if}
</button>

<!--
<input
  bind:checked={value}
  on:change={() => onChange(value)}
  type="checkbox"
  disabled={!enabled}
  name={name}
>
-->