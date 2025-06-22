# 会員登録webフォーム 開発仕様書

## プロジェクト概要

会員登録のためのwebフォームシステム。入力画面、確認画面、完了画面の3画面構成で、複雑な申込みオプションとリアルタイムバリデーションを持つ。

## 技術スタック

### フロントエンド
- **React** (TypeScript推奨)
- レスポンシブデザイン対応

### バックエンド
- **Go** (Gin framework推奨)
- **Wire** (Dependency Injection)
- クリーンアーキテクチャ採用

### データベース
- **PostgreSQL** (Aurora)
- 接続プール使用

### インフラ・運用
- **AWS ECS** (Fargate)
- **ECR** (コンテナレジストリ)
- **CloudWatch** (ログ・モニタリング)
- **AWS Backup** + Aurora スナップショット

### ログレベル
- ローカル環境：`debug`
- 開発環境：`info`
- 本番環境：`warn`

## 画面構成

### 1. 入力画面
- 個人情報入力
- プラン・オプション選択
- リアルタイムバリデーション

### 2. 確認画面
- 入力内容確認
- サーバサイドバリデーション実施
- エラー時は入力画面に戻る

### 3. 完了画面
- 登録完了メッセージ表示

## セッション管理

- **タイムアウト時間：4時間**
- タイムアウト時は警告表示
- **一時保存機能：必須**（別作業後の再開対応）

## 入力項目仕様

### 個人情報

#### 氏名
- **姓**：最大15文字、必須
- **名**：最大15文字、必須

#### 氏名カナ
- **姓カナ**：全角カタカナ限定、最大15文字、必須
- **名カナ**：全角カタカナ限定、最大15文字、必須

#### 電話番号（3分割入力）
- **制限ルール：**
  - フリーダイヤル（0120、0800等）：NG
  - 11桁の場合：0X0で始まる（090、080、070等の携帯番号）
  - 市外局番：2～5桁
  - 市内局番：1～4桁
  - 契約番号：4桁固定

#### 郵便番号（分割入力）
- **前3桁**：数字3桁、必須
- **後4桁**：数字4桁、必須
- **住所自動検索機能：**
  - 検索ボタン配置
  - API連携で住所取得
  - 見つからない場合は専用メッセージ表示

#### 住所（詳細分割）
- **都道府県**：必須
- **市区町村**：必須（政令指定都市の区も考慮）
- **町名**：条件付き必須（住所によっては不要）
- **丁目**：任意
- **番地**：必須
- **号**：任意
- **建物名**：任意
- **部屋番号**：任意

#### メールアドレス
- **メールアドレス**：最大256文字、RFC準拠、必須
- **メールアドレス（確認用）**：上記と同一チェック、必須

## プラン・オプション仕様

### プラン種類
- **Aプラン**
- **Bプラン**

### オプション種類
- **AAオプション**：Aプランのみ選択可能
- **BBオプション**：Bプランのみ選択可能
- **ABオプション**：A・B両プラン共通

### オプション選択
- **複数選択可能**
- **組み合わせ制限なし**（自由選択）

## API連携仕様

### 在庫状況API
- **形式**：JSON
- **応答時間**：0.5秒以内
- **レスポンス内容**：オプション毎の在庫数
- **呼び出しタイミング：**
  - 画面表示時
  - サーバサイド入力チェック時
- **在庫0の場合**：申し込み不可
- **API失敗時**：申し込み不可

### 住所検索API
- 郵便番号から住所情報を取得
- 自動入力対応

### 地域制限API
- **判定レベル**：市区町村
- **制限対象**：特定オプションのみ
- プラン自体の制限はなし

## バリデーション仕様

### リアルタイムバリデーション
- **実行タイミング**：フォーカス離脱時（onBlur）
- **相関チェック**：全入力完了後
- **必須チェック**：「次へ」ボタン押下時

### エラー表示
- **表示位置**：各入力フィールドの下
- **表示内容**：単項目のみ（複数エラーは一つずつ）
- **色・デザイン**：後日指定

