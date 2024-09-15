<script lang="ts">
  import type { ActionData } from './$types';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  export let form: ActionData;
</script>

{ #if form?.error }
  <p>Error: {form.error}</p>
{/if}
<form method="POST">
  <input name="username" type="text" placeholder="Username" />
  <input name="password" type="password" placeholder="Password" />
  <button type="submit">Login</button>
  <a href="/register">Register</a>
</form>