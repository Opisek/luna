import { fetchJson, fetchResponse } from "../net";
import type { UserData } from "../../../types/settings";

class Users {
  public currentUser = $state("");
  public users = $state<UserData[]>([]);

  public async fetchAll() {
    await fetchJson("/api/users?all=true").then((data: { current: string, users: UserData[] }) => {
      this.currentUser = data.current;
      data.users.forEach(async x => {
        x.created_at = new Date((x.created_at as unknown as string).replace("Z", ""));
      });
      this.users = data.users;
    });
  }

  public async disableUser(id: string) {
    return fetchResponse(`/api/users/${id}/disable`, { method: "POST" }).then(() => {
      this.users = this.users.map(user => {
        if (user.id === id) user.enabled = false;
        return user;
      });
    });
  }

  public async enableUser(id: string) {
    return fetchResponse(`/api/users/${id}/enable`, { method: "POST" }).then(() => {
      this.users = this.users.map(user => {
        if (user.id === id) user.enabled = true;
        return user;
      });
    });
  }

  public async deleteUser(id: string, password: string) {
    const formData = new FormData();
    formData.append("password", password);

    return fetchResponse(`/api/users/${id}/delete`, { method: "DELETE", body: formData }).then(() => {
      this.users = this.users.filter(user => user.id !== id);
    });
  }
}

let users: Users | null = $state(null);
export function getUsers() {
  if (users === null) {
    users = new Users();
  }
  return users;
}

export function resetUsers() {
  users = null;
}