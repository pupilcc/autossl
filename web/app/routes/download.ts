import type { LoaderFunctionArgs } from "react-router";

import { downloadFromBackend } from "@/lib/backend.server";

export async function loader({ params, request }: LoaderFunctionArgs) {
  const file = params.file ?? "";
  if (!/^[A-Za-z0-9_-]+\.(crt|key)$/.test(file)) {
    throw new Response("Not found", { status: 404 });
  }

  const response = await downloadFromBackend(file, request.method === "HEAD" ? "HEAD" : "GET");
  const headers = new Headers();
  for (const name of ["content-type", "content-length", "etag", "last-modified"]) {
    const value = response.headers.get(name);
    if (value) headers.set(name, value);
  }

  return new Response(request.method === "HEAD" ? null : response.body, {
    status: response.status,
    headers,
  });
}
