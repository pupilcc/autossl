export const messages = {
  en: {
    root: {
      pageLoadFailed: "The page could not be loaded. Please try again later.",
      notFound: "The requested page was not found.",
      problem: "Something went wrong",
      backToConsole: "Back to console",
    },
    login: {
      metaTitle: "Sign in · AutoSSL",
      subtitle: "Certificate console",
      username: "Username",
      password: "Password",
      showPassword: "Show password",
      hidePassword: "Hide password",
      missingCredentials: "Enter your username and password.",
      invalidCredentials: "The username or password is incorrect.",
      serviceUnavailable: "The service is unavailable. Please try again later.",
      signingIn: "Signing in",
      signIn: "Sign in",
    },
    certificates: {
      metaTitle: "Certificate console · AutoSSL",
      loadFailed: "The certificate list could not be loaded.",
      invalidDomain: "Enter a valid base domain, such as example.com.",
      created: (domain: string) => `Certificate for ${domain} was created.`,
      invalidCertificate: "The certificate identifier is invalid.",
      deleted: "Certificate deleted.",
      unknownAction: "Unknown action.",
      operationFailed: "The operation failed. Please try again later.",
      issueTitle: "01 / Issue certificate",
      dnsVerification: "DNS-01 validation",
      baseDomain: "Base domain",
      domainValidation: "Enter a base domain, such as example.com",
      domainHint: "Do not include a protocol or * prefix.",
      issuing: "Issuing",
      issue: "Issue certificate",
      coverage: "Coverage",
      rootAndWildcard: "Root domain + wildcard",
      keyAlgorithm: "Key algorithm",
      copySuccess: (label: string) => `${label} copied.`,
      copyFailed: "Copy failed. Allow clipboard access and try again.",
      copy: (label: string) => `Copy ${label}`,
      notConfigured: (label: string) => `${label} not configured`,
      delete: "Delete",
      deleteTitle: (domain: string) => `Delete ${domain}?`,
      deleteDescription:
        "The certificate record and its files will be deleted. This cannot be undone.",
      cancel: "Cancel",
      confirmDelete: "Delete certificate",
      brandLine: "CERTIFICATE CONTROL",
      apiOnline: "API ONLINE",
      logout: "Sign out",
      logoutShort: "Exit",
      infrastructure: "/ INFRASTRUCTURE",
      resources: "Certificate resources",
      description: "Issue, renew, and distribute root-domain and wildcard certificates.",
      domainCount: (count: number) => (count === 1 ? "domain" : "domains"),
      hosted: "hosted",
      listTitle: "02 / Certificate resources",
      total: "TOTAL",
      empty: "No certificate resources",
      validation: "Validation",
      algorithm: "Algorithm",
      certificateUrl: "certificate URL",
      privateKeyUrl: "private key URL",
    },
  },
  "zh-CN": {
    root: {
      pageLoadFailed: "页面暂时无法加载，请稍后重试。",
      notFound: "未找到请求的页面。",
      problem: "出现问题",
      backToConsole: "返回控制台",
    },
    login: {
      metaTitle: "登录 · AutoSSL",
      subtitle: "证书控制台",
      username: "用户名",
      password: "密码",
      showPassword: "显示密码",
      hidePassword: "隐藏密码",
      missingCredentials: "请输入用户名和密码。",
      invalidCredentials: "用户名或密码不正确。",
      serviceUnavailable: "暂时无法连接服务，请稍后重试。",
      signingIn: "正在登录",
      signIn: "登录",
    },
    certificates: {
      metaTitle: "证书控制台 · AutoSSL",
      loadFailed: "证书列表加载失败。",
      invalidDomain: "请输入有效的基础域名，例如 example.com。",
      created: (domain: string) => `${domain} 的证书已创建。`,
      invalidCertificate: "证书标识无效。",
      deleted: "证书已删除。",
      unknownAction: "未知操作。",
      operationFailed: "操作失败，请稍后重试。",
      issueTitle: "01 / 签发证书",
      dnsVerification: "DNS-01 验证",
      baseDomain: "基础域名",
      domainValidation: "请输入基础域名，例如 example.com",
      domainHint: "无需填写协议或 * 前缀。",
      issuing: "正在签发",
      issue: "签发证书",
      coverage: "覆盖范围",
      rootAndWildcard: "根域名 + 通配符",
      keyAlgorithm: "密钥算法",
      copySuccess: (label: string) => `${label}已复制。`,
      copyFailed: "复制失败，请允许浏览器访问剪贴板后重试。",
      copy: (label: string) => `复制${label}`,
      notConfigured: (label: string) => `${label}未配置`,
      delete: "删除",
      deleteTitle: (domain: string) => `删除 ${domain}？`,
      deleteDescription: "证书记录及对应文件将被删除，此操作无法撤销。",
      cancel: "取消",
      confirmDelete: "确认删除",
      brandLine: "证书控制",
      apiOnline: "API 在线",
      logout: "退出登录",
      logoutShort: "退出",
      infrastructure: "/ 基础设施",
      resources: "证书资源",
      description: "签发、续期并分发根域名与通配符证书。",
      domainCount: () => "个域名",
      hosted: "正在托管",
      listTitle: "02 / 证书资源",
      total: "共",
      empty: "暂无证书资源",
      validation: "验证",
      algorithm: "算法",
      certificateUrl: "证书地址",
      privateKeyUrl: "私钥地址",
    },
  },
} as const;

export type Locale = keyof typeof messages;

const LOCALE_COOKIE = "autossl_locale";

function parseLocale(value: string | null | undefined): Locale | null {
  return value === "en" || value === "zh-CN" ? value : null;
}

export function getRequestedLocale(request: Request) {
  return parseLocale(new URL(request.url).searchParams.get("lang"));
}

export function getLocale(request: Request): Locale {
  const requestedLocale = getRequestedLocale(request);
  if (requestedLocale) return requestedLocale;

  const cookieLocale = request.headers
    .get("cookie")
    ?.split(";")
    .map((part) => part.trim().split("=", 2))
    .find(([name]) => name === LOCALE_COOKIE)?.[1];
  const savedLocale = parseLocale(cookieLocale);
  if (savedLocale) return savedLocale;

  const language = request.headers
    .get("accept-language")
    ?.split(",", 1)[0]
    .trim()
    .toLowerCase();
  return language?.startsWith("zh") ? "zh-CN" : "en";
}

export function createLocaleCookie(request: Request, locale: Locale) {
  const protocol = request.headers.get("X-Forwarded-Proto") ?? new URL(request.url).protocol;
  const secure = protocol.replace(":", "") === "https" ? "; Secure" : "";
  return `${LOCALE_COOKIE}=${locale}; Path=/; SameSite=Lax; Max-Age=31536000${secure}`;
}
