## 認証API

### 使用したライブラリ
- github.com/dgrijalva/jwt-go v3.2.0
- github.com/go-sql-driver/mysql v1.5.0
- github.com/gorilla/mux v1.7.4
- github.com/jinzhu/gorm v1.9.15
- github.com/subosito/gotenv v1.2.0 
- golang.org/x/crypto 

### URL


##### ユーザー登録
```
POST {host_url}/signup
```

##### ログイン
```
POST {host_url}/login
```

##### トークン認証
```
GET {host_url}/verify
```