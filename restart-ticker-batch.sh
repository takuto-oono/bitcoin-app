#!/bin/bash

echo "バッチサービスを再起動しています..."

# バッチサービスのみをビルドして再起動
docker-compose -f production.yml build --no-cache ticker-batch
docker-compose -f production.yml stop ticker-batch
docker-compose -f production.yml up -d ticker-batch

echo "バッチサービスの再起動が完了しました。"
