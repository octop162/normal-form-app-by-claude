# 会員登録webフォーム

React/TypeScript + Go + PostgreSQLによる3画面構成の会員登録フォームシステム

## 技術スタック

- **フロントエンド**: React + TypeScript + Vite
- **バックエンド**: Go + Gin Framework
- **データベース**: PostgreSQL
- **インフラ**: Docker + Docker Compose
- **開発環境**: ローカル開発 + Docker

## 必要な環境

- **Node.js**: v18以上
- **Go**: v1.21以上
- **Docker**: v20以上
- **Docker Compose**: v2以上

## 🚀 開発環境セットアップ

### 1. リポジトリクローン

```bash
git clone https://github.com/octop162/normal-form-app-by-claude.git
cd normal-form-app-by-claude
```

### 2. 環境変数設定

```bash
# .envファイルをコピー（既に存在する場合はスキップ）
cp .env.example .env

# 必要に応じて.envファイルを編集
vim .env
```

### 3. 依存関係インストール

```bash
# Go依存関係
go mod download

# React依存関係
cd frontend
npm install
cd ..
```

## 🏃 開発環境起動

### オプション1: 全サービス一括起動（推奨）

```bash
# PostgreSQL + Backend + Frontend をまとめて起動
docker-compose up -d postgres
go run cmd/server/main.go &
cd frontend && npm run dev &
```

### オプション2: サービス個別起動

#### PostgreSQLデータベース起動

```bash
docker-compose up -d postgres

# ログ確認
docker-compose logs postgres
```

#### Goバックエンド起動

```bash
# 開発モード（ホットリロード無し）
go run cmd/server/main.go

# バックグラウンド起動
go run cmd/server/main.go &
```

#### Reactフロントエンド起動

```bash
cd frontend
npm run dev

# 特定ポート指定
npm run dev -- --port 3000

# 外部アクセス許可
npm run dev -- --host 0.0.0.0
```

## 📍 アクセス情報

開発環境起動後、以下のURLでアクセスできます：

| サービス | URL | 説明 |
|---------|-----|------|
| **React Frontend** | http://localhost:5173 | 開発用フロントエンド |
| **Go Backend** | http://localhost:8080 | RESTful API |
| **Health Check** | http://localhost:8080/health | サーバー状態確認 |
| **API Test** | http://localhost:8080/api/v1/ping | API接続テスト |
| **PostgreSQL** | localhost:5432 | データベース |

## 🧪 動作確認

### ヘルスチェック

```bash
# Go serverの状態確認
curl http://localhost:8080/health

# 期待される応答
# {"service":"normal-form-app","status":"ok","version":"1.0.0"}
```

### API接続テスト

```bash
# APIエンドポイントテスト
curl http://localhost:8080/api/v1/ping

# 期待される応答
# {"message":"pong"}
```

### データベース接続テスト

```bash
# PostgreSQLコンテナに接続
docker exec -it normal-form-db psql -U postgres -d normal_form_db

# 接続後、SQLで確認
\dt  -- テーブル一覧
SELECT * FROM health_check;  -- 初期データ確認
\q   -- 終了
```

## 🛑 サービス停止

### 全サービス停止

```bash
# 実行中のプロセス停止
pkill -f "go run"
pkill -f "vite"
pkill -f "npm"

# Dockerコンテナ停止
docker-compose down
```

### 個別サービス停止

```bash
# PostgreSQL停止
docker-compose stop postgres

# Go/React はCtrl+Cまたは該当プロセスを停止
```

## 🐳 Docker開発環境

### 完全Docker環境での起動

```bash
# 全サービスをDockerで起動
docker-compose --profile backend --profile frontend up

# バックグラウンド起動
docker-compose --profile backend --profile frontend up -d
```

### Dockerログ確認

```bash
# 全サービスのログ
docker-compose logs -f

# 特定サービスのログ
docker-compose logs -f postgres
docker-compose logs -f backend
docker-compose logs -f frontend
```

## 📁 プロジェクト構造

