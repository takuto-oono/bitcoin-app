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


### prod環境

再起動
```docker-compose -f production.yml build --no-cache && docker-compose -f production.yml down && docker-compose -f production.yml up -d```

初回起動
```docker-compose -f production.yml up --build -d```

ログ出力
```docker-compose -f production.yml logs -f```
