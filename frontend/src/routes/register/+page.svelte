<script lang="ts">
  import type { ActionData } from './$types';

  import ToggleInput from '../../components/forms/ToggleInput.svelte';
  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';

  import { beforeNavigate } from '$app/navigation';
  import { page } from '$app/state';

  import { isValidEmail, isValidPassword, isValidRepeatPassword, isValidUsername, valid } from '$lib/client/validation';
  import { queueNotification } from '$lib/client/notifications';
  import SimplePage from '../../components/layout/SimplePage.svelte';

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

  let password: string = $state("");

  let usernameValidity: Validity = $state(valid);
  let emailValidity: Validity = $state(valid);
  let passwordValidity: Validity = $state(valid);
  let passwordRepeatValidity: Validity = $state(valid);

  let canSubmit: boolean = $derived(usernameValidity?.valid && emailValidity?.valid && passwordValidity?.valid && passwordRepeatValidity?.valid);
</script>

<SimplePage>
  <Form title="Register" submittable={canSubmit}>
    <TextInput
      name="username"
      placeholder="Username"
      validation={isValidUsername}
      bind:validity={usernameValidity}
    />
    <TextInput
      name="email"
      placeholder="Email"
      validation={isValidEmail}
      bind:validity={emailValidity}
    />
    <TextInput
      name="password"
      placeholder="Password"
      password={true}
      validation={isValidPassword}
      bind:value={password}
      bind:validity={passwordValidity}
    />
    <TextInput
      name="passwordRepeat"
      placeholder="Repeat Password"
      password={true}
      validation={isValidRepeatPassword(password)}
      bind:validity={passwordRepeatValidity}
    />
    <ToggleInput
      name="remember"
      description="Remember me"
    />
    <Link href="/login?redirect={encodeURIComponent(redirect)}">Already got an account?</Link>
  </Form>
</SimplePage>