### サーバサイドバリデーション
- 確認画面遷移時に全項目チェック
- 重複チェックは不要
- エラー時は入力画面に戻る

## セキュリティ要件

### 通信セキュリティ
- **HTTPS必須**
- データベース保存時の暗号化は不要

### セキュリティレベル
- **一般的な企業レベル**
- CSRF対策実装
- アクセス制御は不要（将来のログイン機能で対応）

## データベース設計

### 主要テーブル（想定）
```sql
-- ユーザー情報
users (
    id SERIAL PRIMARY KEY,
    last_name VARCHAR(15) NOT NULL,
    first_name VARCHAR(15) NOT NULL,
    last_name_kana VARCHAR(15) NOT NULL,
    first_name_kana VARCHAR(15) NOT NULL,
    phone1 VARCHAR(5) NOT NULL,
    phone2 VARCHAR(4) NOT NULL,
    phone3 VARCHAR(4) NOT NULL,
    postal_code1 CHAR(3) NOT NULL,
    postal_code2 CHAR(4) NOT NULL,
    prefecture VARCHAR(10) NOT NULL,
    city VARCHAR(50) NOT NULL,
    town VARCHAR(50),
    chome VARCHAR(10),
    banchi VARCHAR(10) NOT NULL,
    go VARCHAR(10),
    building VARCHAR(100),
    room VARCHAR(20),
    email VARCHAR(256) NOT NULL,
    plan_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 申し込みオプション
user_options (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    option_type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- セッション管理（一時保存）
user_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_data JSONB NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## アーキテクチャ設計

### ディレクトリ構成（Go）
```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/          # HTTPハンドラー
│   ├── service/          # ビジネスロジック
│   ├── repository/       # データアクセス層
│   ├── model/           # ドメインモデル
│   ├── dto/             # Data Transfer Object
│   └── middleware/      # ミドルウェア
├── pkg/
│   ├── database/        # DB接続
│   ├── validator/       # バリデーター
│   └── logger/          # ログ
├── api/                 # API仕様書
├── deployments/         # Docker, ECS設定
└── scripts/            # 各種スクリプト
```

### レイヤー構成
```
Handler → Service → Repository → Database
   ↓         ↓         ↓
Interface Interface Interface
```

### 主要インターフェース
```go
type UserService interface {
    CreateUser(ctx context.Context, req *dto.UserCreateRequest) error
    ValidateUserData(ctx context.Context, req *dto.UserCreateRequest) error
    SaveTemporaryData(ctx context.Context, sessionID string, data *dto.UserCreateRequest) error
    LoadTemporaryData(ctx context.Context, sessionID string) (*dto.UserCreateRequest, error)
}

type UserRepository interface {
    Save(ctx context.Context, user *model.User) error
    ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type OptionService interface {
    GetAvailableOptions(ctx context.Context, planType string, region string) ([]model.Option, error)
    CheckInventory(ctx context.Context, optionIDs []string) (map[string]int, error)
}

type AddressService interface {
    SearchByPostalCode(ctx context.Context, postalCode string) (*model.Address, error)
}
```

## API設計

### 実装済みエンドポイント（Phase 1）

#### ヘルスチェック系
```
GET    /health                          # 総合ヘルスチェック
GET    /health/live                     # Liveness Probe（Kubernetes対応）
GET    /health/ready                    # Readiness Probe（Kubernetes対応）
```

#### テスト・デバッグ系
```
GET    /api/v1/ping                     # API疎通確認
```

### 実装予定エンドポイント（Phase 2以降）

#### ユーザー登録系
```
POST   /api/v1/users                    # ユーザー登録
POST   /api/v1/users/validate           # バリデーション
```

#### セッション管理系
```
POST   /api/v1/sessions                 # セッション作成（一時保存）
GET    /api/v1/sessions/{id}            # セッション取得
PUT    /api/v1/sessions/{id}            # セッション更新
DELETE /api/v1/sessions/{id}            # セッション削除
```

#### オプション・在庫系
```
GET    /api/v1/options                  # オプション一覧取得
POST   /api/v1/options/check-inventory  # 在庫チェック
```

#### 住所・地域系
```
GET    /api/v1/address/search           # 住所検索（郵便番号）
POST   /api/v1/region/check             # 地域制限チェック
```

#### マスターデータ系
```
GET    /api/v1/prefectures              # 都道府県一覧取得
GET    /api/v1/plans                    # プラン一覧取得
```

### レスポンス形式
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
    Code    string            `json:"code"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}
