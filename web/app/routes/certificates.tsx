import {
  Activity,
  Check,
  Clipboard,
  FileKey2,
  Globe2,
  KeyRound,
  LoaderCircle,
  LogOut,
  RefreshCw,
  ShieldCheck,
  Trash2,
} from "lucide-react";
import { useEffect, useId, useRef } from "react";
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
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
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
    <aside className="rounded-lg border bg-card p-5 lg:sticky lg:top-6">
      <div className="flex items-start gap-3">
        <span className="grid size-10 shrink-0 place-items-center rounded-md bg-primary/10 text-primary">
          <FileKey2 aria-hidden="true" className="size-5" />
        </span>
        <div>
          <h2 className="font-semibold">签发新证书</h2>
          <p className="mt-1 text-sm leading-5 text-muted-foreground">
            同时覆盖根域名和通配符域名。
          </p>
        </div>
      </div>
      <fetcher.Form
        ref={formRef}
        method="post"
        className="mt-6 space-y-4"
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
        <Button type="submit" className="w-full" disabled={submitting}>
          {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Check aria-hidden="true" />}
          {submitting ? "正在签发" : "签发证书"}
        </Button>
      </fetcher.Form>

      <dl className="mt-6 space-y-3 border-t pt-5 text-sm">
        <div className="flex items-center justify-between gap-4">
          <dt className="text-muted-foreground">验证方式</dt>
          <dd className="font-medium">DNS-01</dd>
        </div>
        <div className="flex items-center justify-between gap-4">
          <dt className="text-muted-foreground">密钥算法</dt>
          <dd className="font-medium">ECDSA P-256</dd>
        </div>
      </dl>
    </aside>
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

function CertificateLink({ label, value }: { label: string; value: string }) {
  const copy = async () => {
    try {
      await copyText(value);
      toast.success(`${label}已复制。`);
    } catch {
      toast.error("复制失败，请手动选择链接。");
    }
  };

  if (!value) return <span className="text-muted-foreground">未配置</span>;

  return (
    <div className="flex min-w-0 items-center gap-1 rounded-md bg-secondary/55 py-1 pl-3 pr-1">
      <code className="min-w-0 flex-1 truncate text-xs text-foreground" title={value}>
        {value}
      </code>
      <Button
        type="button"
        variant="ghost"
        size="icon"
        onClick={copy}
        aria-label={`复制${label}`}
        title={`复制${label}`}
      >
        <Clipboard aria-hidden="true" />
      </Button>
    </div>
  );
}

