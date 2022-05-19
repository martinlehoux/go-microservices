import { User } from "../resources/users";

export interface HumanResourcesGateway {
  setToken(token: string): void;
  register(payload: RegisterPayload): Promise<void>;
  listUsers(): Promise<User[]>;
}

type RegisterPayload = {
  preferred_name: string;
};

export class ApiHumanResourcesGateway implements HumanResourcesGateway {
  private readonly apiUrl = new URL("http://localhost:3000");
  private token: string | null = null;

  setToken(token: string): void {
    this.token = token;
  }

  async register(payload: RegisterPayload): Promise<void> {
    this.checkToken(this.token);
    const url = new URL(`/hr/users/register`, this.apiUrl);
    const res = await fetch(url.toString(), {
      method: "POST",
      body: JSON.stringify(payload),
      headers: this.headers,
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`registration failed: ${data.error}`);
    }
    return;
  }

  async listUsers(): Promise<User[]> {
    const url = new URL(`/hr/users`, this.apiUrl);
    const res = await fetch(url.toString(), { headers: this.headers });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`listing users failed: ${data.error}`);
    }
    return data.items;
  }

  private checkToken(token: string | null): asserts token is string {
    if (token === null) {
      throw new Error("missing token");
    }
  }

  private get headers() {
    const headers = new Headers();
    headers.set("Content-Type", "application/json");
    if (this.token) {
      headers.set("Authorization", `Bearer ${this.token}`);
    }
    return headers;
  }
}
