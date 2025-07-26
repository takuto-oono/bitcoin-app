#!/bin/bash

echo "DRFサービスを再起動しています..."

# DRFサービスのみをビルドして再起動
docker-compose -f production.yml build --no-cache drf
docker-compose -f production.yml stop drf
docker-compose -f production.yml up -d drf

echo "DRFサービスの再起動が完了しました。"
