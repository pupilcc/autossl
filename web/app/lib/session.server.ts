import { randomUUID } from "node:crypto";
import { redirect } from "react-router";

const COOKIE_NAME = "autossl_session";
const MAX_AGE_SECONDS = 72 * 60 * 60;

type StoredSession = {
  token: string;
  expiresAt: number;
};

declare global {
  var autoSslSessions: Map<string, StoredSession> | undefined;
}

const sessions = (globalThis.autoSslSessions ??= new Map());

function cookieValue(request: Request) {
  const cookie = request.headers.get("Cookie");
  if (!cookie) return null;

  for (const part of cookie.split(";")) {
    const [name, ...value] = part.trim().split("=");
    if (name === COOKIE_NAME) return decodeURIComponent(value.join("="));
  }
  return null;
}

function secureAttribute(request: Request) {
  const protocol = request.headers.get("X-Forwarded-Proto") ?? new URL(request.url).protocol;
  return protocol.replace(":", "") === "https" ? "; Secure" : "";
}

function serializeCookie(request: Request, value: string, maxAge: number) {
  return `${COOKIE_NAME}=${encodeURIComponent(value)}; Path=/; HttpOnly; SameSite=Strict; Max-Age=${maxAge}${secureAttribute(request)}`;
}

export function getSession(request: Request) {
  const id = cookieValue(request);
  if (!id) return null;

  const session = sessions.get(id);
  if (!session || session.expiresAt <= Date.now()) {
    sessions.delete(id);
    return null;
  }
  return { id, ...session };
}

export function requireSession(request: Request) {
  const session = getSession(request);
  if (!session) throw redirect("/login");
  return session;
}

export function createSessionCookie(request: Request, token: string) {
  const id = randomUUID();
  sessions.set(id, {
    token,
    expiresAt: Date.now() + MAX_AGE_SECONDS * 1000,
  });
  return serializeCookie(request, id, MAX_AGE_SECONDS);
}

export function destroySessionCookie(request: Request) {
  const id = cookieValue(request);
  if (id) sessions.delete(id);
  return serializeCookie(request, "", 0);
}
