// This class periodically checks /api/health to determine the reachability of the different parts of the application.

import { fetchJson } from "./net";

export enum Reachability {
  Unkown,
  None,
  Frontend,
  Backend,
  Database
}

class Connectivity {
  private reachable: Reachability;
  private promise: Promise<Reachability> | null = null;
  
  constructor() {
    this.reachable = Reachability.Database;
  }

  private async checkConnection() {
    await fetchJson("/api/health").then(() => {
      this.reachable = Reachability.Database
    }).catch((err) => {
      if (err.message === "Could not contact server" || err.message.includes("NetworkError")) {
        this.reachable = Reachability.None;
      } else if (err.message === "The backend is not reachable") {
        this.reachable = Reachability.Frontend;
      } else if (err.message === "Database error") {
        this.reachable = Reachability.Backend;
      } else {
        this.reachable = Reachability.Unkown;
      }
    })

    this.promise = null;
    return this.reachable;
  }

  async check() {
    if (this.promise === null) {
      this.promise = this.checkConnection();
    }
    return this.promise;
  }
}

let connectivity: Connectivity | null = null;
export function getConnectivity() {
  if (connectivity === null) {
    connectivity = new Connectivity();
  }
  return connectivity;
}
