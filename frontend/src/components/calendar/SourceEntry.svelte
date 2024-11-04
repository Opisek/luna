<script lang="ts">
  import { focusIndicator } from "$lib/client/decoration";
  import { calendars, faultySources, fetchSourceCalendars } from "$lib/client/repository";
  import { collapsedSources, setSourceCollapse } from "../../lib/client/localStorage";
  import CollapseToggle from "../interactive/CollapseToggle.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";
  import SourceModal from "../modals/SourceModal.svelte";

  export let source: SourceModel;

  let hasErrored = false;
  faultySources.subscribe((faulty) => {
    hasErrored = faulty.has(source.id);
  });

  let hasCals = false;
  calendars.subscribe(async () => {
    hasCals = (await fetchSourceCalendars(source.id)).length > 0;
  })

  let showModal: () => any;

  $: if (source && source.id) setSourceCollapse(source.id, source.collapsed);
  collapsedSources.subscribe((collapsed) => {
    source.collapsed = collapsed.has(source.id);
  });
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/colors.scss";

  div {
    color: $foregroundFaded;
    height: 1.25em;
    display: flex;
    justify-content: space-between;
    align-items: center;
    align-content: center;
  }

  span {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    gap: $gapTiny;
    align-items: center;
  }

  button {
    all: unset;
    cursor: pointer;
    display: inline;
    width: max-content;
    position: relative;
  }

  //span :global(button) {
  //  opacity: 0;
  //  transition: all $animationSpeed $cubic;
  //}

  //span:hover :global(button) {
  //  opacity: 1;
  //}
</style>

<div>
  <button on:click={showModal} use:focusIndicator={{ type: "underline", ignoreParent: true }}>
    {source.name}
  </button>
  <span>
    {#if hasCals}
      <CollapseToggle bind:collapsed={source.collapsed}/>
    {/if}
    {#if hasErrored}
      <Tooltip msg="An error occurred trying to retrieve calendars from this source." error={true}/>
    {/if}
  </span>
  <!--
  <IconButton callback={showModal}>
    <PencilIcon size={16}/>
    <BoltIcon size={16}/>
    <CogIcon size={16}/>
  </IconButton>
  -->

  <SourceModal
    bind:showModal={showModal}
    source={source}
  />
</div>
