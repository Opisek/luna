<script lang="ts">
  import type { ActionData } from './$types';

  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import SimplePage from '../../components/layout/SimplePage.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';
  import ToggleInput from '../../components/forms/ToggleInput.svelte';
  import { ColorKeys } from '../../types/colors';

  import { beforeNavigate } from '$app/navigation';
  import { page } from '$app/state';

  import { isValidPassword, isValidUsername, valid } from '$lib/client/validation';
  import { queueNotification } from '$lib/client/notifications';

  interface Props {
    form: ActionData;
  }

  let { form = $bindable() }: Props = $props();

  // TODO: easier way to prevent double notifications?
  let alreadyShownError = $state(false);
  $effect(() => {
    ((form) => {
      if (form?.error && !alreadyShownError) {
        alreadyShownError = true;
        queueNotification(ColorKeys.Danger, form.error);
      }
    })(form)
  });
  beforeNavigate(() => {
    form = null;
    alreadyShownError = false;
  });

  const redirect = $derived(page.url.searchParams.get('redirect') || "/");

  let usernameValidity: Validity = $state(valid);
  let passwordValidity: Validity = $state(valid);

  let canSubmit: boolean = $derived(usernameValidity?.valid && passwordValidity?.valid);
</script>

<SimplePage>
  <Form title="Login" submittable={canSubmit}>
    <TextInput
      name="username"
      placeholder="Username"
      validation={isValidUsername}
      bind:validity={usernameValidity}
    />
    <TextInput
      name="password"
      placeholder="Password"
      password={true}
      validation={isValidPassword}
      bind:validity={passwordValidity}
    />
    <ToggleInput
      name="remember"
      description="Remember me"
    />
    <Link href="/register?redirect={encodeURIComponent(redirect)}">No account yet?</Link>
    <Link href="/recover?redirect={encodeURIComponent(redirect)}">Forgot password?</Link>
  </Form>
</SimplePage>