```
normal-form-app-by-claude/
├── cmd/server/main.go          # Go アプリケーション エントリーポイント
├── internal/                   # Go 内部パッケージ
│   ├── handler/               # HTTPハンドラー
│   ├── service/               # ビジネスロジック
│   ├── repository/            # データアクセス層
│   ├── model/                 # ドメインモデル
│   ├── dto/                   # Data Transfer Object
│   └── middleware/            # ミドルウェア
├── pkg/                       # Go 共有パッケージ
│   ├── database/              # DB接続
│   ├── validator/             # バリデーター
│   └── logger/                # ログ
├── frontend/                  # React アプリケーション
│   ├── src/
│   │   ├── components/        # Reactコンポーネント
│   │   ├── pages/             # ページコンポーネント
│   │   ├── hooks/             # カスタムフック
│   │   ├── services/          # API通信
│   │   ├── types/             # TypeScript型定義
│   │   └── utils/             # ユーティリティ
│   ├── package.json
│   └── vite.config.ts
├── scripts/init.sql           # データベース初期化スクリプト
├── docker-compose.yml         # Docker Compose設定
├── .env                       # 環境変数（ローカル設定）
├── .env.example              # 環境変数例
└── README.md                 # このファイル
```

## 🔧 開発コマンド

### Go 関連

```bash
# 依存関係追加
go get github.com/some/package

# 依存関係整理
go mod tidy

# テスト実行
go test ./...

# ビルド
go build -o app cmd/server/main.go

# フォーマット
go fmt ./...
```

### React 関連

```bash
cd frontend

# 依存関係追加
npm install package-name

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# プレビュー
npm run preview

# リント
npm run lint
```

### Docker 関連

```bash
# イメージビルド
docker-compose build

# コンテナ再作成
docker-compose up --build

# ボリューム削除（データベースリセット）
docker-compose down -v

# 未使用リソース削除
docker system prune
```

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### ポート競合エラー

```bash
# ポート使用状況確認
lsof -i :8080  # Go server
lsof -i :5173  # React dev server  
lsof -i :5432  # PostgreSQL

# プロセス停止
kill -9 <PID>
```

#### PostgreSQL接続エラー

```bash
# コンテナ状態確認
docker-compose ps

# PostgreSQLログ確認
docker-compose logs postgres

# データベース再起動
docker-compose restart postgres
```

#### Go依存関係エラー

```bash
# モジュールキャッシュクリア
go clean -modcache

# 依存関係再取得
go mod download
go mod tidy
```

#### React起動エラー

```bash
# node_modules再インストール
cd frontend
rm -rf node_modules package-lock.json
npm install

# キャッシュクリア
npm run dev -- --force
```

#### Docker関連エラー

```bash
# Dockerシステム情報
docker system df

# 未使用リソース削除
docker system prune -a

# ボリューム確認
docker volume ls

# ネットワーク確認
docker network ls
```

## 📋 開発フロー

1. **issue確認**: GitHub Projectsでタスク選択
2. **ブランチ作成**: `git checkout -b feature/task-name`
3. **開発**: ローカル環境で実装・テスト
4. **動作確認**: 全サービス起動して統合テスト
5. **コミット**: `git commit -m "feat: description"`
6. **Push**: `git push origin feature/task-name`
7. **PR作成**: GitHub上でPull Request作成

## 📚 参考情報

- **プロジェクト仕様**: [CLAUDE.md](./CLAUDE.md)
- **GitHub Issues**: https://github.com/octop162/normal-form-app-by-claude/issues
- **GitHub Projects**: https://github.com/users/octop162/projects/7
- **API仕様**: 今後 `/api` ディレクトリに追加予定
- **E2Eテスト**: 今後 Playwright で実装予定

## 📞 サポート

- **Issues**: バグ報告や機能要望は GitHub Issues へ
- **Discussions**: 質問や議論は GitHub Discussions へ
- **Documentation**: 詳細仕様は CLAUDE.md を参照

---

**Happy Coding! 🚀**