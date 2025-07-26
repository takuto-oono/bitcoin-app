#!/bin/bash

echo "Golangサービスを再起動しています..."

# Golangサーバーのみをビルドして再起動
docker-compose -f production.yml build --no-cache golang-server
docker-compose -f production.yml stop golang-server
docker-compose -f production.yml up -d golang-server

echo "Golangサービスの再起動が完了しました。"
