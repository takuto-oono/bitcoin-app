FROM python:3.13.5-alpine3.21

# 環境変数を設定
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# 作業ディレクトリを設定
WORKDIR /app

# requirements.txtをコピーして依存関係をインストール
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# プロジェクトファイルをコピー
COPY . .

# ポートを公開
EXPOSE 8001

# 起動コマンド
CMD ["python", "main.py", "--env", "env/.env.prod"]
