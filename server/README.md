# ライブラリ
gormConnect()内で.envファイル（環境変数定義）から定数を取得する
```bash
go get github.com/joho/godotenv
```

Usersテーブルにパスワードをそのまま保存するとセキュリティ的に危ないので、これを使って暗号化して保存
```bash
go get golang.org/x/crypto/bcrypt
```

gorm
```bash
$ go get github.com/jinzhu/gorm
```

"github.com/gin-contrib/sessions"
"github.com/gin-contrib/sessions/cookie"
