const apiBaseUrl = (process.env.AUTOSSL_API_URL ?? "http://127.0.0.1:1323").replace(
  /\/$/,
  "",
);

export type Certificate = {
  code: string;
  domain: string;
  dnsNames: string[] | null;
  cert: string;
  key: string;
};

export class BackendError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
  }
}

async function backendRequest<T>(
  path: string,
  init: RequestInit = {},
  token?: string,
): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);

  const response = await fetch(`${apiBaseUrl}${path}`, { ...init, headers });
  if (!response.ok) {
    let message = `后端请求失败（${response.status}）`;
    try {
      const body = (await response.json()) as { message?: string };
      if (body.message) message = body.message;
    } catch {
      // Keep the status-based message for non-JSON responses.
    }
    throw new BackendError(message, response.status);
  }

  if (response.status === 204 || response.headers.get("content-length") === "0") {
    return undefined as T;
  }
  return (await response.json()) as T;
}

export async function login(username: string, password: string) {
  return backendRequest<{ token: string }>("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
}

export function listCertificates(token: string) {
  return backendRequest<Certificate[]>("/list", {}, token);
}

export function createCertificate(token: string, domain: string) {
  return backendRequest<void>(
    "/generate",
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ domain }),
    },
    token,
  );
}

export function deleteCertificate(token: string, code: string) {
  return backendRequest<void>(
    `/${encodeURIComponent(code)}`,
    { method: "DELETE" },
    token,
  );
}

export function downloadFromBackend(file: string, method: string) {
  return fetch(`${apiBaseUrl}/dl/${encodeURIComponent(file)}`, { method });
}
