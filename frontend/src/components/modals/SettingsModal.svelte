<script lang="ts">
  import { Code, LockKeyhole, Monitor, User } from "lucide-svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";

  interface Props {
    showModal?: () => any;
    hideModal?: () => any;
  }

  let {
    showModal = $bindable(),
    hideModal = $bindable(NoOp),
  }: Props = $props();

  showModal = () => {
    showModalInternal();
  };

  hideModal = () => {
    hideModalInternal();
  };

  let showModalInternal = $state(NoOp);
  let hideModalInternal = $state(NoOp);

  let selectedCategory = $state("general");
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div {
    box-sizing: border-box;
    width: 50vw;
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 1fr;
    gap: dimensions.$gapMiddle;
  }
</style>

<Modal
  title={"Settings"}
  bind:showModal={showModalInternal}
  bind:hideModal={hideModalInternal}
>
  <div class="container">
    <ButtonList
      bind:value={selectedCategory}
      options={[
        [
          { name: "Account", value: "general", icon: User },
          { name: "Appearance", value: "appearance", icon: Monitor },
          { name: "Developer", value: "developer", icon: Code }
        ],
        [
          { name: "Administrative", value: "admin", icon: LockKeyhole },
        ],
      ]} 
    />
    <main>
      {#if selectedCategory === "general"}
        General Settings
      {:else if selectedCategory === "appearance"}
        Appearance Settings
      {:else if selectedCategory === "developer"}
        Developer Settings
      {:else if selectedCategory === "admin"}
        Administrative Settings
      {/if}
    </main>
  </div>
</Modal>