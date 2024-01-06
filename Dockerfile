# goバージョン
FROM golang:1.20-bullseye
# コンテナ内の作業ディレクトリを作成
WORKDIR /app
# アップデートとgitのインストール
RUN apt-get update && apt-get install -y git