function DeleteCertificate({ code, domain }: { code: string; domain: string }) {
  const fetcher = useFetcher<ActionResult>();
  useActionToast(fetcher);
  const submitting = fetcher.state !== "idle";
  const formId = `delete-${code}-${useId().replaceAll(":", "")}`;

  return (
    <>
      <fetcher.Form id={formId} method="post">
        <input type="hidden" name="intent" value="delete" />
        <input type="hidden" name="code" value={code} />
      </fetcher.Form>
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
            <AlertDialogAction form={formId} type="submit" disabled={submitting}>
              {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Trash2 aria-hidden="true" />}
              确认删除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

export default function CertificatesPage() {
  const { certificates } = useLoaderData<typeof loader>();

  return (
    <div className="min-h-screen">
      <header className="border-b bg-card">
        <div className="mx-auto flex h-[72px] max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <span className="grid size-10 place-items-center rounded-md bg-primary text-primary-foreground">
              <ShieldCheck aria-hidden="true" className="size-5.5" />
            </span>
            <div>
              <p className="font-semibold leading-5">AutoSSL</p>
              <p className="hidden text-xs text-muted-foreground sm:block">证书运维控制台</p>
            </div>
          </div>
          <div className="flex items-center gap-2 sm:gap-4">
            <span className="hidden items-center gap-2 text-sm text-muted-foreground sm:flex">
              <span className="size-2 rounded-full bg-primary" aria-hidden="true" />
              服务已连接
            </span>
            <Form method="post">
              <input type="hidden" name="intent" value="logout" />
              <Button type="submit" variant="ghost" size="sm">
                <LogOut aria-hidden="true" />
                <span className="hidden sm:inline">退出登录</span>
                <span className="sm:hidden">退出</span>
              </Button>
            </Form>
          </div>
        </div>
      </header>

      <main>
        <section className="border-b bg-card">
          <div className="mx-auto max-w-7xl px-4 pb-0 pt-8 sm:px-6 sm:pt-10 lg:px-8">
            <div className="max-w-2xl">
              <p className="text-sm font-medium text-primary">证书管理</p>
              <h1 className="mt-2 text-2xl font-semibold sm:text-[28px]">证书控制台</h1>
              <p className="mt-2 text-sm leading-6 text-muted-foreground sm:text-base">
                集中签发 SSL 证书，并为部署脚本提供稳定的证书与私钥地址。
              </p>
            </div>

            <div className="mt-8 grid divide-y border-t sm:grid-cols-3 sm:divide-x sm:divide-y-0">
              <div className="flex items-center gap-3 py-4 sm:pr-6">
                <Activity aria-hidden="true" className="size-5 text-primary" />
                <div>
                  <p className="text-xs text-muted-foreground">托管证书</p>
                  <p className="mt-0.5 text-sm font-semibold tabular-nums">{certificates.length} 个域名</p>
                </div>
              </div>
              <div className="flex items-center gap-3 py-4 sm:px-6">
                <Globe2 aria-hidden="true" className="size-5 text-primary" />
                <div>
                  <p className="text-xs text-muted-foreground">签发范围</p>
                  <p className="mt-0.5 text-sm font-semibold">根域名 + 通配符</p>
                </div>
              </div>
              <div className="flex items-center gap-3 py-4 sm:pl-6">
                <RefreshCw aria-hidden="true" className="size-5 text-primary" />
                <div>
                  <p className="text-xs text-muted-foreground">续期策略</p>
                  <p className="mt-0.5 text-sm font-semibold">每日自动检查</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div className="mx-auto grid max-w-7xl gap-6 px-4 py-6 sm:px-6 sm:py-8 lg:px-8 xl:grid-cols-[minmax(0,1fr)_19rem] xl:gap-8">
          <div className="xl:order-2">
            <AddCertificateForm />
          </div>

          <section aria-labelledby="certificate-list-title" className="min-w-0 xl:order-1">
            <div className="mb-4 flex items-end justify-between gap-4">
              <div>
                <div className="flex items-center gap-2">
                  <h2 id="certificate-list-title" className="text-lg font-semibold">证书与分发地址</h2>
                  <Badge variant="secondary">{certificates.length}</Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">复制地址后即可用于自动化部署。</p>
              </div>
            </div>

            {certificates.length === 0 ? (
              <div className="grid min-h-72 place-items-center rounded-lg border border-dashed bg-card px-6 text-center">
                <div>
                  <span className="mx-auto grid size-12 place-items-center rounded-md bg-secondary text-primary">
                    <KeyRound aria-hidden="true" className="size-6" />
                  </span>
                  <h3 className="mt-4 font-semibold">暂无域名证书</h3>
                  <p className="mt-1 text-sm text-muted-foreground">新签发的证书会显示在这里。</p>
                </div>
              </div>
            ) : (
              <>
                <div className="hidden overflow-hidden rounded-lg border bg-card xl:block">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-[22%]">覆盖域名</TableHead>
                        <TableHead className="w-[31%]">证书地址</TableHead>
                        <TableHead className="w-[31%]">私钥地址 · 敏感</TableHead>
                        <TableHead className="w-[16%] text-right">
                          <span className="sr-only">操作</span>
                        </TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {certificates.map((certificate) => (
                        <TableRow key={certificate.code}>
                          <TableCell>
                            <p className="font-medium">{certificate.domain}</p>
                            <p className="mt-1 text-xs text-muted-foreground">*.{certificate.domain}</p>
                          </TableCell>
                          <TableCell>
                            <CertificateLink
                              label={`${certificate.domain} 的证书地址`}
                              value={certificate.cert}
                            />
                          </TableCell>
                          <TableCell>
                            <CertificateLink
                              label={`${certificate.domain} 的私钥地址`}
                              value={certificate.key}
                            />
                          </TableCell>
                          <TableCell className="text-right">
                            <DeleteCertificate code={certificate.code} domain={certificate.domain} />
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>

                <div className="space-y-3 xl:hidden">
                  {certificates.map((certificate) => (
                    <article
                      key={certificate.code}
                      className="overflow-hidden rounded-lg border bg-card"
                    >
                      <div className="flex items-start justify-between gap-4 px-4 py-4">
                        <div className="min-w-0">
                          <h3 className="truncate font-semibold">{certificate.domain}</h3>
                          <p className="mt-1 truncate text-xs text-muted-foreground">*.{certificate.domain}</p>
                        </div>
                        <DeleteCertificate code={certificate.code} domain={certificate.domain} />
                      </div>
                      <dl className="space-y-4 border-t bg-secondary/20 px-4 py-4">
                        <div className="min-w-0 space-y-2">
                          <dt className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                            <ShieldCheck aria-hidden="true" className="size-3.5" />
                            证书地址
                          </dt>
                          <dd>
                            <CertificateLink
                              label={`${certificate.domain} 的证书地址`}
                              value={certificate.cert}
                            />
                          </dd>
                        </div>
                        <div className="min-w-0 space-y-2">
                          <dt className="flex items-center gap-1.5 text-xs font-medium text-muted-foreground">
                            <KeyRound aria-hidden="true" className="size-3.5" />
                            私钥地址 · 敏感
                          </dt>
                          <dd>
                            <CertificateLink
                              label={`${certificate.domain} 的私钥地址`}
                              value={certificate.key}
                            />
                          </dd>
                        </div>
                      </dl>
                    </article>
                  ))}
                </div>
              </>
            )}
          </section>
        </div>
      </main>
    </div>
  );
}
