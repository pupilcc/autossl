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
  type ActionFunctionArgs,
  type LoaderFunctionArgs,
  type MetaFunction,
} from "react-router";
import { toast } from "sonner";

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
import {
  assertSameOrigin,
  destroySessionCookie,
  requireSession,
} from "@/lib/session.server";

export const meta: MetaFunction = () => [{ title: "证书控制台 · AutoSSL" }];

export type ActionResult = { ok: boolean; message: string };

export async function loader({ request }: LoaderFunctionArgs) {
  const session = requireSession(request);
  try {
    return { certificates: await listCertificates(session.token) };
  } catch (error) {
    if (error instanceof BackendError && error.status === 401) {
      return redirect("/login", {
        headers: { "Set-Cookie": destroySessionCookie(request) },
      });
    }
    throw new Response("证书列表加载失败。", { status: 502 });
  }
}

export async function action({ request }: ActionFunctionArgs): Promise<ActionResult | Response> {
  assertSameOrigin(request);
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
      if (!validDomain) return { ok: false, message: "请输入有效的基础域名，例如 example.com。" };

      await createCertificate(session.token, domain);
      return { ok: true, message: `${domain} 的证书已创建。` };
    }

    if (intent === "delete") {
      const code = String(formData.get("code") ?? "");
      if (!/^[A-Za-z0-9_-]+$/.test(code)) {
        return { ok: false, message: "证书标识无效。" };
      }
      await deleteCertificate(session.token, code);
      return { ok: true, message: "证书已删除。" };
    }

    return { ok: false, message: "未知操作。" };
  } catch (error) {
    if (error instanceof BackendError && error.status === 401) {
      return redirect("/login", {
        headers: { "Set-Cookie": destroySessionCookie(request) },
      });
    }
    return {
      ok: false,
      message: error instanceof BackendError ? error.message : "操作失败，请稍后重试。",
    };
  }
}

