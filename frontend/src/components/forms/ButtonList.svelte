<script lang="ts" generics="T">
  import { addRipple, focusIndicator } from "../../lib/client/decoration";
  import { EmptyOption } from "../../lib/client/placeholders";

  // This component is used for category lists (currently only in the settings modal)

  interface Props {
    value: T;
    options: Option<T>[][];
  }

  let {
    value = $bindable(),
    options
  }: Props = $props();

  let selected: Option<T> = $derived(options.flat().filter(option => option.value === value)[0] || options[0] || EmptyOption);
</script>

<style lang="scss">
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";

  aside {
    display: flex;
    flex-direction: column;
    overflow: auto;
    padding-right: dimensions.$gapSmall;
    margin-right: -(dimensions.$gapSmall);
  }

  .option {
    border: 0;
    outline: 0;
    font-family: inherit;
    font-size: inherit;
    width: 100%;
    background-color: colors.$backgroundSecondary;
    color: colors.$foregroundSecondary;
    padding: dimensions.$gapSmall;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    gap: dimensions.$gapSmall;
    position: relative;
    overflow: hidden;
    cursor: pointer;
    flex-shrink: 0;
  }

  .first {
    border-top-left-radius: dimensions.$borderRadius;
    border-top-right-radius: dimensions.$borderRadius;
  }

  .last {
    border-bottom-left-radius: dimensions.$borderRadius;
    border-bottom-right-radius: dimensions.$borderRadius;
  }


  :not(.option:last-child).last {
    margin-bottom: dimensions.$gapSmall;
  }

  .selected {
    background-color: colors.$backgroundAccent;
    color: colors.$foregroundAccent;
    --barFocusIndicatorColor: #{colors.$barFocusIndicatorColorAlt};
  }

  .option :global(*) {
    pointer-events: none;
  }

  p {
    margin: 0;
    padding: 0;
    flex-grow: 1;
    text-align: center;
  }
</style>

<aside>
  {#each options as block}
    {#each block as option, i}
      {@const Icon = option.icon}
      <button
        class="option"
        class:first={i === 0}
        class:last={i === block.length - 1}
        class:selected={option.value === value}
        onclick={() => value = option.value}
        onmousedown={addRipple}
        use:focusIndicator
      >
        <Icon size={20}/>
        <p>
          {option.name}
        </p>
      </button>
    {/each}
  {/each}
</aside>