```

## フロントエンド設計

### ディレクトリ構成（React）
```
src/
├── components/
│   ├── common/          # 共通コンポーネント
│   ├── forms/           # フォームコンポーネント
│   └── layout/          # レイアウトコンポーネント
├── pages/
│   ├── UserInput.tsx    # 入力画面
│   ├── UserConfirm.tsx  # 確認画面
│   └── UserComplete.tsx # 完了画面
├── hooks/               # カスタムフック
├── services/            # API通信
├── types/               # 型定義
├── utils/               # ユーティリティ
└── validation/          # バリデーション
```

### 状態管理
- **React Context** または **Redux Toolkit**
- セッション管理とフォームデータの永続化

### バリデーション
- **react-hook-form** + **yup** または **zod**
- リアルタイムバリデーション実装

## 外部API仕様

### 在庫確認API
```javascript
// Request
POST /external/api/inventory/check
{
  "option_ids": ["AA", "BB", "AB"]
}

// Response
{
  "success": true,
  "data": {
    "AA": 10,
    "BB": 0,
    "AB": 5
  }
}
```

### 地域制限API
```javascript
// Request
POST /external/api/region/check
{
  "prefecture": "東京都",
  "city": "渋谷区",
  "option_ids": ["AA", "BB"]
}

// Response
{
  "success": true,
  "data": {
    "AA": true,
    "BB": false
  }
}
```

## デプロイ・運用

### Docker設定
```dockerfile
# Goアプリケーション用
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### ECS設定
- **Fargate**使用
- **Application Load Balancer**
- **Auto Scaling**設定
- **CloudWatch Logs**統合

### 環境変数
```
# データベース
DB_HOST=
DB_PORT=5432
DB_NAME=
DB_USER=
DB_PASSWORD=

# 外部API
INVENTORY_API_URL=
REGION_API_URL=
ADDRESS_API_URL=

# その他
LOG_LEVEL=warn
SESSION_TIMEOUT=4h
```

## コード標準化・品質管理

### 命名規則

#### Go言語
- **パッケージ名**：小文字のみ（例：`user`, `handler`, `repository`）
- **公開関数/メソッド**：パスカルケース（例：`GetUserByID`, `ValidateEmail`）
- **非公開関数/メソッド**：キャメルケース（例：`validateEmail`, `parseRequest`）
- **構造体**：パスカルケース（例：`UserService`, `DatabaseConfig`）
- **インターフェース**：動詞+er形式推奨（例：`UserRepository`, `Validator`）
- **ファイル名**：スネークケース（例：`user_service.go`, `database_config.go`）
- **略語**：全て大文字（例：`ID`, `URL`, `HTTP`, `API`）

#### TypeScript/React
- **コンポーネント名**：パスカルケース（例：`UserInputForm`, `AddressSearchButton`）
- **変数/関数名**：キャメルケース（例：`handleSubmit`, `validateForm`）
- **ファイル名**：
  - コンポーネント：パスカルケース（例：`UserInput.tsx`）
  - その他：キャメルケース（例：`userService.ts`, `apiClient.ts`）
- **Props/State**：
  - キャメルケース使用
  - boolean値：`is`, `has`, `can`で開始（例：`isLoading`, `hasError`, `canSubmit`）
- **定数**：UPPER_SNAKE_CASE（例：`API_BASE_URL`, `MAX_RETRY_COUNT`）