function useActionToast(fetcher: ReturnType<typeof useFetcher<ActionResult>>, onSuccess?: () => void) {
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

function AddCertificateForm() {
  const fetcher = useFetcher<ActionResult>();
  const formRef = useRef<HTMLFormElement>(null);
  useActionToast(fetcher, () => formRef.current?.reset());
  const submitting = fetcher.state !== "idle";

  return (
    <section className="border-2 border-foreground bg-card" aria-labelledby="issue-title">
      <div className="flex min-h-11 items-center justify-between gap-4 bg-foreground px-4 py-2 text-card">
        <h2 id="issue-title" className="font-mono text-xs font-semibold uppercase">
          01 / 签发证书
        </h2>
        <span className="flex items-center gap-2 text-xs text-card/70">
          <ShieldCheck aria-hidden="true" className="size-4 text-primary" />
          DNS-01 验证
        </span>
      </div>
      <fetcher.Form
        ref={formRef}
        method="post"
        className="grid gap-4 p-4 sm:p-5 lg:grid-cols-[minmax(0,1fr)_10rem] lg:items-end"
      >
        <input type="hidden" name="intent" value="create" />
        <div className="space-y-2">
          <Label htmlFor="domain">基础域名</Label>
          <Input
            id="domain"
            name="domain"
            inputMode="url"
            autoComplete="off"
            placeholder="example.com"
            className="h-12 rounded-none border-foreground shadow-none"
            pattern="(?=.{1,253}$)([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?"
            title="请输入基础域名，例如 example.com"
            aria-describedby="domain-hint"
            required
            disabled={submitting}
          />
          <p id="domain-hint" className="text-xs leading-5 text-muted-foreground">
            无需填写协议或 * 前缀。
          </p>
        </div>
        <Button type="submit" className="h-12 w-full rounded-none" disabled={submitting}>
          {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Check aria-hidden="true" />}
          {submitting ? "正在签发" : "签发证书"}
        </Button>
      </fetcher.Form>

      <dl className="grid border-t bg-secondary/45 text-xs sm:grid-cols-2 sm:divide-x">
        <div className="flex items-center justify-between gap-4 px-4 py-3 sm:px-5">
          <dt className="text-muted-foreground">覆盖范围</dt>
          <dd className="font-medium">根域名 + 通配符</dd>
        </div>
        <div className="flex items-center justify-between gap-4 border-t px-4 py-3 sm:border-t-0 sm:px-5">
          <dt className="text-muted-foreground">密钥算法</dt>
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
  sensitive = false,
}: {
  label: string;
  value: string;
  sensitive?: boolean;
}) {
  const copy = async () => {
    try {
      await copyText(value);
      toast.success(`${label}已复制。`);
    } catch {
      toast.error("复制失败，请允许浏览器访问剪贴板后重试。");
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
      aria-label={value ? `复制${label}` : `${label}未配置`}
    >
      <Clipboard aria-hidden="true" />
      {value ? `复制${label}` : `${label}未配置`}
    </Button>
  );
}

function DeleteCertificate({ code, domain }: { code: string; domain: string }) {
  const fetcher = useFetcher<ActionResult>();
  useActionToast(fetcher);
  const submitting = fetcher.state !== "idle";

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button type="button" variant="ghost" size="sm" className="text-destructive hover:bg-destructive/8 hover:text-destructive">
          <Trash2 aria-hidden="true" />
          删除
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>删除 {domain}？</AlertDialogTitle>
          <AlertDialogDescription>
            证书记录及对应文件将被删除，此操作无法撤销。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={submitting}>取消</AlertDialogCancel>
          <AlertDialogAction
            type="button"
            disabled={submitting}
            onClick={() => fetcher.submit({ intent: "delete", code }, { method: "post" })}
          >
            {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Trash2 aria-hidden="true" />}
            确认删除
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

export default function CertificatesPage() {
  const { certificates } = useLoaderData<typeof loader>();

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
              <p className="hidden font-mono text-[11px] text-card/55 sm:block">CERTIFICATE CONTROL</p>
            </div>
          </div>
          <div className="flex items-center gap-2 sm:gap-4">
            <span className="hidden items-center gap-2 font-mono text-xs text-card/70 sm:flex">
              <span className="size-2 bg-primary" aria-hidden="true" />
              API ONLINE
            </span>
            <Form method="post">
              <input type="hidden" name="intent" value="logout" />
              <Button
                type="submit"
                variant="ghost"
                size="sm"
                className="text-card hover:bg-card/10 hover:text-card"
              >
                <LogOut aria-hidden="true" />
                <span className="hidden sm:inline">退出登录</span>
                <span className="sm:hidden">退出</span>
              </Button>
            </Form>
          </div>
        </div>
      </header>

      <main>
        <section className="border-b-2 border-foreground bg-card">
          <div className="mx-auto flex max-w-7xl flex-col gap-6 px-4 py-8 sm:px-6 sm:py-10 md:flex-row md:items-end md:justify-between lg:px-8">
            <div>
              <p className="font-mono text-xs font-semibold text-primary">/ INFRASTRUCTURE</p>
              <h1 className="mt-3 text-3xl font-semibold sm:text-4xl">证书资源</h1>
              <p className="mt-2 max-w-xl text-sm leading-6 text-muted-foreground">
                签发、续期并分发根域名与通配符证书。
              </p>
            </div>
            <div className="flex items-end gap-3 border-l-4 border-primary pl-4 md:text-right">
              <strong className="text-4xl font-semibold tabular-nums sm:text-5xl">{certificates.length}</strong>
              <span className="pb-1 text-xs leading-5 text-muted-foreground">
                个域名<br />正在托管
              </span>
            </div>
          </div>
        </section>

        <div className="mx-auto max-w-7xl space-y-10 px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
          <AddCertificateForm />

          <section aria-labelledby="certificate-list-title">
            <div className="mb-4 flex items-center justify-between gap-4 border-b-2 border-foreground pb-3">
              <h2 id="certificate-list-title" className="font-mono text-xs font-semibold uppercase">
                02 / 证书资源
              </h2>
              <span className="font-mono text-xs text-muted-foreground">TOTAL {certificates.length}</span>
            </div>

            {certificates.length === 0 ? (
              <div className="grid min-h-64 place-items-center border-x border-b bg-card px-6 text-center">
                <div>
                  <span className="mx-auto grid size-12 place-items-center border-2 border-foreground text-primary">
                    <KeyRound aria-hidden="true" className="size-6" />
                  </span>
                  <h3 className="mt-4 font-semibold">暂无证书资源</h3>
                </div>
              </div>
            ) : (
              <div>
                {certificates.map((certificate) => (
                  <article
                    key={certificate.code}
                    className="grid gap-5 border-x border-t bg-card p-4 transition-colors hover:bg-secondary/25 last:border-b sm:p-5 lg:grid-cols-[minmax(0,1fr)_12rem_auto] lg:items-center"
                  >
                    <div className="flex min-w-0 items-start gap-3">
                      <span className="mt-1 size-2.5 shrink-0 bg-primary" aria-hidden="true" />
                      <div className="min-w-0">
                        <h3 className="truncate text-lg font-semibold">{certificate.domain}</h3>
                        <p className="mt-1 truncate font-mono text-xs text-muted-foreground">
                          {(certificate.dnsNames?.length
                            ? certificate.dnsNames
                            : [certificate.domain, `*.${certificate.domain}`]
                          ).join(" + ")}
                        </p>
                      </div>
                    </div>

                    <dl className="grid grid-cols-2 gap-3 border-y py-3 text-xs sm:max-w-sm lg:block lg:border-y-0 lg:border-l lg:py-0 lg:pl-5">
                      <div>
                        <dt className="text-muted-foreground">验证</dt>
                        <dd className="mt-1 font-medium">DNS-01</dd>
                      </div>
                      <div className="lg:mt-3">
                        <dt className="text-muted-foreground">算法</dt>
                        <dd className="mt-1 font-medium">ECDSA P-256</dd>
                      </div>
                    </dl>

                    <div className="flex flex-wrap gap-2 lg:justify-end">
                      <CopyCertificate label="证书地址" value={certificate.cert} />
                      <CopyCertificate label="私钥地址" value={certificate.key} sensitive />
                      <DeleteCertificate code={certificate.code} domain={certificate.domain} />
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
