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

  import { isValidEmail, isValidInviteCode, isValidPassword, isValidRepeatPassword, isValidUsername, valid } from '$lib/client/validation';
  import { queueNotification } from '$lib/client/notifications';
  import Title from '../../components/layout/Title.svelte';
  import Paragraph from '../../components/layout/Paragraph.svelte';
  import Box from '../../components/layout/Box.svelte';
  import { browser } from '$app/environment';

  interface PageProps {
    form: ActionData;
    data: {
      registrationEnabled: boolean;
    }
  }

  let {
    form = $bindable(),
    data,
  }: PageProps = $props();

  afterNavigate(() => {
    if (form?.error) queueNotification(ColorKeys.Danger, form.error);
    if (browser) localStorage.clear();
  });

  const redirect = $derived(page.url.searchParams.get('redirect') || "/");

  let password: string = $state("");

  let usernameValidity: Validity = $state(valid);
  let emailValidity: Validity = $state(valid);
  let passwordValidity: Validity = $state(valid);
  let passwordRepeatValidity: Validity = $state(valid);

  let inviteCode: string = $state("");
  let inviteCodeValidity: Validity = $state(valid);

  function inviteCodeFormatting(_, event: Event | null) {
    if (inviteCode.length == 0) return;

    // To simplify formatting, delete all hyphens first
    let tmpCode = inviteCode;
    tmpCode = tmpCode.replace(/-/g, "");

    // Remove all leading and trailing spaces
    tmpCode = tmpCode.trim();

    // Make sure the invite code is upper case
    tmpCode = tmpCode.toUpperCase();

    // Remove all illegal characters
    for (let i = 0; i < tmpCode.length;) {
      if (!tmpCode[i].match(/[A-Z0-9]/)) tmpCode = tmpCode.slice(0, i) + tmpCode.slice(i + 1);
      else i++;
    }

    // If a backspace or delete key deleted a hyphen, remove the character before or after it
    let amountOfHyphens = 0;
    for (let i = 0; i < inviteCode.length; i++) if (inviteCode[i] == "-") amountOfHyphens++;
    if (event && "inputType" in event && ((inviteCode.length >= 9 && amountOfHyphens == 1) || (inviteCode.length >= 4 && amountOfHyphens == 0))) {
      if (event.inputType === "deleteContentBackward") {
        if (inviteCode.length == 4 || (inviteCode.length > 4 && inviteCode[4] != "-")) {
          tmpCode = tmpCode.slice(0, 3) + tmpCode.slice(4);
        } else if (inviteCode.length == 9 || (inviteCode.length > 9 && inviteCode[9] != "-")) {
          tmpCode = tmpCode.slice(0, 7) + tmpCode.slice(8);
        }
      } else if (event.inputType === "deleteContentForward") {
        if (inviteCode.length >= 9 && inviteCode[8] != "-") {
          tmpCode = tmpCode.slice(0, 8) + tmpCode.slice(9);
        } else if (inviteCode.length >= 4 && inviteCode[3] != "-") {
          tmpCode = tmpCode.slice(0, 4) + tmpCode.slice(5);
        }
      }
    }

    // Remove all characters past the invite code length
    if (tmpCode.length > 12) tmpCode = tmpCode.slice(0, 12);

    // Add hyphens in the right places
    if (tmpCode.length >= 4) tmpCode = tmpCode.slice(0, 4) + "-" + tmpCode.slice(4);
    if (tmpCode.length >= 9) tmpCode = tmpCode.slice(0, 9) + "-" + tmpCode.slice(9);

    // Set the invite code to the formatted one
    inviteCode = tmpCode;
  }

  let canSubmit: boolean = $derived(
    usernameValidity?.valid && emailValidity?.valid && passwordValidity?.valid && passwordRepeatValidity?.valid
    && (data.registrationEnabled || inviteCodeValidity.valid)
  );
</script>

<SimplePage>
  {#if data.registrationEnabled}
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
        name="password_repeat"
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
  {:else}
    <Box>
      <Title>Registration</Title>
      <Paragraph>
        To register, please input an invite code.
        Please contact the administrator if you don't have one.
      </Paragraph>
      <TextInput
        name="invite_code"
        placeholder="Invite Code"
        bind:value={inviteCode}
        validation={isValidInviteCode}
        bind:validity={inviteCodeValidity}
        onChange={inviteCodeFormatting}
        onInput={inviteCodeFormatting}
      />
      <Link href="/login?redirect={encodeURIComponent(redirect)}">Already got an account?</Link>
    </Box>
  {/if}
</SimplePage>