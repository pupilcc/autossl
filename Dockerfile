FROM node:24-bookworm-slim AS web-builder
WORKDIR /src/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build && npm prune --omit=dev

FROM golang:1.21.1 AS go-builder
RUN apt-get update && apt-get install -y gcc libc6-dev
WORKDIR /go/src/autossl
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/autossl .

FROM node:24-bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates curl openssl cron tini \
    && rm -rf /var/lib/apt/lists/*
RUN curl https://get.acme.sh | sh
WORKDIR /root/
COPY --from=go-builder /go/bin/autossl .
COPY --from=web-builder /src/web/build ./web/build
COPY --from=web-builder /src/web/node_modules ./web/node_modules
COPY --from=web-builder /src/web/package.json ./web/package.json
COPY --from=web-builder /src/web/server.mjs ./web/server.mjs
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh
ENV AUTOSSL_ADDR=127.0.0.1:1323 \
    AUTOSSL_API_URL=http://127.0.0.1:1323 \
    NODE_ENV=production \
    PORT=3000
EXPOSE 3000
ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["./docker-entrypoint.sh"]
