# =====================
# 1️⃣ 构建阶段（Go）
# =====================
FROM golang:1.22-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download || true

COPY . .

RUN go mod tidy || true
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# =====================
# 2️⃣ 运行阶段（Debian）
# =====================
FROM debian:12-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /app/app
COPY --from=builder /app/migrations /app/migrations

RUN mkdir -p /app/logs /app/uploads

EXPOSE 8080

CMD ["/app/app"]
