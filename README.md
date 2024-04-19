## マイグレーション

```
go run src/migrate/migrate.go -migrate
```

## テーブル全削除

```
go run src/migrate/migrate.go -drop
```

## コンパイル

```
go run src/main.go
```

## 注意点

refresh_token は最初のログインの時しか返してくれない.(refresh_token は有効時間が長いので発行しすぎは注意しなければならないから) もし、開発時に refresh_token を発行しなおさないといけないってなったら、以下のリンクにアクセスして、permission を消すとまた発行しなおされる.

https://myaccount.google.com/u/0/permissions?pli=1

## 今後のメモ

- 企業マスタのレコード削除の際は、物理削除なら ID で外部結合してる関連のテーブルの該当レコードも全削除する。論理削除なら、レコードは削除しつつ、データは S3 に保管したい