#### データベース
- **テーブル名**：スネークケース、複数形（例：`users`, `user_options`）
- **カラム名**：スネークケース、単数形（例：`user_id`, `created_at`）
- **インデックス名**：`idx_テーブル名_カラム名`（例：`idx_users_email`）
- **外部キー名**：`fk_テーブル名_参照テーブル名`（例：`fk_user_options_users`）

#### API
- **エンドポイント**：REST準拠、複数形（例：`/api/v1/users`, `/api/v1/options`）
- **パラメータ名**：スネークケース（例：`user_id`, `option_type`）
- **レスポンスフィールド**：スネークケース（例：`created_at`, `is_active`）

### 静的解析ツール設定

#### Go言語
- **golangci-lint**：包括的なlinter集合
- **gofmt**：コードフォーマット
- **goimports**：import文整理
- **go vet**：静的解析
- **staticcheck**：高度な静的解析
- **ineffassign**：無効な代入検出
- **misspell**：スペルチェック

#### TypeScript/React
- **ESLint**：コード品質チェック
- **Prettier**：コードフォーマット
- **TypeScript Compiler**：型チェック
- **stylelint**：CSS品質チェック
- **Playwright**：E2Eテストフレームワーク

#### 設定ファイル管理
```
├── .golangci.yml          # Go静的解析設定
├── .eslintrc.js           # ESLint設定
├── .prettierrc            # Prettier設定
├── tsconfig.json          # TypeScript設定
├── playwright.config.ts   # Playwright E2E設定
├── .editorconfig          # エディタ統一設定
└── .vscode/
    ├── settings.json      # VS Code設定
    └── extensions.json    # 推奨拡張機能
```

### Git管理・CI/CD品質ゲート

#### Pre-commit Hooks
- コードフォーマット自動実行
- 静的解析自動実行
- テスト自動実行（高速なもののみ）

#### CI/CDパイプライン品質チェック
```yaml
# GitHub Actions例
name: Quality Gate
on: [push, pull_request]

jobs:
  go-quality:
    steps:
      - golangci-lint実行
      - go test実行
      - カバレッジ測定（80%以上）

  frontend-quality:
    steps:
      - ESLint実行
      - Prettier チェック
      - TypeScript型チェック
      - Jest テスト実行

  e2e-tests:
    needs: [go-quality, frontend-quality]
    steps:
      - データベース起動（PostgreSQL）
      - Goサーバー起動
      - Reactアプリ起動
      - Playwright E2Eテスト実行
      - 複数ブラウザテスト（Chrome, Firefox, Safari）
      - モバイルレスポンシブテスト
```

#### 品質基準
- **静的解析**：エラー0件必須
- **テストカバレッジ**：80%以上
- **型安全性**：TypeScript strict mode
- **セキュリティ**：gosec, npm audit クリア
- **E2Eテスト**：主要ユーザーフロー100%カバー
- **クロスブラウザ対応**：Chrome, Firefox, Safari対応
- **レスポンシブ対応**：モバイル・タブレット・デスクトップ

### AI実装ガイドライン

#### Claude Code使用時の標準化
- 命名規則の厳密な遵守
- インターフェース駆動開発
- 単一責任の原則
- 依存関係注入パターン

#### 実装パターン統一
```go
// Go: 標準的なHandler実装パターン
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.userService.CreateUser(c.Request.Context(), &req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
```

```typescript
// TypeScript: 標準的なComponent実装パターン
interface UserInputProps {
  onSubmit: (data: UserFormData) => void;
  isLoading: boolean;
}

export const UserInput: React.FC<UserInputProps> = ({ onSubmit, isLoading }) => {
  const { register, handleSubmit, formState: { errors } } = useForm<UserFormData>();

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* フォーム実装 */}
    </form>
  );
};
```

## 開発・テスト方針

### テスト戦略
- **Unit Test**：各レイヤーの単体テスト（カバレッジ80%以上）
- **Integration Test**：API統合テスト
- **E2E Test**：Playwrightによる画面遷移・ユーザーフローテスト
- **Static Analysis Test**：静的解析による品質担保

