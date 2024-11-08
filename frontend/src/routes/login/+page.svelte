<script lang="ts">
  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';
  import type { ActionData } from './$types';
  import { queueNotification } from '../../lib/client/notifications';
  import { page } from '$app/stores';
  import CheckboxInput from '../../components/forms/CheckboxInput.svelte';
  import { isValidPassword, isValidUsername } from '../../lib/client/validation';

  export let form: ActionData;

  $: if (form?.error) queueNotification("failure", form.error);

  const redirect = $page.url.searchParams.get('redirect') || "/";
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
  <Form title="Login">
    <TextInput name="username" placeholder="Username" validation={isValidUsername}/>
    <TextInput name="password" placeholder="Password" password={true} validation={isValidPassword}/>
    <CheckboxInput name="remember" description="Remember me"/>
    <Link href="/register?redirect={encodeURIComponent(redirect)}">No account yet?</Link>
    <Link href="/recover?redirect={encodeURIComponent(redirect)}">Forgot password?</Link>
  </Form>
</div>