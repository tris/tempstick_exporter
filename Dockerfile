FROM golang:1.20-alpine as builder
MAINTAINER Tristan Horn <tristan+docker@ethereal.net>
WORKDIR /app
RUN apk add --no-cache upx
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o tempstick_exporter .
RUN upx --lzma tempstick_exporter

FROM scratch
COPY --from=builder /app/tempstick_exporter /tempstick_exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/tempstick_exporter"]
