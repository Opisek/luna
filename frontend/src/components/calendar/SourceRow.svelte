<script lang="ts">
  import { BoltIcon, CogIcon, PencilIcon } from "lucide-svelte";
  import IconButton from "../interactive/IconButton.svelte";
  import SourceModal from "../modals/SourceModal.svelte";

  export let source: SourceModel;

  let showModal: () => any;
</script>

<style lang="scss">
  @import "../../styles/animations.scss";
  @import "../../styles/dimensions.scss";
  @import "../../styles/colors.scss";

  span {
    margin: 0;
    color: $foregroundFaded;
    height: 1.25em;
    display: flex;
    justify-content: space-between;
    align-items: center;
    align-content: center;
    cursor: pointer;
  }


  //span :global(button) {
  //  opacity: 0;
  //  transition: all $animationSpeed $cubic;
  //}

  //span:hover :global(button) {
  //  opacity: 1;
  //}
</style>

<span on:click={showModal}>
  {source.name}
  <!--
  <IconButton callback={showModal}>
    <PencilIcon size={16}/>
    <BoltIcon size={16}/>
    <CogIcon size={16}/>
  </IconButton>
  -->
</span>

<SourceModal
  bind:showModal={showModal}
  source={source}
/>