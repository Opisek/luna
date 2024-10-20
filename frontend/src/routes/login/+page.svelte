<script lang="ts">
  import type { ActionData } from './$types';
  import TextInput from '../../components/forms/TextInput.svelte';
  import Form from '../../components/forms/Form.svelte';
  import { queueNotification } from '../../lib/client/notifications';

  export let form: ActionData;

  $: if (form?.error) queueNotification("failure", form.error);
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
    <TextInput name="username" placeholder="Username"/>
    <TextInput name="password" placeholder="Password" password={true}/>
  </Form>
</div>