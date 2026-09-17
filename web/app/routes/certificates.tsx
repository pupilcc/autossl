import {
  Check,
  Clipboard,
  KeyRound,
  LoaderCircle,
  LogOut,
  ShieldCheck,
  Trash2,
} from "lucide-react";
import { useEffect, useRef } from "react";
import {
  Form,
  redirect,
  useFetcher,
  useLoaderData,
  useRevalidator,
  type ActionFunctionArgs,
  type LoaderFunctionArgs,
  type MetaFunction,
} from "react-router";
import { toast } from "sonner";

import { LanguageSwitcher } from "@/components/language-switcher";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  BackendError,
  createCertificate,
  deleteCertificate,
  listCertificates,
} from "@/lib/backend.server";
import { getLocale, messages, type Locale } from "@/lib/i18n";
import {
  assertSameOrigin,
  destroySessionCookie,
  requireSession,
} from "@/lib/session.server";

export const meta: MetaFunction<typeof loader> = ({ data }) => [
  { title: messages[data?.locale ?? "en"].certificates.metaTitle },
];

export type ActionResult = { ok: boolean; message: string };
type ActionFetcher = ReturnType<typeof useFetcher<ActionResult>>;

export async function loader({ request }: LoaderFunctionArgs) {
  const locale = getLocale(request);
  const text = messages[locale].certificates;
  const session = requireSession(request);
  try {
    return { certificates: await listCertificates(session.token), locale };
  } catch (error) {
    if (error instanceof BackendError && error.status === 401) {
      return redirect("/login", {
        headers: { "Set-Cookie": destroySessionCookie(request) },
      });
    }
    throw new Response(text.loadFailed, { status: 502 });
  }
}

export async function action({ request }: ActionFunctionArgs): Promise<ActionResult | Response> {
  assertSameOrigin(request);
  const locale = getLocale(request);
  const text = messages[locale].certificates;
  const session = requireSession(request);
  const formData = await request.formData();
  const intent = String(formData.get("intent") ?? "");

  if (intent === "logout") {
    return redirect("/login", {
      headers: { "Set-Cookie": destroySessionCookie(request) },
    });
  }

  try {
    if (intent === "create") {
      const domain = String(formData.get("domain") ?? "")
        .trim()
        .toLowerCase()
        .replace(/\.$/, "");
      const validDomain = /^(?=.{1,253}$)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/.test(domain);
      if (!validDomain) return { ok: false, message: text.invalidDomain };

      await createCertificate(session.token, domain);
      return { ok: true, message: text.created(domain) };
    }

    if (intent === "delete") {
      const code = String(formData.get("code") ?? "");
      if (!/^[A-Za-z0-9_-]+$/.test(code)) {
        return { ok: false, message: text.invalidCertificate };
      }
      await deleteCertificate(session.token, code);
      return { ok: true, message: text.deleted };
    }

    return { ok: false, message: text.unknownAction };
  } catch (error) {
    if (error instanceof BackendError && error.status === 401) {
      return redirect("/login", {
        headers: { "Set-Cookie": destroySessionCookie(request) },
      });
    }
    return {
      ok: false,
      message:
        error instanceof BackendError && locale === "en"
          ? error.message
          : text.operationFailed,
    };
  }
}

function useActionToast(fetcher: ActionFetcher, onSuccess?: () => void) {
  const handled = useRef<ActionResult | undefined>(undefined);

  useEffect(() => {
    if (fetcher.state !== "idle" || !fetcher.data || fetcher.data === handled.current) return;
    handled.current = fetcher.data;
    if (fetcher.data.ok) {
      toast.success(fetcher.data.message);
      onSuccess?.();
    } else {
      toast.error(fetcher.data.message);
    }
  }, [fetcher.data, fetcher.state, onSuccess]);
}