#### E2Eテスト詳細（Playwright）
- **テスト対象**：
  - 会員登録フロー（入力→確認→完了）
  - バリデーションエラーハンドリング
  - 外部API連携（在庫確認・住所検索・地域制限）
  - セッション管理・一時保存機能
  - レスポンシブデザイン対応
- **テスト環境**：
  - 複数ブラウザ（Chrome, Firefox, Safari）
  - モバイルデバイス（Pixel 5等）
  - CI/CD統合
- **レポート**：HTML形式、スクリーンショット・動画付き

### 開発フロー
1. **ローカル開発**（Docker Compose + 静的解析 + Playwright）
2. **Pull Request**（CI/CD品質ゲート + E2Eテスト）
3. **開発環境デプロイ**（ECS）
4. **統合テスト実施**（E2Eテスト含む）
5. **本番デプロイ**

### 監視・アラート
- **CloudWatch Metrics**
- **CloudWatch Alarms**
- エラー率、レスポンス時間の監視
- **コード品質メトリクス**（SonarQube等）

## 今後の拡張予定

### ログイン機能
- 将来的にログイン機能追加予定
- その際にアクセス制御を実装

### プラン・オプション拡張
- 柔軟性を持たせた設計
- 新プラン・オプション追加に対応

### 多言語対応
- 国際化（i18n）対応の準備

## E2Eテスト仕様（Playwright）

### テストシナリオ

#### 1. 正常系フロー
```typescript
// 基本的な会員登録フロー
test('会員登録正常フロー', async ({ page }) => {
  // 入力画面
  await page.goto('/');
  await page.fill('[data-testid="last-name"]', '田中');
  await page.fill('[data-testid="first-name"]', '太郎');
  // ... 全項目入力
  await page.click('[data-testid="next-button"]');

  // 確認画面
  await expect(page.locator('[data-testid="confirm-name"]')).toContainText('田中 太郎');
  await page.click('[data-testid="submit-button"]');

  // 完了画面
  await expect(page.locator('[data-testid="success-message"]')).toBeVisible();
});
```

#### 2. バリデーションエラー
```typescript
test('必須項目未入力エラー', async ({ page }) => {
  await page.goto('/');
  await page.click('[data-testid="next-button"]');

  await expect(page.locator('[data-testid="last-name-error"]')).toContainText('姓は必須です');
  await expect(page.locator('[data-testid="first-name-error"]')).toContainText('名は必須です');
});

test('メールアドレス形式エラー', async ({ page }) => {
  await page.goto('/');
  await page.fill('[data-testid="email"]', 'invalid-email');
  await page.blur('[data-testid="email"]');

  await expect(page.locator('[data-testid="email-error"]')).toContainText('正しいメールアドレスを入力してください');
});
```

#### 3. 外部API連携
```typescript
test('住所自動検索', async ({ page }) => {
  await page.goto('/');
  await page.fill('[data-testid="postal-code-1"]', '100');
  await page.fill('[data-testid="postal-code-2"]', '0001');
  await page.click('[data-testid="address-search-button"]');

  await expect(page.locator('[data-testid="prefecture"]')).toHaveValue('東京都');
  await expect(page.locator('[data-testid="city"]')).toHaveValue('千代田区');
});

test('在庫切れオプション制御', async ({ page }) => {
  // モックAPIで在庫0を返すよう設定
  await page.route('/api/v1/options/check-inventory', route => {
    route.fulfill({ json: { AA: 0, BB: 5, AB: 3 } });
  });

  await page.goto('/');
  await page.selectOption('[data-testid="plan"]', 'A');

  await expect(page.locator('[data-testid="option-AA"]')).toBeDisabled();
  await expect(page.locator('[data-testid="option-AB"]')).toBeEnabled();
});
```

