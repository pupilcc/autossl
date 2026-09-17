import { Eye, EyeOff, KeyRound, LoaderCircle, ShieldCheck } from "lucide-react";
import { useState } from "react";
import {
  Form,
  redirect,
  useActionData,
  useLoaderData,
  useNavigation,
  type ActionFunctionArgs,
  type LoaderFunctionArgs,
  type MetaFunction,
} from "react-router";

import { LanguageSwitcher } from "@/components/language-switcher";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { BackendError, login } from "@/lib/backend.server";
import { getLocale, messages } from "@/lib/i18n";
import {
  assertSameOrigin,
  createSessionCookie,
  getSession,
} from "@/lib/session.server";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
  { title: messages[data?.locale ?? "en"].login.metaTitle },
];

export function loader({ request }: LoaderFunctionArgs) {
  if (getSession(request)) return redirect("/");
  return { locale: getLocale(request) };
}

export async function action({ request }: ActionFunctionArgs) {
  assertSameOrigin(request);
  const text = messages[getLocale(request)].login;
  const formData = await request.formData();
  const username = String(formData.get("username") ?? "").trim();
  const password = String(formData.get("password") ?? "");

  if (!username || !password) return { error: text.missingCredentials };

  try {
    const result = await login(username, password);
    return redirect("/", {
      headers: { "Set-Cookie": createSessionCookie(request, result.token) },
    });
  } catch (error) {
    if (error instanceof BackendError && error.status === 401) {
      return { error: text.invalidCredentials };
    }
    return { error: text.serviceUnavailable };
  }
}

export default function LoginPage() {
  const { locale } = useLoaderData<typeof loader>();
  const text = messages[locale].login;
  const actionData = useActionData<typeof action>();
  const navigation = useNavigation();
  const [showPassword, setShowPassword] = useState(false);
  const submitting = navigation.state === "submitting";

  return (
    <main className="relative grid min-h-screen place-items-center overflow-hidden px-4 py-10">
      <div className="absolute inset-x-0 top-0 h-1 bg-primary" />
      <section className="w-full max-w-sm rounded-lg bg-card p-7 shadow-[0_18px_55px_rgba(16,35,25,0.11)] sm:p-8">
        <div className="mb-8 flex items-start justify-between gap-3">
          <div className="flex items-center gap-3">
            <span className="grid size-11 place-items-center rounded-md bg-primary text-primary-foreground">
              <ShieldCheck aria-hidden="true" className="size-6" />
            </span>
            <div>
              <h1 className="text-xl font-semibold">AutoSSL</h1>
              <p className="mt-0.5 text-sm text-muted-foreground">{text.subtitle}</p>
            </div>
          </div>
          <LanguageSwitcher locale={locale} />
        </div>

        <Form method="post" className="space-y-5">
          <div className="space-y-2">
            <Label htmlFor="username">{text.username}</Label>
            <Input
              id="username"
              name="username"
              autoComplete="username"
              autoFocus
              required
              disabled={submitting}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">{text.password}</Label>
            <div className="relative">
              <Input
                id="password"
                name="password"
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                className="pr-11"
                required
                disabled={submitting}
              />
              <button
                type="button"
                className="absolute right-0 top-0 grid size-11 place-items-center rounded-sm text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                onClick={() => setShowPassword((value) => !value)}
                aria-label={showPassword ? text.hidePassword : text.showPassword}
              >
                {showPassword ? <EyeOff aria-hidden="true" /> : <Eye aria-hidden="true" />}
              </button>
            </div>
          </div>

          {actionData?.error ? (
            <div role="alert" className="flex gap-2 rounded-md bg-destructive/8 px-3 py-2.5 text-sm text-destructive">
              <KeyRound aria-hidden="true" className="mt-0.5 size-4 shrink-0" />
              <span>{actionData.error}</span>
            </div>
          ) : null}

          <Button type="submit" className="w-full" disabled={submitting}>
            {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : null}
            {submitting ? text.signingIn : text.signIn}
          </Button>
        </Form>
      </section>
    </main>
  );
}
