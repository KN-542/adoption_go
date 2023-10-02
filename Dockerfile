# goバージョン
FROM golang:1.20-bullseye
# コンテナ内の作業ディレクトリを作成
WORKDIR /app
# アップデートとgitのインストール
RUN apt-get update && apt-get install -y git
# # boiler-plateディレクトリの作成
# RUN mkdir /go/src/github.com/boiler-plate
# # ワーキングディレクトリの設定
# WORKDIR /go/src/github.com/boiler-plate
# # ホストのファイルをコンテナの作業ディレクトリに移行
# ADD . /go/src/github.com/boiler-plate
# # パッケージのインポート
# RUN go get -u golang.org/x/tools/cmd/goimports
