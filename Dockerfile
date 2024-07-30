FROM golang:1.21.1 as builder
RUN apt-get update && apt-get install -y gcc libc6-dev
WORKDIR /go/src/autossl
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -o /go/bin/autossl .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates curl openssl \
    && rm -rf /var/lib/apt/lists/*
RUN curl https://get.acme.sh | sh
WORKDIR /root/
COPY --from=builder /go/bin/autossl .
EXPOSE 1323
CMD ["./autossl"]