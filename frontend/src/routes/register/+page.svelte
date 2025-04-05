<script lang="ts">
  import type { ActionData } from './$types';

  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import SimplePage from '../../components/layout/SimplePage.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';
  import ToggleInput from '../../components/forms/ToggleInput.svelte';
  import { ColorKeys } from '../../types/colors';

  import { afterNavigate } from '$app/navigation';
  import { page } from '$app/state';

  import { isValidEmail, isValidPassword, isValidRepeatPassword, isValidUsername, valid } from '$lib/client/validation';
  import { queueNotification } from '$lib/client/notifications';

  interface Props {
    form: ActionData;
  }

  let { form = $bindable() }: Props = $props();

  afterNavigate(() => {
    if (form?.error) queueNotification(ColorKeys.Danger, form.error);
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