import { page } from "$app/state";

import { fetchJson, fetchResponse } from "../net";

export class OauthClients {
  public clients = $state<OauthClientModel[]>([]);
  public clientTokens = $state<Map<string, OauthTokensModel[]>>(new Map());
  public tokenClients = $state<Map<string, string>>(new Map());

  public async fetch() {
    await fetchJson("/api/oauth/clients").then((data: { clients: OauthClientModel[] }) => {
      this.clients = data.clients;
    });
  }

  public async fetchTokens() {
    await fetchJson("/api/oauth/tokens").then((data: { tokens: OauthTokensModel[] }) => {
      const clientsMap = new Map<string, OauthTokensModel[]>();
      const tokensMap = new Map<string, string>();
      for (const tokensEntry of data.tokens) {
        if (!clientsMap.has(tokensEntry.client_id)) clientsMap.set(tokensEntry.client_id, []);
        clientsMap.get(tokensEntry.client_id)?.push(tokensEntry)
        tokensMap.set(tokensEntry.id, tokensEntry.client_id);
      }
      this.clientTokens = clientsMap;
      this.tokenClients = tokensMap;
    });
  }

  public async getClientDetails(id: string) {
    return await fetchJson(`/api/oauth/clients/${id}`).then((data: { client: OauthClientModel }) => {
      return data.client;
    });
  }

  public async registerClient(client: OauthClientModel) {
    const formData = new FormData();
    formData.append("name", client.name);
    formData.append("client_id", client.client_id);
    formData.append("client_secret", client.client_secret);
    formData.append("base_url", client.base_url);
    formData.append("scope", client.scope);

    return await fetchJson("/api/oauth/clients", { method: "PUT", body: formData }).then((data: { client: OauthClientModel }) => {
      this.clients.push(data.client);
      return data.client;
    });
  }

  public async updateClient(client: OauthClientModel) {
    const formData = new FormData();
    formData.append("name", client.name);
    formData.append("client_id", client.client_id);
    formData.append("client_secret", client.client_secret);
    formData.append("base_url", client.base_url);
    formData.append("scope", client.scope);
    
    return await fetchJson(`/api/oauth/clients/${client.id}`, { method: "PATCH", body: formData }).then((data: { client: OauthClientModel }) => {
      this.clients = this.clients.map(x => x.id == client.id ? data.client : x);
      return data.client;
    });
  }

  public async deleteClient(id: string) {
    return fetchResponse(`/api/oauth/clients/${id}`, { method: "DELETE" }).then(() => {
      this.clients = this.clients.filter(x => x.id != id);
    });
  }
}

export function getOauthClients(): OauthClients {
  return page.data.singletons.oauthClients;
}