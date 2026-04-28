FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server /app/server

# 复制 GoFrame 配置和国际化目录
COPY --from=builder /app/manifest/config /app/config
COPY --from=builder /app/manifest/i18n /app/i18n

EXPOSE 8002

ENTRYPOINT ["/app/server"]