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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 本番用の軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates musl-dev gcc
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/server .
COPY --from=builder /app/toml ./toml
COPY --from=builder /app/env ./env

# ポートを公開
EXPOSE 8080

# 環境変数を設定
ENV GIN_MODE=release

# 起動コマンド
CMD ["./server", "-toml", "toml/prod.toml", "-env", "env/.env.prod"]
