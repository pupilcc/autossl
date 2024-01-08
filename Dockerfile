FROM golang:1.21.1 as builder
WORKDIR /go/src/autossl
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/autossl .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/autossl .
EXPOSE 1323
CMD ["./autossl"]