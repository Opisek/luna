<script lang="ts">
  import { run } from 'svelte/legacy';

  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';
  import type { ActionData } from './$types';
  import { queueNotification } from '../../lib/client/notifications';
  import { page } from '$app/stores';
  import CheckboxInput from '../../components/forms/CheckboxInput.svelte';
  import { isValidEmail, isValidPassword, isValidRepeatPassword, isValidUsername } from '../../lib/client/validation';
  import { beforeNavigate } from '$app/navigation';

  interface Props {
    form: ActionData;
  }

  let { form = $bindable() }: Props = $props();

  // TODO: easier way to prevent double notifications?
  let alreadyShownError = $state(false);
  run(() => {
    ((form) => {
      if (form?.error && !alreadyShownError) {
        alreadyShownError = true;
        queueNotification("failure", form.error);
      }
    })(form)
  });
  beforeNavigate(() => {
    form = null;
    alreadyShownError = false;
  });

  const redirect = $page.url.searchParams.get('redirect') || "/";

  let password: string = $state("");

  let usernameValidity: Validity = $state();
  let emailValidity: Validity = $state();
  let passwordValidity: Validity = $state();
  let passwordRepeatValidity: Validity = $state();

  let canSubmit: boolean = $derived(usernameValidity?.valid && emailValidity?.valid && passwordValidity?.valid && passwordRepeatValidity?.valid);
  
</script>

<style lang="scss">
  div {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
  }
</style>

<div>
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
    <CheckboxInput
      name="remember"
      description="Remember me"
    />
    <Link href="/login?redirect={encodeURIComponent(redirect)}">Already got an account?</Link>
  </Form>
</div>


