# # goバージョン
# FROM golang:1.20-bullseye
# # コンテナ内の作業ディレクトリを作成
# WORKDIR /app
# # アップデートとgitのインストール
# RUN apt-get update && apt-get install -y git

# ECRプッシュ用

# Goのバージョンを指定
FROM golang:1.20-bullseye

# コンテナ内の作業ディレクトリを作成
WORKDIR /app

# アップデートとgitのインストール（Goモジュールのダウンロードに必要）
RUN apt-get update && apt-get install -y git

# Goの依存関係をコンテナにコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコンテナにコピー
COPY . .

# Goアプリケーションをビルド（イメージビルド時）
RUN go build -o app src/main.go

# 必要なポートを公開
EXPOSE 8080

# コンテナ起動時に実行するコマンド（マイグレーションを先に実行してアプリケーションを起動）
CMD ["sh", "-c", "go run src/migrate/migrate.go -drop && go run src/migrate/migrate.go -migrate && ./app"]