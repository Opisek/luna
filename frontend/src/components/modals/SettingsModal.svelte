<script lang="ts">
  import { Code, LockKeyhole, Monitor, User } from "lucide-svelte";
  import { NoOp } from "../../lib/client/placeholders";
  import ButtonList from "../forms/ButtonList.svelte";
  import Modal from "./Modal.svelte";
  import TextInput from "../forms/TextInput.svelte";
  import { isValidEmail, isValidPassword, isValidUsername } from "../../lib/client/validation";
  import ToggleInput from "../forms/ToggleInput.svelte";
  import SelectButtons from "../forms/SelectButtons.svelte";
  import Image from "../layout/Image.svelte";
  import FileUpload from "../forms/FileUpload.svelte";
  import Horizontal from "../layout/Horizontal.svelte";

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

  const categories: Option[][] = [
    [
      { name: "Account", value: "account", icon: User },
      { name: "Appearance", value: "appearance", icon: Monitor },
      { name: "Developer", value: "developer", icon: Code }
    ],
    [
      { name: "Administrative", value: "admin", icon: LockKeyhole },
    ],
  ]
  let selectedCategory = $state("account");

  // Account Settings
  let profilePictureType = $state("gravatar");
  let profilePictureFiles: FileList | null = $state(null);
</script>

<style lang="scss">
  @use "../../styles/dimensions.scss";

  div.container {
    box-sizing: border-box;
    display: grid;
    grid-template-columns: auto 1fr;
    grid-template-rows: 1fr;
    gap: dimensions.$gapMiddle;
    min-width: 30vw;
    min-height: 50vh;
  }

  main {
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
  }

  div.pfpButtons {
    display: flex;
    flex-direction: column;
    gap: dimensions.$gapMiddle;
    width: 100%;
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
      options={categories} 
    />
    <main>
      {#if selectedCategory === "account"}
        <TextInput
          name="username"
          placeholder="Username"
          validation={isValidUsername}
        />
        <TextInput
          name="email"
          placeholder="Email"
          validation={isValidEmail}
        />
        <TextInput
          name="password"
          placeholder="New Password"
          password={true}
          validation={isValidPassword}
        />
        <ToggleInput
          name="searchable" 
          description="Allow other users to find me"
        />
        <Horizontal position="justify">
          <div class="pfpButtons">
            <SelectButtons
              name="pfp_type"
              placeholder="Profile Picture"
              bind:value={profilePictureType}
              options={[
                { name: "Gravatar", value: "gravatar" },
                { name: "Upload File", value: "database" },
                { name: "Internet Link", value: "remote" }
              ]}
            />
          </div>
          <Image
            src="https://opisek.net/_app/immutable/assets/portrait.2-Ny2279.webp/"
            alt="Profile Picture"
          />
        </Horizontal>
        {#if profilePictureType === "database"}
          <FileUpload
            name="pfp_file"
            placeholder="Profile Picture File"
            files={profilePictureFiles}
          />
        {:else if profilePictureType === "remote"}
          <TextInput
            name="pfp_link"
            placeholder="Profile Picture Link"
          />
        {/if}
      {:else if selectedCategory === "appearance"}
        Appearance Settings
      {:else if selectedCategory === "developer"}
        Developer Settings
      {:else if selectedCategory === "admin"}
        <ToggleInput
          name="debug_mode" 
          description="Display IDs"
        />
      {/if}
    </main>
  </div>
</Modal>