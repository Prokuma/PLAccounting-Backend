# PLAccounting Backend
青色申告兼家計簿用の自作会計システムのバックエンド

## 開発環境
- Go 1.21.0
- PostgreSQL 15.4
- Redis 7.2.0

## 環境構築
### JWT用公開鍵・秘密鍵生成
PKCS8形式で生成する。パスワード設定はなし。
なお、公開鍵に関しては、`ssh-keygen -m PKCS8`のみで対応不可能なため、変換を行う。
```shell
ssh-keygen -m PKCS8
ssh-keygen -i -f .jwt_pub -m pkcs8 > .jwt_rsa.pub.pkcs8
```

生成後は、それぞれのファイルのパスを環境変数として設定する。
```shell
JWT_PRIVATE_KEY_PATH=./.jwt_rsa
JWT_PUBLIC_KEY_PATH=./.jwt_rsa.pub.pkcs8
```

### PostgreSQL, Redisサーバの構築及び設定
適宜それぞれのサーバを立ち上げ、接続情報を環境変数にて設定する。
```shell
POSTGRESQL_HOST=localhost
POSTGRESQL_USER=user
POSTGRESQL_PASSWORD=password
POSTGRESQL_DBNAME=test
POSTGRESQL_SSLMODE=false
POSTGRESQL_TIMEZONE=Asia/Tokyo
REDIS_HOST=localhost
REDIS_PASSWORD=password
REDIS_PORT=6379
REDIS_DB=0
```

### ハッシュソルト設定（任意）
予め用意されたハッシュテーブルによるパスワード復号化対策。
完全ではないが事故時のリスクを減らせる。
```shell
HASH_SALT=qawsedrftgyhujikolp
```

### 依存パッケージ導入
```bash
go mod tidy
```

### 起動
`go run`利用（ビルド後実行）
```bash
go run main.go
```

ビルド
```bash
go build
```

## API仕様
後ほど作成する