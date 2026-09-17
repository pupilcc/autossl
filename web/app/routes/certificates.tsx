import {
  Check,
  Clipboard,
  KeyRound,
  LoaderCircle,
  LogOut,
  Plus,
  ShieldCheck,
  Trash2,
  X,
} from "lucide-react";
import { useEffect, useId, useRef, useState } from "react";
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

export const meta: MetaFunction = () => [{ title: "证书管理 · AutoSSL" }];

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

function AddCertificateForm({ onClose }: { onClose: () => void }) {
  const fetcher = useFetcher<ActionResult>();
  const formRef = useRef<HTMLFormElement>(null);
  const onSuccess = () => {
    formRef.current?.reset();
    onClose();
  };
  useActionToast(fetcher, onSuccess);
  const submitting = fetcher.state !== "idle";

  return (
    <section className="border-y bg-card">
      <fetcher.Form
        ref={formRef}
        method="post"
        className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-5 sm:px-6 lg:flex-row lg:items-end lg:px-8"
      >
        <input type="hidden" name="intent" value="create" />
        <div className="min-w-0 flex-1 space-y-2">
          <Label htmlFor="domain">基础域名</Label>
          <Input
            id="domain"
            name="domain"
            inputMode="url"
            autoComplete="off"
            placeholder="example.com"
            pattern="(?=.{1,253}$)([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?"
            title="请输入基础域名，例如 example.com"
            required
            disabled={submitting}
          />
        </div>
        <div className="flex gap-2">
          <Button type="submit" disabled={submitting}>
            {submitting ? <LoaderCircle className="animate-spin" aria-hidden="true" /> : <Check aria-hidden="true" />}
            {submitting ? "正在签发" : "确认添加"}
          </Button>
          <Button type="button" variant="outline" onClick={onClose} disabled={submitting}>
            <X aria-hidden="true" />
            取消
          </Button>
        </div>
      </fetcher.Form>
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
    <div className="flex min-w-0 items-center gap-2">
      <code className="min-w-0 flex-1 truncate text-xs text-foreground" title={value}>
        {value}
      </code>
      <Button type="button" variant="ghost" size="sm" onClick={copy} aria-label={`复制${label}`}>
        <Clipboard aria-hidden="true" />
        复制
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
  const [adding, setAdding] = useState(false);

  return (
    <div className="min-h-screen">
      <header className="border-b bg-card">
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <span className="grid size-9 place-items-center rounded-md bg-primary text-primary-foreground">
              <ShieldCheck aria-hidden="true" className="size-5" />
            </span>
            <span className="font-semibold">AutoSSL</span>
          </div>
          <Form method="post">
            <input type="hidden" name="intent" value="logout" />
            <Button type="submit" variant="ghost" size="sm">
              <LogOut aria-hidden="true" />
              退出
            </Button>
          </Form>
        </div>
      </header>

      <main>
        <div className="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-8 sm:flex-row sm:items-center sm:justify-between sm:px-6 lg:px-8">
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-semibold">域名证书</h1>
            <Badge variant="secondary">{certificates.length}</Badge>
          </div>
          <Button type="button" onClick={() => setAdding((value) => !value)} aria-expanded={adding}>
            {adding ? <X aria-hidden="true" /> : <Plus aria-hidden="true" />}
            {adding ? "收起" : "添加域名"}
          </Button>
        </div>

        {adding ? <AddCertificateForm onClose={() => setAdding(false)} /> : null}

        <section className="mx-auto max-w-7xl px-4 pb-12 pt-6 sm:px-6 lg:px-8">
          {certificates.length === 0 ? (
            <div className="grid min-h-72 place-items-center rounded-lg border border-dashed bg-card px-6 text-center">
              <div>
                <span className="mx-auto grid size-12 place-items-center rounded-full bg-secondary text-primary">
                  <KeyRound aria-hidden="true" className="size-6" />
                </span>
                <h2 className="mt-4 text-base font-semibold">暂无域名证书</h2>
                <Button type="button" className="mt-5" onClick={() => setAdding(true)}>
                  <Plus aria-hidden="true" />
                  添加第一个域名
                </Button>
              </div>
            </div>
          ) : (
            <>
              <div className="hidden overflow-hidden rounded-lg bg-card shadow-[0_8px_30px_rgba(16,35,25,0.07)] lg:block">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead className="w-[18%]">域名</TableHead>
                      <TableHead className="w-[33%]">证书链接</TableHead>
                      <TableHead className="w-[33%]">私钥链接</TableHead>
                      <TableHead className="w-[16%] text-right">操作</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {certificates.map((certificate) => (
                      <TableRow key={certificate.code}>
                        <TableCell className="font-medium">{certificate.domain}</TableCell>
                        <TableCell>
                          <CertificateLink label="证书链接" value={certificate.cert} />
                        </TableCell>
                        <TableCell>
                          <CertificateLink label="私钥链接" value={certificate.key} />
                        </TableCell>
                        <TableCell className="text-right">
                          <DeleteCertificate code={certificate.code} domain={certificate.domain} />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              <div className="space-y-3 lg:hidden">
                {certificates.map((certificate) => (
                  <article
                    key={certificate.code}
                    className="overflow-hidden rounded-lg bg-card shadow-[0_8px_30px_rgba(16,35,25,0.07)]"
                  >
                    <h2 className="px-4 py-4 font-semibold">{certificate.domain}</h2>
                    <dl className="space-y-4 border-y px-4 py-4">
                      <div className="min-w-0 space-y-1">
                        <dt className="text-xs font-medium text-muted-foreground">证书链接</dt>
                        <dd>
                          <CertificateLink label="证书链接" value={certificate.cert} />
                        </dd>
                      </div>
                      <div className="min-w-0 space-y-1">
                        <dt className="text-xs font-medium text-muted-foreground">私钥链接</dt>
                        <dd>
                          <CertificateLink label="私钥链接" value={certificate.key} />
                        </dd>
                      </div>
                    </dl>
                    <div className="flex justify-end px-3 py-2">
                      <DeleteCertificate code={certificate.code} domain={certificate.domain} />
                    </div>
                  </article>
                ))}
              </div>
            </>
          )}
        </section>
      </main>
    </div>
  );
}
