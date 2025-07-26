# bitcoin-app

## 起動方法

### dev環境

#### mysql


起動コマンド
```brew services start mysql```

停止コマンド
```brew services stop mysql```

```brew install mysql```をしている前提。

#### drf

`bitcoin-app/python/drf`に移動。

パッケージインストール
```pip install -r requirements.txt```

サーバー起動
```python manage.py runserver```

マイグレーション
```python manage.py migrate```

#### golang server

`bitcoin-app/golang`に移動

パッケージインストール
```go mod tidy```

サーバー起動
```go run cmd/server/main.go```

#### ticker batch

`bitcoin-app/golang`に移動

パッケージインストール
```go mod tidy```

サーバー起動
```go run cmd/ticker_batch/main.go```

mysqlへの接続
```mysql -h localhost -P 3306 -u root mysql_api_local```

prodのtickersテーブルの情報を開発環境にコピーする (結構時間がかかる)
```cd ~/bitcoin-app/golang && go run cmd/import_tickers/main.go```


### prod環境

golangサーバの再起動
```cd ~/bitcoin-app && sh restart-golang-server.sh```

drfの再起動
```cd ~/bitcoin-app && sh restart-drf.sh```

ticker-batchの再起動
```cd ~/bitcoin-app && sh restart-ticker-batch.sh```

すべてのサービスの再起動
```docker-compose -f production.yml build --no-cache && docker-compose -f production.yml down && docker-compose -f production.yml up -d```

初回起動
```docker-compose -f production.yml up --build -d```

ログ出力
```docker-compose -f production.yml logs -f```

mysqlへの接続
```docker exec -it bitcoin-mysql-prod bash``` &&
```mysql -h mysql -P 3306 -u bitcoin_user -p bitcoin_app``` (passは`bitcoin_password`)
