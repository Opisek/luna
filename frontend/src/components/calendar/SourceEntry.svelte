<script lang="ts">
  import CollapseToggle from "../interactive/CollapseToggle.svelte";
  import Spinner from "../decoration/Spinner.svelte";
  import Tooltip from "../interactive/Tooltip.svelte";

  import { NoOp } from "$lib/client/placeholders";
  import { focusIndicator } from "$lib/client/decoration";
  import { getMetadata } from "$lib/client/data/metadata.svelte";
  import { getRepository } from "$lib/client/data/repository.svelte";
  import { draggable } from "$lib/client/reordering.svelte";

  import { getContext } from "svelte";
  import { queueNotification } from "../../lib/client/notifications";
  import { ColorKeys } from "../../types/colors";
  import { t } from "@sveltia/i18n";
  import { extendUniqueElementId, generateUniqueElementId } from "$lib/common/dom";

  interface Props {
    source: SourceModel;
    calendars: CalendarModel[];
    idSeed: string;
  }

  let {
    source = $bindable(),
    calendars,
    idSeed,
  }: Props = $props();
  let sourceEntryId = $derived(generateUniqueElementId(["sourceentry", source.id], idSeed));

  let ariaControls = $derived(calendars.map(x => generateUniqueElementId(["calendarentry", x.id], idSeed)).join(" "));

  const metadata = getMetadata();
  const repository = getRepository();

  let hasErrored = $derived(source && metadata.faultySources.has(source.id));
  let isLoading = $derived(source && metadata.loadingSources.has(source.id));
  let sourceCollapsed = $state(source && metadata.collapsedSources.has(source.id));

  let hasCals = $derived(repository.calendars.filter(x => x.source === source.id).length > 0);

  let showModal: ((source: SourceModel) => Promise<SourceModel>) = getContext("showSourceModal");
  function showModalInternal() {
    showModal(source).then(newSource => source = newSource).catch(NoOp);
  }

  $effect(() => {
    sourceCollapsed = metadata.collapsedSources.has(source.id);
  });
  $effect(() => {
    if (source && source.id) metadata.setSourceCollapse(source.id, sourceCollapsed);
  });

  async function reorderSource(newIndex: number) {
    await repository.changeSourceDisplayOrder(source, newIndex).catch((err) => {
      queueNotification(ColorKeys.Danger, err);
    });
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/colors.scss";

  div {
    color: color-mix(in srgb, colors.$foregroundPrimary 50%, transparent);
    height: 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    align-content: center;
    user-select: none;
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

<div
  class="sourceEntry"
  use:draggable={{ ownClass: "sourceEntry", childClasses: ["calendarEntry"], callback: reorderSource}}
  id={sourceEntryId}
  aria-label={t("source.aria", { values: { name: source.name } })}
  aria-owns={ariaControls}
  aria-level="1"
  role="treeitem"
  aria-selected={document.activeElement?.id === sourceEntryId}
  aria-details={hasErrored ? extendUniqueElementId(["tooltip"], sourceEntryId) : undefined}
>
  <button onclick={showModalInternal} use:focusIndicator={{ type: "underline" }}>
    {source.name}
  </button>
  <span>
    {#if isLoading}
      <Spinner/>
    {/if}
    {#if hasCals}
      <CollapseToggle bind:collapsed={sourceCollapsed} ariaControls={ariaControls}/>
    {/if}
    {#if hasErrored}
      <Tooltip describes={sourceEntryId} error={true}>{t("source.error.calendars.tooltip")}</Tooltip>
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