function AddCertificateForm({ locale }: { locale: Locale }) {
  const text = messages[locale].certificates;
  const fetcher = useFetcher<ActionResult>();
  const formRef = useRef<HTMLFormElement>(null);
  useActionToast(fetcher, () => formRef.current?.reset());
  const submitting = fetcher.state !== "idle";

  return (
    <section className="border-2 border-foreground bg-card" aria-labelledby="issue-title">
      <div className="flex min-h-11 items-center justify-between gap-4 bg-foreground px-4 py-2 text-card">
        <h2 id="issue-title" className="font-mono text-xs font-semibold uppercase">
          {text.issueTitle}
        </h2>
        <span className="flex items-center gap-2 text-xs text-card/70">
          <ShieldCheck aria-hidden="true" className="size-4 text-primary" />
          {text.dnsVerification}
        </span>
      </div>
      <fetcher.Form
        ref={formRef}
        method="post"
        className="grid gap-x-4 gap-y-2 p-4 sm:p-5 lg:grid-cols-[minmax(0,1fr)_10rem]"
      >
        <input type="hidden" name="intent" value="create" />
        <Label htmlFor="domain" className="lg:col-start-1 lg:row-start-1">{text.baseDomain}</Label>
        <Input
          id="domain"
          name="domain"
          inputMode="url"
          autoComplete="off"
          placeholder="example.com"
          className="h-12 rounded-none border-foreground shadow-none lg:col-start-1 lg:row-start-2"
          pattern="(?=.{1,253}$)([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?"
          title={text.domainValidation}
          aria-describedby="domain-hint"
          required
          disabled={submitting}
        />
        <p id="domain-hint" className="text-xs leading-5 text-muted-foreground lg:col-start-1 lg:row-start-3">
          {text.domainHint}
        </p>
        <Button
          type="submit"
          className="mt-2 h-12 w-full rounded-none lg:col-start-2 lg:row-start-2 lg:mt-0"
          disabled={submitting}
        >
          {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Check aria-hidden="true" />}
          {submitting ? text.issuing : text.issue}
        </Button>
      </fetcher.Form>

      <dl className="grid border-t bg-secondary/45 text-xs sm:grid-cols-2 sm:divide-x">
        <div className="flex items-center justify-between gap-4 px-4 py-3 sm:px-5">
          <dt className="text-muted-foreground">{text.coverage}</dt>
          <dd className="font-medium">{text.rootAndWildcard}</dd>
        </div>
        <div className="flex items-center justify-between gap-4 border-t px-4 py-3 sm:border-t-0 sm:px-5">
          <dt className="text-muted-foreground">{text.keyAlgorithm}</dt>
          <dd className="font-medium">ECDSA P-256</dd>
        </div>
      </dl>
    </section>
  );
}

async function copyText(value: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value);
    return;
  }
  const textarea = document.createElement("textarea");
  textarea.value = value;
  textarea.style.position = "fixed";
  textarea.style.opacity = "0";
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand("copy");
  textarea.remove();
}

function CopyCertificate({
  label,
  value,
  locale,
  sensitive = false,
}: {
  label: string;
  value: string;
  locale: Locale;
  sensitive?: boolean;
}) {
  const text = messages[locale].certificates;
  const copy = async () => {
    try {
      await copyText(value);
      toast.success(text.copySuccess(label));
    } catch {
      toast.error(text.copyFailed);
    }
  };

  return (
    <Button
      type="button"
      variant="outline"
      onClick={copy}
      disabled={!value}
      className={
        sensitive
          ? "border-accent bg-accent text-accent-foreground hover:bg-accent/80"
          : undefined
      }
      aria-label={value ? text.copy(label) : text.notConfigured(label)}
    >
      <Clipboard aria-hidden="true" />
      {value ? text.copy(label) : text.notConfigured(label)}
    </Button>
  );
}

