<script lang="ts">
  import Horizontal from "../layout/Horizontal.svelte";
  import Title from "../layout/Title.svelte";
  import type { Snippet } from "svelte";
  import { ColorKeys } from "../../types/colors";
  import { enhance } from "$app/forms";
  import type { ActionResult } from "@sveltejs/kit";
  import { NoOp } from "../../lib/client/placeholders";
  import IconButton from "../interactive/IconButton.svelte";
  import { Send } from "lucide-svelte";

  interface Props {
    title: string;
    submittable?: boolean;
    callback?: (result: ActionResult) => void;
    children?: Snippet;
  }

  let {
    title,
    submittable = true,
    callback = NoOp,
    children,
  }: Props = $props();

  let registerButtonPromise: (promise: Promise<any>) => void = $state(NoOp);
  let promiseResult = $state(NoOp);

  function onSubmit(e: SubmitEvent) {
    if (submittable) {
      registerButtonPromise(new Promise<void>(res => {
        promiseResult = () => {
          promiseResult = NoOp;
          res();
        }
      }));
    } else {
      e.preventDefault();
      return; // TODO: add some user feedback when a form fails to submit
    }
  }

  function onResult({ result, update }: { result: ActionResult; update: () => void }) {
    promiseResult();
    callback(result);
    update();
  }
</script>

<style lang="scss">
  @use "../../styles/animations.scss";
  @use "../../styles/colors.scss";
  @use "../../styles/dimensions.scss";
  @use "../../styles/decorations.scss";

  form {
    border-radius: dimensions.$borderRadius;
    max-width: 50vw;
    min-width: 30em;
    padding: dimensions.$gapLarge dimensions.$gapLarger dimensions.$gapLarger dimensions.$gapLarger;
    border-radius: dimensions.$borderRadius;
    display: flex;
    flex-direction: column;
    flex-wrap: nowrap;
    gap: dimensions.$gapMiddle;
    box-shadow: decorations.$boxShadow;
  }
</style>

<form method="POST" onsubmit={onSubmit} use:enhance={() => onResult}>
  <Title>{title}</Title>
  {@render children?.()}
  <Horizontal position="right">
    <IconButton type="submit" bind:externalLoading={registerButtonPromise} color={ColorKeys.Success} enabled={submittable} alt="Submit" canRenderAsButton={true}>
      <Send/>
    </IconButton>
  </Horizontal>
</form>