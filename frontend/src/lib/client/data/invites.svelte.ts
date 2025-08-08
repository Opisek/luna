import { page } from "$app/state";

import { fetchJson, fetchResponse } from "../net";

export class RegistrationInvites {
  public activeInvites = $state<RegistrationInvite[]>([]);

  public async fetch() {
    await fetchJson("/api/invites").then((data: { invites: RegistrationInvite[] }) => {
      data.invites.forEach(x => {
        x.created_at = new Date((x.created_at as unknown as string).replace("Z", ""));
        x.expires_at = new Date((x.expires_at as unknown as string).replace("Z", ""));
      })
      this.activeInvites = data.invites;
    });
  }

  public async createInvite(duration: number) {
    const formData = new FormData();
    formData.append("duration", duration.toString());

    return await fetchJson("/api/invites", { method: "PUT", body: formData }).then((data: { invite: RegistrationInvite }) => {
      data.invite.created_at = new Date((data.invite.created_at as unknown as string).replace("Z", ""));
      data.invite.expires_at = new Date((data.invite.expires_at as unknown as string).replace("Z", ""));
      this.activeInvites.push(data.invite);
      return data.invite;
    });
  }

  public async revokeInvite(id: string) {
    return fetchResponse(`/api/invites/${id}`, { method: "DELETE" }).then(() => {
      this.activeInvites = this.activeInvites.filter(x => x.invite_id != id);
    });
  }

  public async revokeInvites() {
    return fetchResponse(`/api/invites`, { method: "DELETE" }).then(() => {
      this.activeInvites = [];
    });
  }
}

export function getRegistrationInvites() {
  return page.data.singletons.invites;
}