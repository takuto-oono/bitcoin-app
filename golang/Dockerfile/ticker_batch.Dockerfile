FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# 必要なシステムパッケージをインストール
RUN apk add --no-cache git

# go.modとgo.sumをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ticker_batch ./cmd/ticker_batch

# 本番用の軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/ticker_batch .
COPY --from=builder /app/toml ./toml
COPY --from=builder /app/env ./env

# 起動コマンド
CMD ["./ticker_batch", "-toml", "toml/prod.toml", "-env", "env/.env.prod"]