#### 4. セッション管理
```typescript
test('一時保存・復元機能', async ({ page }) => {
  await page.goto('/');
  await page.fill('[data-testid="last-name"]', '田中');
  await page.fill('[data-testid="first-name"]', '太郎');

  // ページリロード
  await page.reload();

  // データが復元されることを確認
  await expect(page.locator('[data-testid="last-name"]')).toHaveValue('田中');
  await expect(page.locator('[data-testid="first-name"]')).toHaveValue('太郎');
});

test('セッションタイムアウト警告', async ({ page }) => {
  // タイムアウト時間を短縮してテスト
  await page.addInitScript(() => {
    window.SESSION_TIMEOUT = 1000; // 1秒
  });

  await page.goto('/');
  await page.waitForTimeout(1500);

  await expect(page.locator('[data-testid="timeout-warning"]')).toBeVisible();
});
```

#### 5. レスポンシブデザイン
```typescript
test('モバイル表示確認', async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 667 });
  await page.goto('/');

  // モバイル専用要素の確認
  await expect(page.locator('[data-testid="mobile-menu"]')).toBeVisible();
  await expect(page.locator('[data-testid="desktop-menu"]')).toBeHidden();
});
```

### テスト設定

#### Playwright設定（playwright.config.ts）
```typescript
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [
    ['html'],
    ['junit', { outputFile: 'test-results/junit.xml' }]
  ],
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure'
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    }
  ],
  webServer: [
    {
      command: './server',
      url: 'http://localhost:8080',
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'npm run dev',
      url: 'http://localhost:3000',
      reuseExistingServer: !process.env.CI,
    }
  ]
});
```

### CI/CD統合

#### GitHub Actions E2E設定
```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    timeout-minutes: 60
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_DB: testdb
          POSTGRES_USER: testuser
          POSTGRES_PASSWORD: testpass
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432

    steps:
    - uses: actions/checkout@v4

    - uses: actions/setup-node@v4
      with:
        node-version: 18

    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: npm ci

    - name: Build Go server
      run: go build -o server cmd/server/main.go

    - name: Start servers
      run: |
        ./server &
        npm run dev &

    - name: Install Playwright Browsers
      run: npx playwright install --with-deps

    - name: Run Playwright tests
      run: npx playwright test

    - uses: actions/upload-artifact@v4
      if: always()
      with:
        name: playwright-report
        path: playwright-report/
        retention-days: 30

    - uses: actions/upload-artifact@v4
      if: always()
      with:
        name: test-results
        path: test-results/
        retention-days: 30
```

---

## 開発開始時の注意点

1. **セキュリティ**：HTTPS必須、CSRF対策実装
2. **パフォーマンス**：API応答時間0.5秒以内を厳守
3. **ユーザビリティ**：エラーメッセージの分かりやすさ
4. **保守性**：クリーンアーキテクチャの徹底
5. **スケーラビリティ**：将来の機能拡張を考慮した設計

この仕様書を基に、段階的な開発を開始してください。不明点や仕様変更があれば随時更新してください。

## プロジェクト管理

### 実装管理方法
- **GitHub Issues**: 全タスクはGitHub issueで管理されています
- **フェーズ構成**: 7フェーズ、32個のissueに分割済み
- **ラベル体系**: `phase-1`～`phase-7`、技術カテゴリ、作業種別

### 実装順序とマイルストーン
1. **MVP完成**: フェーズ1-3完了時点（Week 1-5）
2. **機能完成**: フェーズ1-5完了時点（Week 1-9）
3. **本番準備完了**: 全フェーズ完了時点（Week 1-12）

### 優先度
- 🔴 **高**: フェーズ1-3（基盤構築、バックエンドコア、外部API）
- 🟡 **中**: フェーズ4-5（フロントエンド、セキュリティ・品質）
- 🟢 **低**: フェーズ6-7（インフラ・デプロイ、運用準備）

## GitHub開発フロー

### GitHub Issues管理

#### Issue構成
- **総数**: 32個のissue（全フェーズ対応）
- **命名規則**: `【Phase X】タスク名`
- **ラベル体系**:
  - **フェーズ**: `phase-1` ～ `phase-7`
  - **技術領域**: `backend`, `frontend`, `database`, `api`, `infrastructure`
  - **作業種別**: `setup`, `documentation`, `testing`, `security`