function DeleteCertificate({
  code,
  domain,
  fetcher,
  locale,
}: {
  code: string;
  domain: string;
  fetcher: ActionFetcher;
  locale: Locale;
}) {
  const text = messages[locale].certificates;
  const submitting = fetcher.state !== "idle";

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button type="button" variant="ghost" size="sm" className="text-destructive hover:bg-destructive/8 hover:text-destructive">
          <Trash2 aria-hidden="true" />
          {text.delete}
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{text.deleteTitle(domain)}</AlertDialogTitle>
          <AlertDialogDescription>
            {text.deleteDescription}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={submitting}>{text.cancel}</AlertDialogCancel>
          <AlertDialogAction
            type="button"
            disabled={submitting}
            onClick={() =>
              fetcher.submit(
                { intent: "delete", code },
                { method: "post", defaultShouldRevalidate: false },
              )
            }
          >
            {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Trash2 aria-hidden="true" />}
            {text.confirmDelete}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export default function CertificatesPage() {
  const { certificates, locale } = useLoaderData<typeof loader>();
  const text = messages[locale].certificates;
  const deleteFetcher = useFetcher<ActionResult>();
  const { revalidate } = useRevalidator();
  useActionToast(deleteFetcher, revalidate);

  return (
    <div className="min-h-screen">
      <header className="border-b-4 border-primary bg-foreground text-card">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <span className="grid size-9 place-items-center bg-primary text-primary-foreground">
              <ShieldCheck aria-hidden="true" className="size-5" />
            </span>
            <div>
              <p className="font-semibold leading-5">AutoSSL</p>
              <p className="hidden font-mono text-[11px] text-card/55 sm:block">{text.brandLine}</p>
            </div>
          </div>
          <div className="flex items-center gap-2 sm:gap-4">
            <span className="hidden items-center gap-2 font-mono text-xs text-card/70 sm:flex">
              <span className="size-2 bg-primary" aria-hidden="true" />
              {text.apiOnline}
            </span>
            <LanguageSwitcher locale={locale} inverted />
            <Form method="post">
              <input type="hidden" name="intent" value="logout" />
              <Button
                type="submit"
                variant="ghost"
                size="sm"
                className="text-card hover:bg-card/10 hover:text-card"
              >
                <LogOut aria-hidden="true" />
                <span className="hidden sm:inline">{text.logout}</span>
                <span className="sm:hidden">{text.logoutShort}</span>
              </Button>
            </Form>
          </div>
        </div>
      </header>

      <main>
        <section className="border-b-2 border-foreground bg-card">
          <div className="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-8 sm:px-6 sm:py-10 md:flex-row md:items-end md:justify-between lg:px-8">
            <div>
              <p className="font-mono text-xs font-semibold text-primary">{text.infrastructure}</p>
              <h1 className="mt-3 text-3xl font-semibold sm:text-4xl">{text.resources}</h1>
              <p className="mt-2 max-w-xl text-sm leading-6 text-muted-foreground">
                {text.description}
              </p>
            </div>
            <div className="flex items-end gap-3 border-l-4 border-primary pl-4 md:text-right">
              <strong className="text-4xl font-semibold tabular-nums sm:text-5xl">{certificates.length}</strong>
              <span className="pb-1 text-xs leading-5 text-muted-foreground">
                {text.domainCount(certificates.length)}<br />{text.hosted}
              </span>
            </div>
          </div>
        </section>

        <div className="mx-auto max-w-7xl space-y-10 px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
          <AddCertificateForm locale={locale} />

          <section aria-labelledby="certificate-list-title">
            <div className="mb-4 flex items-center justify-between gap-4 border-b-2 border-foreground pb-3">
              <h2 id="certificate-list-title" className="font-mono text-xs font-semibold uppercase">
                {text.listTitle}
              </h2>
              <span className="font-mono text-xs text-muted-foreground">{text.total} {certificates.length}</span>
            </div>

            {certificates.length === 0 ? (
              <div className="grid min-h-64 place-items-center border-x border-b bg-card px-6 text-center">
                <div>
                  <span className="mx-auto grid size-12 place-items-center border-2 border-foreground text-primary">
                    <KeyRound aria-hidden="true" className="size-6" />
                  </span>
                  <h3 className="mt-4 font-semibold">{text.empty}</h3>
                </div>
              </div>
            ) : (
              <div>
                {certificates.map((certificate) => (
                  <article
                    key={certificate.code}
                    className="grid gap-5 border-x border-t bg-card p-4 transition-colors hover:bg-secondary/25 last:border-b sm:p-5 lg:grid-cols-[minmax(0,1fr)_12rem_auto] lg:items-center"
                  >
                    <div className="grid min-w-0 grid-cols-[auto_minmax(0,1fr)] gap-x-3">
                      <span className="size-2.5 self-center bg-primary" aria-hidden="true" />
                      <h3 className="truncate text-lg font-semibold">{certificate.domain}</h3>
                      <p className="col-start-2 mt-1 truncate font-mono text-xs text-muted-foreground">
                        {(certificate.dnsNames?.length
                          ? certificate.dnsNames
                          : [certificate.domain, `*.${certificate.domain}`]
                        ).join(" + ")}
                      </p>
                    </div>

                    <dl className="grid grid-cols-2 gap-3 border-y py-3 text-xs sm:max-w-sm lg:block lg:border-y-0 lg:border-l lg:py-0 lg:pl-5">
                      <div>
                        <dt className="text-muted-foreground">{text.validation}</dt>
                        <dd className="mt-1 font-medium">DNS-01</dd>
                      </div>
                      <div className="lg:mt-3">
                        <dt className="text-muted-foreground">{text.algorithm}</dt>
                        <dd className="mt-1 font-medium">ECDSA P-256</dd>
                      </div>
                    </dl>

                    <div className="flex flex-wrap gap-2 lg:justify-end">
                      <CopyCertificate label={text.certificateUrl} value={certificate.cert} locale={locale} />
                      <CopyCertificate label={text.privateKeyUrl} value={certificate.key} locale={locale} sensitive />
                      <DeleteCertificate
                        code={certificate.code}
                        domain={certificate.domain}
                        fetcher={deleteFetcher}
                        locale={locale}
                      />
                    </div>
                  </article>
                ))}
              </div>
            )}
          </section>
        </div>
      </main>
    </div>
  );
}
