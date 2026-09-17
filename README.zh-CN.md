# AutoSSL

[English](README.md)

**其他域名只需添加一条 CNAME，AutoSSL 就能自动申请和续期证书。**

先把一个域名（例如 `alias.com`）托管在 Cloudflare，并将它的 API Token 配置给 AutoSSL。其他需要证书的域名，即使使用不同的 DNS 服务商，也只需添加一条 CNAME，指向 `alias.com` 下的对应验证记录。AutoSSL 随后通过 [acme.sh DNS 别名模式](https://github.com/acmesh-official/acme.sh/wiki/DNS-alias-mode)完成证书签发和续期，不再需要这些域名各自的 DNS API 凭据。

证书仍然签发给目标域名，而不是 `alias.com`。AutoSSL 会为目标基础域名和通配符域名申请一张 ECDSA P-256 证书，再通过固定的证书与私钥 URL 提供给部署脚本。

## AutoSSL 解决什么问题

| 直接使用目标域名的 DNS API | 使用 AutoSSL |
| --- | --- |
| 保存每个目标 DNS 区域的 API 凭据。 | 只保存验证区域的一份 Cloudflare Token。 |
| 依赖每家目标 DNS 服务商提供 API。 | 目标 DNS 服务商只需支持 CNAME。 |
| 每台服务器分别签发、续期和部署证书。 | 集中签发和续期，目标服务器只需下载文件。 |

目标域名无需使用 Cloudflare，无需向 AutoSSL 提供自身 DNS API，也无需更换 NS。CNAME 只需配置一次，后续自动续期会继续使用它。

## 工作原理

假设 Cloudflare 托管的统一验证域名是 `alias.com`，现在需要为 `example.com` 签发证书。

1. 在 AutoSSL 中设置 `ACME_ALIAS=alias.com` 和 `ACME_ALIAS_PER_DOMAIN=true`。
2. 在 `example.com` 的 DNS 区域添加：

   ```dns
   _acme-challenge.example.com. CNAME _acme-challenge.example.com.alias.com.
   ```

3. 在 AutoSSL 控制台提交基础域名 `example.com`。
4. AutoSSL 使用 Cloudflare DNS 钩子调用 acme.sh：

   ```sh
   acme.sh --issue --dns dns_cf \
     -d example.com -d '*.example.com' \
     --challenge-alias example.com.alias.com \
     --challenge-alias example.com.alias.com \
     --keylength ec-256
   ```

5. acme.sh 使用 `CF_Token` 在 Cloudflare 区域写入 TXT 验证记录。Let's Encrypt 查询目标域名后会沿着 CNAME 找到该记录：

   ```text
   _acme-challenge.example.com
     CNAME -> _acme-challenge.example.com.alias.com
                 TXT <- acme.sh 通过 Cloudflare API 写入
   ```

6. Let's Encrypt 验证 TXT 记录后，为 `example.com` 和 `*.example.com` 签发证书。AutoSSL 将完整证书链和私钥安装到固定的下载 URL。

AutoSSL 每天在容器时间 01:30 运行 `acme.sh --cron`。请保留所有 CNAME，并持久化 `/root/.acme.sh`，让续期继续使用相同的验证路径和 ACME 账户。

## 配置与部署

### 1. 准备 Cloudflare 验证域名

将 `alias.com` 这样的域名添加到 Cloudflare，并将它专门用作 DNS-01 验证区域。创建仅授权该区域的 API Token，授予 DNS 编辑和区域读取权限，然后记录 Zone ID。

这个 Token 只需要访问 `alias.com`，不需要访问任何目标域名。

### 2. 添加目标域名的 CNAME

为每个目标基础域名添加上文所示的 CNAME。如果目标域名也使用 Cloudflare，请将记录的代理状态设为 **仅限 DNS**。签发完成后不要删除该记录，自动续期仍会使用它。

### 3. 部署 AutoSSL

创建 `docker-compose.yml`：

```yaml
services:
  autossl:
    image: ghcr.io/pupilcc/autossl:master
    container_name: autossl
    restart: always
    ports:
      - "1323:1323"
      - "3000:3000"
    volumes:
      - data:/root/data
      - acme:/root/.acme.sh
    environment:
      PUBLIC_API_URL: https://api.example.com
      PUBLIC_CONSOLE_URL: https://ssl.example.com
      ADMIN_USERNAME: admin
      ADMIN_PASSWORD: replace-with-a-long-random-password
      ACME_CA: letsencrypt
      ACME_EMAIL: admin@example.com
      ACME_DNS: dns_cf
      ACME_ALIAS: alias.com
      ACME_ALIAS_PER_DOMAIN: "true"
      ACME_DEBUG: "false"
      CF_Zone_ID: xxxxxxxx
      CF_Token: xxxxxx

volumes:
  data:
  acme:
```

启动服务：

```sh
docker compose up -d
```

配置 HTTPS 反向代理：将 `api.example.com` 转发到 Go API 的 `1323` 端口，将 `ssl.example.com` 转发到 Web 控制台的 `3000` 端口。登录后只需提交 `example.com` 这样的基础域名，AutoSSL 会自动加入通配符域名。

## 配置项

| 变量 | 用途 |
| --- | --- |
| `PUBLIC_API_URL` | Go API 的公开访问地址，用于生成 `/dl/...` 证书链接，请勿以 `/` 结尾。 |
| `PUBLIC_CONSOLE_URL` | 允许提交管理操作的公开控制台 URL。 |
| `ADMIN_USERNAME` | 控制台用户名。 |
| `ADMIN_PASSWORD` | 控制台密码和 JWT 签名密钥，请使用足够长的随机值。 |
| `ACME_CA` | acme.sh CA 名称，例如 `letsencrypt`。 |
| `ACME_EMAIL` | ACME 账户邮箱。 |
| `ACME_DNS` | acme.sh DNS 钩子；Cloudflare 验证区域使用 `dns_cf`。 |
| `ACME_ALIAS` | 由 Cloudflare 托管的统一验证域名。 |
| `ACME_ALIAS_PER_DOMAIN` | 仅 `true` 为每个目标域名使用 `<目标基础域名>.<ACME_ALIAS>`；其他值使所有域名共用 `ACME_ALIAS`。详见下文。 |
| `ACME_DEBUG` | `true` 表示启用 acme.sh 一级调试日志。 |
| `CF_Zone_ID` | Cloudflare 验证域名所在区域的 Zone ID。 |
| `CF_Token` | 有权管理该区域 DNS 记录的 API Token。 |

### `ACME_ALIAS_PER_DOMAIN` 如何工作

该变量决定 AutoSSL 是为每个目标域名生成独立验证别名，还是让所有目标域名直接共用 `ACME_ALIAS`。只有值为 `true`（不区分大小写）时才会启用；未设置、`false` 或其他值均表示关闭。

假设 `ACME_ALIAS=alias.com`，签发目标为 `example.com`：

| 设置 | 传给 acme.sh 的验证别名 | CNAME 目标 | 结果 |
| --- | --- | --- | --- |
| `ACME_ALIAS_PER_DOMAIN=true` | `example.com.alias.com` | `_acme-challenge.example.com.alias.com` | 每个目标基础域名使用独立记录，推荐。 |
| 未设置或非 `true` | `alias.com` | `_acme-challenge.alias.com` | 所有目标域名共用同一记录，并发验证可能互相干扰。 |

启用后，别名按 `<目标基础域名>.<ACME_ALIAS>` 生成。AutoSSL 会移除输入末尾的点和开头的 `*.`，并转为小写；例如 `*.Example.COM.` 会得到 `example.com.alias.com`。同一张证书中的 `example.com` 和 `*.example.com` 会共用这个别名，但其他目标域名使用各自的别名，因此可以安全地共享同一个 Cloudflare 区域和 API Token。

每个目标域名的 CNAME 必须与当前模式匹配。该变量只在 AutoSSL 发起新的 `acme.sh --issue` 时参与生成参数；`acme.sh --cron` 会使用已保存的签发配置。修改 `ACME_ALIAS` 或此变量时，请保留旧 CNAME，直到相关证书按新配置重新签发并确认续期路径已经更新。

## Cloudflare Universal SSL 冲突

如果目标域名也使用 Cloudflare，Universal SSL 可能在目标区域发布隐藏的 `_acme-challenge` TXT 记录，导致 Let's Encrypt 以 `Incorrect TXT record` 拒绝委派验证。

遇到此问题时，请在目标域名所在区域的 **SSL/TLS -> Edge Certificates** 中关闭 Universal SSL，等待 DNS 缓存过期后重试。此操作会影响依赖 Cloudflare 边缘证书的代理站点，请先确认这些站点已经配置其他有效的边缘证书。详情参见 [Cloudflare 说明](https://developers.cloudflare.com/dns/manage-dns-records/troubleshooting/unexpected-dns-records/#acme_challenge-txt-records)。

## 证书分发

控制台会为每张证书提供 `.crt` 和 `.key` URL。为了供部署脚本使用，这些下载路由不需要身份验证；请将私钥 URL 视为机密信息，并且只通过 HTTPS 暴露 AutoSSL。

在每台目标服务器上使用 [scripts/distribute-certificates.sh](scripts/distribute-certificates.sh)。替换 `urls_and_paths` 中的占位 URL 和路径，保留 `fullchain.pem` 与 `privkey.pem` 文件名，并以 root 身份运行。脚本会原子下载文件，在 OpenSSL 可用时校验内容，修正文件权限，并且只在文件变化且 `nginx -t` 成功后重新加载 Nginx。
