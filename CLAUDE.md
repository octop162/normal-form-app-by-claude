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

### エンドポイント
```
POST   /api/v1/users                    # ユーザー登録
POST   /api/v1/users/validate           # バリデーション
POST   /api/v1/sessions                 # セッション作成
GET    /api/v1/sessions/{id}            # セッション取得
PUT    /api/v1/sessions/{id}            # セッション更新
GET    /api/v1/options                  # オプション取得
POST   /api/v1/options/check-inventory  # 在庫チェック
GET    /api/v1/address/search           # 住所検索
POST   /api/v1/region/check             # 地域制限チェック
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

# 会員登録webフォーム実装TODOリスト

## フェーズ1: 基盤構築 🏗️

### 開発環境セットアップ
- [ ] Goプロジェクト初期化（go mod init）
- [ ] Reactプロジェクト初期化（Create React App + TypeScript）
- [ ] PostgreSQL環境構築（Docker Compose）
- [ ] ディレクトリ構造作成（仕様書通り）
- [ ] 基本的なDockerfile作成（Go/React）
- [ ] Docker Compose設定（ローカル開発用）
- [ ] 環境変数設定（.env ファイル）

### 命名規則・コード規約設定
- [ ] Go命名規則ドキュメント作成
- [ ] TypeScript/React命名規則ドキュメント作成
- [ ] データベース命名規則ドキュメント作成
- [ ] API命名規則ドキュメント作成
- [ ] ファイル・ディレクトリ命名規則ドキュメント作成

### 静的解析ツール設定（Go）
- [ ] golangci-lint設定（.golangci.yml）
- [ ] gofmt設定
- [ ] goimports設定
- [ ] golint設定
- [ ] go vet設定
- [ ] ineffassign設定
- [ ] misspell設定
- [ ] staticcheck設定

### 静的解析ツール設定（TypeScript/React）
- [ ] ESLint設定（.eslintrc.js）
- [ ] Prettier設定（.prettierrc）
- [ ] TypeScript設定（tsconfig.json）
- [ ] Husky設定（pre-commit hooks）
- [ ] lint-staged設定
- [ ] stylelint設定（CSS）
- [ ] Playwright設定（playwright.config.ts）

### エディタ設定
- [ ] VS Code設定（.vscode/settings.json）
- [ ] 推奨拡張機能リスト（.vscode/extensions.json）
- [ ] EditorConfig設定（.editorconfig）

### Go基盤実装
- [ ] 基本的なHTTPサーバー設定（Gin）
- [ ] データベース接続設定（PostgreSQL）
- [ ] ログ設定（レベル別：debug/info/warn）
- [ ] 基本ミドルウェア実装（CORS、ログ、エラーハンドリング）
- [ ] ヘルスチェックエンドポイント（/health）

### データベース設計・実装
- [ ] マイグレーションツール設定（golang-migrate）
- [ ] usersテーブル作成
- [ ] user_optionsテーブル作成
- [ ] user_sessionsテーブル作成
- [ ] インデックス設定
- [ ] 初期マスターデータ投入

## フェーズ2: バックエンドコア機能 ⚙️

### モデル・DTO定義
- [ ] Userモデル定義
- [ ] UserOptionモデル定義
- [ ] UserSessionモデル定義
- [ ] リクエスト/レスポンスDTO定義
- [ ] バリデーションルール定義

### Repository層実装
- [ ] UserRepository interface定義
- [ ] UserRepository実装（PostgreSQL）
- [ ] UserOptionRepository interface定義
- [ ] UserOptionRepository実装
- [ ] SessionRepository interface定義
- [ ] SessionRepository実装

### Service層実装
- [ ] UserService interface定義
- [ ] UserService実装
- [ ] OptionService interface定義
- [ ] OptionService実装
- [ ] AddressService interface定義
- [ ] AddressService実装
- [ ] バリデーションサービス実装

### Handler層実装
- [ ] ユーザー登録エンドポイント（POST /api/v1/users）
- [ ] バリデーションエンドポイント（POST /api/v1/users/validate）
- [ ] セッション管理エンドポイント（CRUD）
- [ ] オプション取得エンドポイント（GET /api/v1/options）
- [ ] 住所検索エンドポイント（GET /api/v1/address/search）

### DI設定（Wire）
- [ ] Wire設定ファイル作成
- [ ] 依存関係注入設定
- [ ] プロバイダー関数定義
- [ ] DIコンテナ初期化

## フェーズ3: 外部API連携 🔌

### 外部API実装
- [ ] HTTPクライアント設定（タイムアウト0.5秒）
- [ ] 在庫確認API連携実装
- [ ] 地域制限API連携実装
- [ ] 住所検索API連携実装
- [ ] APIエラーハンドリング実装
- [ ] リトライ処理実装（必要に応じて）

### 外部APIモック
- [ ] 在庫確認APIモック作成
- [ ] 地域制限APIモック作成
- [ ] 住所検索APIモック作成
- [ ] モック切り替え設定

## フェーズ4: フロントエンド実装 💻

### React基盤実装
- [ ] TypeScript型定義
- [ ] API通信層実装（axios/fetch）
- [ ] 共通コンポーネント作成
- [ ] レイアウトコンポーネント作成
- [ ] ルーティング設定（React Router）

### フォームコンポーネント実装
- [ ] 個人情報入力コンポーネント
  - [ ] 氏名入力（姓・名分離）
  - [ ] カナ入力（全角カタカナ制限）
  - [ ] 電話番号入力（3分割）
  - [ ] 郵便番号入力（3桁-4桁分離）
  - [ ] 住所入力（詳細分割）
  - [ ] メールアドレス入力（確認用含む）
- [ ] プラン・オプション選択コンポーネント
- [ ] 郵便番号検索機能実装

### バリデーション実装
- [ ] react-hook-form + yup設定
- [ ] リアルタイムバリデーション（onBlur）
- [ ] 相関チェック実装
- [ ] エラーメッセージ表示
- [ ] サーバサイドエラー処理

### 画面実装
- [ ] 入力画面（UserInput.tsx）
- [ ] 確認画面（UserConfirm.tsx）
- [ ] 完了画面（UserComplete.tsx）
- [ ] エラー画面
- [ ] ローディング表示

### セッション管理
- [ ] セッション状態管理（Context/Redux）
- [ ] 一時保存機能実装
- [ ] タイムアウト警告実装
- [ ] セッション復元機能

## フェーズ5: セキュリティ・品質 🔒

### セキュリティ実装
- [ ] CSRF対策実装
- [ ] HTTPS設定
- [ ] 入力サニタイゼーション
- [ ] SQLインジェクション対策確認
- [ ] XSS対策実装

### テスト実装
- [ ] Goユニットテスト
  - [ ] Repository層テスト
  - [ ] Service層テスト
  - [ ] Handler層テスト
- [ ] Go統合テスト
- [ ] Reactコンポーネントテスト
- [ ] E2Eテスト（Playwright）
  - [ ] Playwright設定・初期化
  - [ ] 会員登録フロー E2Eテスト
  - [ ] バリデーションエラー E2Eテスト
  - [ ] API連携 E2Eテスト
  - [ ] セッション管理 E2Eテスト
  - [ ] レスポンシブデザイン E2Eテスト
  - [ ] 複数ブラウザ対応テスト（Chrome, Firefox, Safari）
  - [ ] CI/CD用 E2Eテスト設定

### パフォーマンス最適化
- [ ] API応答時間測定・最適化
- [ ] データベースクエリ最適化
- [ ] フロントエンドバンドル最適化
- [ ] 画像・リソース最適化

## フェーズ6: インフラ・デプロイ ☁️

### AWS設定
- [ ] ECRリポジトリ作成
- [ ] Aurora PostgreSQL設定
- [ ] VPC・サブネット設定
- [ ] セキュリティグループ設定

### ECS設定
- [ ] ECSクラスター作成
- [ ] タスク定義作成
- [ ] サービス定義作成
- [ ] Application Load Balancer設定
- [ ] Auto Scaling設定

### CI/CD設定
- [ ] GitHub Actions設定
- [ ] ビルドパイプライン
- [ ] 静的解析パイプライン（Go + TypeScript）
- [ ] テストパイプライン
  - [ ] Goユニット・統合テスト
  - [ ] Reactコンポーネントテスト
  - [ ] E2Eテスト（Playwright）
- [ ] デプロイパイプライン
- [ ] 環境別デプロイ設定
- [ ] コード品質ゲート設定

### 監視・ログ
- [ ] CloudWatch Logs設定
- [ ] CloudWatch Metrics設定
- [ ] アラーム設定
- [ ] ダッシュボード作成

### バックアップ
- [ ] AWS Backup設定
- [ ] Aurora自動バックアップ設定
- [ ] スナップショット設定

## フェーズ7: 運用準備 🚀

### ドキュメント整備
- [ ] API仕様書更新
- [ ] 運用手順書作成
- [ ] 障害対応手順書作成
- [ ] コーディング規約書作成
- [ ] 命名規則ガイド作成
- [ ] AI実装ガイドライン作成
- [ ] README更新

### 負荷テスト
- [ ] 負荷テストシナリオ作成
- [ ] 負荷テスト実行
- [ ] パフォーマンスチューニング

### 本番準備
- [ ] 本番環境設定
- [ ] データ移行手順確認
- [ ] ロールバック手順確認
- [ ] 監視設定確認

---

## 命名規則・コード規約詳細

### Go言語命名規則

#### パッケージ名
- 小文字のみ使用
- 短く、簡潔に
- 例: `user`, `handler`, `repository`

#### 変数・関数名
- キャメルケース使用
- 公開: 大文字開始 `GetUserByID`
- 非公開: 小文字開始 `validateEmail`
- 略語は全て大文字 `ID`, `URL`, `HTTP`

#### 構造体名
- パスカルケース使用
- 例: `UserService`, `DatabaseConfig`

#### インターフェース名
- "er"で終わる場合が多い
- 例: `UserRepository`, `Validator`

#### ファイル名
- スネークケース使用
- 例: `user_service.go`, `database_config.go`

### TypeScript/React命名規則

#### コンポーネント名
- パスカルケース使用
- 例: `UserInputForm`, `AddressSearchButton`

#### 変数・関数名
- キャメルケース使用
- 例: `handleSubmit`, `validateForm`

#### ファイル名
- コンポーネント: パスカルケース `UserInput.tsx`
- その他: キャメルケース `userService.ts`

#### Props・State
- キャメルケース使用
- boolean は `is`, `has`, `can`で開始
- 例: `isLoading`, `hasError`, `canSubmit`

### データベース命名規則

#### テーブル名
- スネークケース、複数形
- 例: `users`, `user_options`

#### カラム名
- スネークケース、単数形
- 例: `user_id`, `created_at`

#### インデックス名
- `idx_テーブル名_カラム名`
- 例: `idx_users_email`

### API命名規則

#### エンドポイント
- REST準拠、複数形使用
- 例: `/api/v1/users`, `/api/v1/options`

#### パラメータ名
- スネークケース使用
- 例: `user_id`, `option_type`

---

## 静的解析ツール設定詳細

### golangci-lint設定例 (.golangci.yml)
```yaml
run:
  timeout: 5m
  go: "1.21"