#### Issue作成時の標準テンプレート
```markdown
## 概要
タスクの目的と背景

## タスク
- [ ] 具体的な作業項目1
- [ ] 具体的な作業項目2

## 受け入れ条件
- [ ] 完了の基準1
- [ ] 完了の基準2

## 参考
CLAUDE.md の該当セクション
```

### GitHub Projects管理

#### プロジェクトボード情報
- **プロジェクト名**: 会員登録webフォーム開発
- **URL**: https://github.com/users/octop162/projects/7
- **管理対象**: 全32 issues

#### ボードカラム構成
- **📝 Todo**: 未着手タスク
- **🚀 In Progress**: 作業中タスク
- **👀 Review**: レビュー・確認中タスク
- **✅ Done**: 完了タスク

#### プロジェクトビュー
- **Board View**: 進捗の可視化
- **Table View**: 詳細情報とフィルタリング
- **Roadmap View**: タイムライン管理

### 自動化ワークフロー

#### Issue自動追加ワークフロー
- **ファイル**: `.github/workflows/add-to-project.yml`
- **トリガー**: Issue作成・再開時
- **機能**: 新しいissueを自動的にプロジェクト#7に追加

#### 必要な設定
1. **Personal Access Token作成**:
   - Scopes: `repo`, `project`
   - Name: `Project Auto-add Token`

2. **Repository Secret設定**:
   - Name: `ADD_TO_PROJECT_PAT`
   - Value: 上記で作成したPAT

### 開発ワークフロー

#### 1. タスク開始フロー
```
1. GitHub Projects でTodoからタスクを選択
2. Issue を "In Progress" に移動
3. ローカルでfeatureブランチ作成
4. 実装・テスト実施
5. Pull Request作成
6. レビュー・マージ
7. Issue を "Done" に移動
```

#### 2. ブランチ戦略
```
main (master)
├── feature/phase1-setup          # フェーズ1実装
├── feature/phase2-backend-core    # フェーズ2実装
├── feature/phase3-external-api    # フェーズ3実装
└── hotfix/security-fix           # 緊急修正
```

#### 3. コミットメッセージ規約
```
Type: Short description

- Detailed explanation
- Reference to issue #XX

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

**Types:**
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: コードスタイル
- `refactor`: リファクタリング
- `test`: テスト
- `chore`: その他

#### 4. Pull Request テンプレート
```markdown
## 概要
このPRで解決する課題と実装内容

## 変更内容
- [ ] 変更点1
- [ ] 変更点2

## テスト
- [ ] 単体テスト実行
- [ ] 統合テスト実行
- [ ] 手動テスト実行

## チェックリスト
- [ ] コードレビュー済み
- [ ] テストが通る
- [ ] ドキュメント更新済み

## 関連Issue
Closes #XX
```

### 品質管理

#### CI/CDパイプライン（予定）
- **静的解析**: golangci-lint, ESLint
- **テスト**: Unit, Integration, E2E
- **ビルド**: Docker image作成
- **デプロイ**: ECS環境への自動デプロイ

#### コードレビュー基準
- **機能要件**: 仕様書通りの実装
- **コード品質**: 命名規則、アーキテクチャ遵守
- **テスト**: 適切なテストカバレッジ
- **セキュリティ**: セキュリティ要件クリア

### 進捗管理

#### 日次確認事項
- [ ] GitHub Projects ボードで進捗確認
- [ ] 作業中issueのステータス更新
- [ ] ブロッカーや課題の特定

#### 週次確認事項
- [ ] フェーズ進捗の確認
- [ ] マイルストーン達成状況確認
- [ ] 次週の作業計画策定

#### 報告フォーマット
```
## Week X Progress Report

### 完了したタスク
- Issue #XX: タスク名

### 進行中のタスク
- Issue #XX: タスク名 (進捗XX%)

### 次週の予定
- Issue #XX: タスク名（着手予定）

### ブロッカー・課題
- 課題があれば記載
```
