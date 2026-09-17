# AutoSSL

AutoSSL is an open-source project built on top of acme.sh, designed to provide an SSL certificate distribution service. With AutoSSL, you can generate SSL certificates on one server and distribute them to other servers via HTTP. This project is inspired by the [vx.link](https://vx.link) SSL certificate service.

## Features

- **Centralized SSL Certificate Generation**: Generate SSL certificates on a single server.
- **HTTP Distribution**: Distribute the generated certificates to other servers via HTTP.
- **DNS Alias Mode**: Uses DNS alias mode for certificate generation. Please refer to the [acme.sh documentation](https://github.com/acmesh-official/acme.sh/wiki/DNS-alias-mode) for more details.
- **Web Console**: Sign in to add and remove domain certificates, then copy certificate links for deployment scripts.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

Create a `docker-compose.yml` file with the following content:

```yaml
version: '3.7'
services:
  autossl:
    image: ghcr.io/pupilcc/autossl:master
    container_name: autossl
    restart: always
    volumes:
      - data:/root/data
      - acme:/root/.acme.sh
    ports:
      - "3000:3000"
    environment:
      - DOMAIN=https://example.com
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=123456
      - ACME_CA=letsencrypt
      - ACME_EMAIL=example@gmail.com
      - ACME_DNS=dns_cf
      - ACME_ALIAS=alias.com
      - ACME_ALIAS_PER_DOMAIN=true
      - ACME_DEBUG=false
      - CF_Zone_ID=xxxxxxxx
      - CF_Token=xxxxxx

volumes:
  acme:
  data:
```

Run the following command to start the service:
```sh
docker-compose up -d
```

Open `http://localhost:3000` and sign in with `ADMIN_USERNAME` and `ADMIN_PASSWORD`.

The container exposes only the web service. It calls the Go API over the container loopback interface, and `/dl/:file` remains available through the web service for certificate distribution.

### Configuration

- `DOMAIN`: The domain for the SSL certificate.
- `ADMIN_USERNAME`: The username for the admin interface.
- `ADMIN_PASSWORD`: The password for the admin interface.
- `ACME_CA`: The Certificate Authority (e.g., letsencrypt).
- `ACME_EMAIL`: The email address for ACME registration.
- `ACME_DNS`: The DNS provider for ACME (e.g., dns_cf for Cloudflare).
- `ACME_ALIAS`: The DNS alias mode for ACME.
- `ACME_ALIAS_PER_DOMAIN`: Appends the requested domain to `ACME_ALIAS` when set to `true`, isolating concurrent DNS challenges by domain.
- `ACME_DEBUG`: Enables acme.sh debug level 1 logging when set to `true`.
- `CF_Zone_ID`: The Cloudflare Zone ID.
- `CF_Token`: The Cloudflare API token.

When `ACME_ALIAS_PER_DOMAIN=true`, `*.example.com` with `ACME_ALIAS=alias.com` requires this DNS record:

```dns
_acme-challenge.example.com. CNAME _acme-challenge.example.com.alias.com.
```

Submit the base domain (`example.com`) to the generate API. AutoSSL issues one certificate for both `example.com` and `*.example.com`; the wildcard already covers `www.example.com`.

#### Cloudflare Universal SSL conflict

When the certificate domain uses Cloudflare DNS, Universal SSL can automatically serve its own `_acme-challenge` TXT records instead of the configured CNAME alias. These records may not appear in the DNS dashboard. As a result, acme.sh can report `All checks succeeded` when checking the alias, while Let's Encrypt fails with `Incorrect TXT record` when querying the original domain.

To resolve this conflict, select the certificate domain's zone in Cloudflare (for example, `example.com`, rather than the alias zone `alias.com`), go to **SSL/TLS → Edge Certificates → Universal SSL**, and disable Universal SSL. Keep the challenge CNAME configured, allow cached DNS responses to expire, then retry issuance.

Disabling Universal SSL affects HTTPS for Cloudflare-proxied sites that rely on its edge certificates. Ensure those sites have another valid edge certificate before disabling it. See [Cloudflare's explanation of automatic ACME TXT records](https://developers.cloudflare.com/dns/manage-dns-records/troubleshooting/unexpected-dns-records/#acme_challenge-txt-records).

### Certificate Download Script

For detailed instructions on how to download the certificates, please refer to the [certificate download script](https://github.com/tmplink/KnowledgeBase/blob/main/vxlink/vxssl.md).

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/pupilcc/autossl/blob/master/LICENSE) file for more details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgments

- [acme.sh](https://github.com/acmesh-official/acme.sh)
- [vx.link SSL Certificate Service](https://vx.link)
