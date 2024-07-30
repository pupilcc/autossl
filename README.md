# AutoSSL

AutoSSL is an open-source project built on top of acme.sh, designed to provide an SSL certificate distribution service. With AutoSSL, you can generate SSL certificates on one server and distribute them to other servers via HTTP. This project is inspired by the [vx.link](https://vx.link) SSL certificate service.

## Features

- **Centralized SSL Certificate Generation**: Generate SSL certificates on a single server.
- **HTTP Distribution**: Distribute the generated certificates to other servers via HTTP.
- **DNS Alias Mode**: Uses DNS alias mode for certificate generation. Please refer to the [acme.sh documentation](https://github.com/acmesh-official/acme.sh/wiki/DNS-alias-mode) for more details.

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
      - "1323:1323"
    environment:
      - DOMAIN=https://example.com
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=123456
      - ACME_CA=letsencrypt
      - ACME_EMAIL=example@gmail.com
      - ACME_DNS=dns_cf
      - ACME_ALIAS=alias.com
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

### Configuration

- `DOMAIN`: The domain for the SSL certificate.
- `ADMIN_USERNAME`: The username for the admin interface.
- `ADMIN_PASSWORD`: The password for the admin interface.
- `ACME_CA`: The Certificate Authority (e.g., letsencrypt).
- `ACME_EMAIL`: The email address for ACME registration.
- `ACME_DNS`: The DNS provider for ACME (e.g., dns_cf for Cloudflare).
- `ACME_ALIAS`: The DNS alias mode for ACME.
- `CF_Zone_ID`: The Cloudflare Zone ID.
- `CF_Token`: The Cloudflare API token.

### Certificate Download Script

For detailed instructions on how to download the certificates, please refer to the [certificate download script](https://github.com/tmplink/KnowledgeBase/blob/main/vxlink/vxssl.md).

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/pupilcc/autossl/blob/master/LICENSE) file for more details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgments

- [acme.sh](https://github.com/acmesh-official/acme.sh)
- [vx.link SSL Certificate Service](https://vx.link)
