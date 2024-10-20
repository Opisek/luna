<script lang="ts">
  import Form from '../../components/forms/Form.svelte';
  import Link from '../../components/forms/Link.svelte';
  import TextInput from '../../components/forms/TextInput.svelte';
  import type { ActionData } from './$types';
  import { queueNotification } from '../../lib/client/notifications';
  import { page } from '$app/stores';

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
  <Form title="Register">
    <TextInput name="username" placeholder="Username"/>
    <TextInput name="email" placeholder="Email"/>
    <TextInput name="password" placeholder="Password" password={true}/>
    <TextInput name="passwordRepeat" placeholder="Repeat Password" password={true}/>
    <Link href="/login?redirect={encodeURIComponent(redirect)}">Already got an account?</Link>
  </Form>
</div>


