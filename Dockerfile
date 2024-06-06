FROM golang:1.21.1 as builder
WORKDIR /go/src/autossl
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/autossl .

FROM debian:12-slim

RUN apt update && apt install -y ca-certificates curl openssl cron
RUN curl https://get.acme.sh | sh

WORKDIR /root/
COPY --from=builder /go/bin/autossl .
EXPOSE 1323
ENTRYPOINT cron && ./autossl