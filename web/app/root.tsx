import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useRouteLoaderData,
  type LinksFunction,
  type LoaderFunctionArgs,
} from "react-router";
import { Toaster } from "sonner";

import { getLocale, messages } from "@/lib/i18n";

import stylesheet from "./app.css?url";

export const links: LinksFunction = () => [
  { rel: "stylesheet", href: stylesheet },
];

export function loader({ request }: LoaderFunctionArgs) {
  return { locale: getLocale(request) };
}

export function headers() {
  return { Vary: "Accept-Language" };
}

export function Layout({ children }: { children: React.ReactNode }) {
  const locale = useRouteLoaderData<typeof loader>("root")?.locale ?? "en";

  return (
    <html lang={locale}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <Toaster position="top-right" richColors closeButton />
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

export function ErrorBoundary({ error }: { error: unknown }) {
  const locale = useRouteLoaderData<typeof loader>("root")?.locale ?? "en";
  const text = messages[locale].root;
  const status = isRouteErrorResponse(error) ? error.status : 500;
  const message =
    status === 404
      ? text.notFound
      : isRouteErrorResponse(error) && typeof error.data === "string"
      ? error.data
      : text.pageLoadFailed;

  return (
    <main className="grid min-h-screen place-items-center px-6">
      <section className="max-w-md text-center">
        <p className="font-mono text-sm text-muted-foreground">{status}</p>
        <h1 className="mt-3 text-2xl font-semibold">{text.problem}</h1>
        <p className="mt-3 text-sm leading-6 text-muted-foreground">{message}</p>
        <a
          className="mt-6 inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          href="/"
        >
          {text.backToConsole}
        </a>
      </section>
    </main>
  );
}
