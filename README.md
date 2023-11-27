## マイグレーション

```
go run src/migrate/migrate.go
```

## コンパイル

```
go run src/main.go
```

## 注意点

refresh_token は最初のログインの時しか返してくれない.(refresh_token は有効時間が長いので発行しすぎは注意しなければならないから) もし、開発時に refresh_token を発行しなおさないといけないってなったら、以下のリンクにアクセスして、permission を消すとまた発行しなおされる.

https://myaccount.google.com/u/0/permissions?pli=1
