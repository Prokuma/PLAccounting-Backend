# PLAccounting Backend
青色申告兼家計簿用の自作会計システムのバックエンド

## 開発環境
- Go 1.21.0
- PostgreSQL 15.4
- Redis 7.2.0
- SMTPサーバ

## 環境構築
### Docker Compose利用
### 環境変数設定
`.env.docker`に環境変数を設定する。設定方法は下記参照。
```shell
cp .env.example .env.docker
```

#### setup.shの実行
フロントエンドのビルドを行う。`node`と`npm`が必要。
```shell
chmod +x ./setup.sh
./setup.sh
```

#### 環境変数の読み込み
```shell
source .env.docker
```

#### Dockerfileのビルド&起動
```shell
docker-compose up -d
``` 

### 手動構築
#### JWT用公開鍵・秘密鍵生成
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

#### PostgreSQL, Redisサーバの構築及び設定
適宜それぞれのサーバを立ち上げ、接続情報を環境変数にて設定する。
```shell
POSTGRES_HOST=localhost
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=test
POSTGRES_SSLMODE=false
POSTGRES_TIMEZONE=Asia/Tokyo
REDIS_HOST=localhost
REDIS_PASSWORD=password
REDIS_PORT=6379
REDIS_DB=0
```

#### メール認証用の送信アカウント設定
無効化するオプションは未実装なので、任意ではなく必須。
```shell
SMTP_HOST=host
SMTP_PORT=465
SMTP_USER=user
SMTP_PASS=password
SMTP_USERADDR=test@test.com
```

#### ハッシュソルト設定（任意）
予め用意されたハッシュテーブルによるパスワード復号化対策。
完全ではないが事故時のリスクを減らせる。
```shell
HASH_SALT=qawsedrftgyhujikolp
```

#### 依存パッケージ導入
```bash
go mod tidy
```

#### 起動
`go run`利用（ビルド後実行）
```bash
go run main.go
```

ビルド
```bash
go build
```

## API仕様
Swaggerで作成しているので当該ファイル（`docs/`以下のファイル）参考。
また、起動後`host:port/swagger/index.html`でもアクセス可