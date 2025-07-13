# ビルドステージ
FROM golang:1.24.4-alpine AS builder

# 必要なビルドツールのインストール
RUN apk add --no-cache git

# 作業ディレクトリの設定
WORKDIR /app

# go.mod と go.sum をコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションのビルド
# CGO_ENABLED=0: 静的リンクを有効にする
# -ldflags="-s -w": バイナリサイズを削減するためのフラグ
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bitcoin-app .

# 実行ステージ
FROM alpine:latest

# 作業ディレクトリの設定
WORKDIR /app

# ビルドステージから実行可能ファイルをコピー
COPY --from=builder /app/bitcoin-app .

# 設定ファイルをコピー
COPY --from=builder /app/toml/prod.toml ./toml/prod.toml

# ポートの公開
EXPOSE 7080

# アプリケーションの実行
# プロダクション環境の設定ファイルを使用
CMD ["./bitcoin-app", "-toml", "toml/prod.toml"]
