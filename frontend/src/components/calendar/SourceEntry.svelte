<script lang="ts">
  import CollapseToggle from "../interactive/CollapseToggle.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMetadata } from "$lib/client/metadata";
  import { getRepository } from "$lib/client/repository";

  import { getContext, untrack } from "svelte";

  interface Props {
    source: SourceModel;
  }

  let {
    source = $bindable()
  }: Props = $props();

  let hasErrored = $state(false);
  getMetadata().faultySources.subscribe((faulty) => {
    if (!source || !source.id) return;
    hasErrored = faulty.has(source.id);
  });

  let isLoading = $state(false);
  getMetadata().loadingSources.subscribe((loading) => {
    if (!source || !source.id) return;
    isLoading = loading.has(source.id) as boolean;
  });

  let hasCals = $state(false);
  getRepository().calendars.subscribe(async (cals) => {
    hasCals = false;
    if (!source) return;
    for (const cal of cals) {
      if (cal.source === source.id) {
        hasCals = true;
        break;
      }
    }
  })

  let showModal: ((source: SourceModel) => Promise<SourceModel>) = getContext("showSourceModal");
  function showModalInternal() {
    showModal(source).then(newSource => source = newSource).catch(NoOp);
  }

  let sourceCollapsed = $state(source ? getMetadata().collapsedSources.has(source.id) : false);
  getMetadata().collapsedSources.subscribe((collapsed) => {
    if (!source || !source.id) return;
    sourceCollapsed = collapsed.has(source.id);
  });
  $effect(() => {
    if (source && source.id) getMetadata().setSourceCollapse(source.id, sourceCollapsed);
  });
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/colors.scss";

  div {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
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
    gap: dimensions.$gapTiny;
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
  <button onclick={showModalInternal} use:focusIndicator={{ type: "underline" }}>
    {source.name}
  </button>
  <span>
    {#if isLoading}
      <Spinner/>
    {/if}
    {#if hasCals}
      <CollapseToggle bind:collapsed={sourceCollapsed}/>
    {/if}
    {#if hasErrored}
      <Tooltip error={true}>An error occurred trying to retrieve calendars from this source.</Tooltip>
    {/if}
  </span>
  <!--
  <IconButton callback={showModal}>
    <PencilIcon size={16}/>
    <BoltIcon size={16}/>
    <CogIcon size={16}/>
  </IconButton>
  -->
</div>
