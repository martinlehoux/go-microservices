export interface AuthenticationGateway {
  register(payload: RegisterPayload): Promise<void>;
  authenticate(payload: AuthenticatePayload): Promise<AuthenticateResponse>;
}

type RegisterPayload = {
  identifier: string;
  password: string;
};

type AuthenticatePayload = {
  identifier: string;
  password: string;
};

type AuthenticateResponse = {
  token: string;
};

export class ApiAuthenticationGateway implements AuthenticationGateway {
  private readonly apiUrl = new URL("http://localhost:3000");

  async authenticate(
    payload: AuthenticatePayload
  ): Promise<AuthenticateResponse> {
    const url = new URL(`/auth/authenticate`, this.apiUrl);
    const res = await fetch(url.toString(), {
      method: "POST",
      body: JSON.stringify(payload),
      headers: { "Content-Type": "application/json" },
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`authentication failed: ${data.error}`);
    }
    return data;
  }

  async register(payload: RegisterPayload): Promise<void> {
    const url = new URL(`/auth/register`, this.apiUrl);
    const res = await fetch(url.toString(), {
      method: "POST",
      body: JSON.stringify(payload),
      headers: { "Content-Type": "application/json" },
    });
    const data = await res.json();
    if (!res.ok) {
      throw new Error(`registration failed: ${data.error}`);
    }
    return;
  }
}
