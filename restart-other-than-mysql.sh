#!/bin/bash

echo "DRFサービスのビルドをしています..."
docker-compose -f production.yml build --no-cache drf

echo "fast-apiサービスのビルドをしています..."
docker-compose -f production.yml build --no-cache fast-api

echo "Golangサービスのビルドをしています..."
docker-compose -f production.yml build --no-cache golang-server

echo "バッチサービスのビルドをしています..."
docker-compose -f production.yml build --no-cache ticker-batch

echo "すべてのサービスのビルドが完了しました。再起動します..."

echo "DRFサービスを再起動しています..."
docker-compose -f production.yml stop drf
docker-compose -f production.yml up -d drf
echo "DRFサービスの再起動が完了しました。"

echo "fast-apiサービスを再起動しています..."
docker-compose -f production.yml stop fast-api
docker-compose -f production.yml up -d fast-api
echo "fast-apiサービスの再起動が完了しました。"

echo "Golangサービスを再起動しています..."
docker-compose -f production.yml stop golang-server
docker-compose -f production.yml up -d golang-server
echo "Golangサービスの再起動が完了しました。"

echo "バッチサービスを再起動しています..."
docker-compose -f production.yml stop ticker-batch
docker-compose -f production.yml up -d ticker-batch
echo "バッチサービスの再起動が完了しました。"

echo "すべてのサービスの再起動が完了しました。"
