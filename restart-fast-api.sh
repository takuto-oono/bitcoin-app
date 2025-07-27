#!/bin/bash

echo "fast-apiサービスを再起動しています..."

# fast-apiサービスのみをビルドして再起動
docker-compose -f production.yml build --no-cache fast-api
docker-compose -f production.yml stop fast-api
docker-compose -f production.yml up -d fast-api

echo "fast-apiサービスの再起動が完了しました。"