linters:
  enable:
    - gofmt
    - golint
    - govet
    - ineffassign
    - misspell
    - staticcheck
    - unused
    - errcheck
    - gosimple
    - typecheck

linters-settings:
  golint:
    min-confidence: 0.8
  gofmt:
    simplify: true
```

### ESLint設定例 (.eslintrc.js)
```javascript
module.exports = {
  extends: [
    'react-app',
    'react-app/jest',
    '@typescript-eslint/recommended',
    'prettier'
  ],
  plugins: ['@typescript-eslint', 'react-hooks'],
  rules: {
    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: 'interface',
        format: ['PascalCase'],
        prefix: ['I']
      },
      {
        selector: 'typeAlias',
        format: ['PascalCase']
      }
    ],
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn'
  }
};
```

### Prettier設定例 (.prettierrc)
```json
{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2
}
```

---

## 推奨実装順序

### Week 1-2: 基盤構築
フェーズ1を完了し、ローカル開発環境を整える

### Week 3-4: バックエンドコア
フェーズ2を完了し、基本的なAPI機能を実装

### Week 5: 外部API連携
フェーズ3を完了し、外部APIとの連携を実装

### Week 6-8: フロントエンド
フェーズ4を完了し、UI/UXを実装

### Week 9: セキュリティ・品質
フェーズ5を完了し、テスト・セキュリティを強化

### Week 10-11: インフラ・デプロイ
フェーズ6を完了し、本番環境を構築

### Week 12: 運用準備
フェーズ7を完了し、本番リリース準備

---

## 注意事項

### 優先度
- 🔴 **高**: フェーズ1-3（機能的な動作確認まで）
- 🟡 **中**: フェーズ4-5（ユーザー体験・品質）
- 🟢 **低**: フェーズ6-7（本番運用準備）

### 並行作業可能項目
- フロントエンドとバックエンドは API仕様確定後に並行開発可能
- テスト実装は各フェーズと並行して実施
- インフラ設定は開発と並行して準備可能

### マイルストーン
1. **MVP完成**: フェーズ1-3完了時点
2. **機能完成**: フェーズ1-5完了時点  
3. **本番準備完了**: 全フェーズ完了時点

このTODOリストを基に、まずはどのフェーズから始めますか？