# コーディング規約・命名規則

このドキュメントは、会員登録webフォームプロジェクトにおけるコーディング規約と命名規則を定義します。

## 目次

1. [Go 命名規則・コーディング規約](#go-命名規則コーディング規約)
2. [TypeScript/React 命名規則・コーディング規約](#typescriptreact-命名規則コーディング規約)
3. [データベース 命名規則](#データベース-命名規則)
4. [API 命名規則](#api-命名規則)
5. [ファイル・ディレクトリ 命名規則](#ファイルディレクトリ-命名規則)

---

## Go 命名規則・コーディング規約

### パッケージ名
- **小文字のみ、アンダースコア不使用**
- 短く、意味のある名前
- 複数形は避ける

```go
// Good
package user
package database
package validator

// Bad
package User
package data_base
package users
```

### 変数・関数名
- **camelCase** を使用
- exportする場合は **PascalCase**
- 省略形は避け、意味のある名前をつける

```go
// Good
var userName string
var userID int
func GetUserByID(id int) (*User, error) {}
func validateEmail(email string) bool {}

// Bad
var un string
var user_name string
func get_user(id int) {}
```

### 構造体・インターフェース名
- **PascalCase** を使用
- インターフェースは機能を表す名前 + "er" または意味のある名前

```go
// Good
type User struct {}
type UserRepository interface {}
type EmailSender interface {}
type Validator interface {}

// Bad
type user struct {}
type IUser interface {}
type UserRepositoryInterface interface {}
```

### 定数名
- **PascalCase** を使用（export する場合）
- **camelCase** を使用（package 内のみの場合）
- 関連する定数はグループ化

```go
// Good
const (
    MaxRetryCount = 3
    DefaultTimeout = 30 * time.Second
)

const (
    StatusActive = iota
    StatusInactive
    StatusPending
)

// Bad
const MAX_RETRY_COUNT = 3
const max_retry_count = 3
```

### エラー処理
- エラーメッセージは小文字で開始
- context を含める
- ラップする場合は `fmt.Errorf` または `errors.Wrap` を使用

```go
// Good
return fmt.Errorf("failed to create user: %w", err)
return errors.New("invalid email format")

// Bad
return errors.New("Failed to create user")
return err // context なし
```

### コメント
- 公開関数・構造体には必ずコメントを記述
- コメントは英語で記述
- 関数名で開始

```go
// GetUserByID retrieves a user by their unique identifier.
// Returns an error if the user is not found or database error occurs.
func GetUserByID(id int) (*User, error) {
    // implementation
}

// User represents a registered user in the system.
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

---

## TypeScript/React 命名規則・コーディング規約

### 変数・関数名
- **camelCase** を使用
- boolean変数は `is`, `has`, `can`, `should` で開始

```typescript
// Good
const userName = 'john';
const isLoggedIn = true;
const hasPermission = false;
const getUserData = () => {};

// Bad
const user_name = 'john';
const UserName = 'john';
const loggedIn = true; // boolean であることが不明確
```

### 型・インターフェース名
- **PascalCase** を使用
- インターフェースは `I` プレフィックスを使用しない
- Propsは `ComponentNameProps` の形式

```typescript
// Good
interface User {
  id: number;
  name: string;
  email: string;
}

interface UserCardProps {
  user: User;
  onEdit: (user: User) => void;
}

type UserStatus = 'active' | 'inactive' | 'pending';

// Bad
interface IUser {}
interface userProps {}
interface Props {} // どのコンポーネントのPropsか不明
```

### React コンポーネント名
- **PascalCase** を使用
- ファイル名とコンポーネント名を一致させる
- HOC は `with` プレフィックス
- hooks は `use` プレフィックス

```typescript
// Good
// components/UserCard.tsx
export const UserCard: React.FC<UserCardProps> = ({ user, onEdit }) => {
  // implementation
};

// hooks/useUser.ts
export const useUser = (userId: number) => {
  // implementation
};

// hoc/withAuth.tsx
export const withAuth = <P extends object>(Component: React.ComponentType<P>) => {
  // implementation
};

// Bad
// components/usercard.tsx
export const userCard = () => {};
// components/UserCard.tsx
export const Card = () => {}; // ファイル名と不一致
```

### 定数名
- **SCREAMING_SNAKE_CASE** を使用（グローバル定数）
- **camelCase** を使用（コンポーネント内定数）

```typescript
// constants/api.ts
export const API_BASE_URL = 'https://api.example.com';
export const MAX_RETRY_COUNT = 3;

// component内
const UserCard: React.FC<UserCardProps> = ({ user }) => {
  const defaultAvatar = '/images/default-avatar.png';
  // implementation
};
```

### ファイル・フォルダ構成
```
src/
├── components/          # 再利用可能なコンポーネント
│   ├── common/         # 汎用コンポーネント
│   └── forms/          # フォーム専用コンポーネント
├── pages/              # ページコンポーネント
├── hooks/              # カスタムフック
├── services/           # API通信・外部サービス
├── types/              # 型定義
├── utils/              # ユーティリティ関数
├── constants/          # 定数
└── styles/             # スタイル関連
```

---

## データベース 命名規則

### テーブル名
- **snake_case** を使用
- 複数形を使用
- 接頭辞は使用しない

```sql
-- Good
CREATE TABLE users (...);
CREATE TABLE user_profiles (...);
CREATE TABLE membership_applications (...);

-- Bad
CREATE TABLE User (...);
CREATE TABLE tbl_users (...);
CREATE TABLE user (...); -- 単数形
```

### カラム名
- **snake_case** を使用
- 意味のある名前を使用
- 型を含めない

```sql
-- Good
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bad
CREATE TABLE users (
    ID INTEGER, -- 大文字
    fname VARCHAR(100), -- 省略形
    email_str VARCHAR(255), -- 型を含む
    createdAt TIMESTAMP -- camelCase
);
```

### インデックス名
- `idx_` プレフィックス + テーブル名 + カラム名

```sql
-- Good
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Bad
CREATE INDEX email_index ON users(email);
CREATE INDEX users_email ON users(email);
```

### 制約名
- 制約タイプのプレフィックス + テーブル名 + カラム名

```sql
-- Good
ALTER TABLE users ADD CONSTRAINT fk_users_profile_id 
    FOREIGN KEY (profile_id) REFERENCES user_profiles(id);
ALTER TABLE users ADD CONSTRAINT uk_users_email 
    UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT ck_users_email_format 
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Bad
ALTER TABLE users ADD CONSTRAINT users_profile_fkey 
    FOREIGN KEY (profile_id) REFERENCES user_profiles(id);
```

---

## API 命名規則

### REST API エンドポイント
- **kebab-case** を使用
- リソース指向の設計
- 動詞は HTTP メソッドで表現

```
# Good
GET    /api/v1/users                    # ユーザー一覧取得
GET    /api/v1/users/{id}              # 特定ユーザー取得
POST   /api/v1/users                   # ユーザー作成
PUT    /api/v1/users/{id}              # ユーザー更新
DELETE /api/v1/users/{id}              # ユーザー削除
GET    /api/v1/users/{id}/profile      # ユーザープロフィール取得
POST   /api/v1/membership-applications # 会員登録申請

# Bad
GET    /api/v1/getUsers               # 動詞を含む
POST   /api/v1/user_create           # snake_case & 動詞
GET    /api/v1/Users                 # PascalCase
```

### JSON フィールド名
- **camelCase** を使用
- 明確で説明的な名前

```json
// Good
{
  "id": 1,
  "firstName": "田中",
  "lastName": "太郎",
  "emailAddress": "tanaka@example.com",
  "phoneNumber": "090-1234-5678",
  "birthDate": "1990-01-01",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}

// Bad
{
  "ID": 1,                    // 大文字
  "first_name": "田中",       // snake_case
  "fname": "田中",            // 省略形
  "email": "tanaka@example.com", // emailAddress の方が明確
  "created": "2024-01-01T10:00:00Z" // created_at の方が明確
}
```

### HTTP ステータスコード
```
200 OK          - 成功（GET, PUT）
201 Created     - 作成成功（POST）
204 No Content  - 成功（DELETE）
400 Bad Request - リクエストエラー
401 Unauthorized - 認証エラー
403 Forbidden   - 認可エラー
404 Not Found   - リソース未発見
409 Conflict    - 競合エラー
422 Unprocessable Entity - バリデーションエラー
500 Internal Server Error - サーバーエラー
```

---

## ファイル・ディレクトリ 命名規則

### ディレクトリ構造
```
normal-form-app-by-claude/
├── cmd/                    # アプリケーションエントリーポイント
│   └── server/            # サーバーアプリケーション
├── internal/              # 内部パッケージ
│   ├── handler/           # HTTPハンドラー
│   ├── service/           # ビジネスロジック
│   ├── repository/        # データアクセス層
│   ├── model/             # ドメインモデル
│   ├── dto/               # Data Transfer Object
│   └── middleware/        # ミドルウェア
├── pkg/                   # 公開可能パッケージ
│   ├── database/          # DB接続・設定
│   ├── validator/         # バリデーター
│   └── logger/            # ログ
├── frontend/              # React アプリケーション
│   └── src/
│       ├── components/    # コンポーネント
│       ├── pages/         # ページ
│       ├── hooks/         # カスタムフック
│       ├── services/      # API通信
│       ├── types/         # 型定義
│       ├── utils/         # ユーティリティ
│       └── constants/     # 定数
├── scripts/               # スクリプト
├── docs/                  # ドキュメント
├── tests/                 # テスト
└── docker/                # Docker関連ファイル
```

### ファイル命名規則

#### Go ファイル
- **snake_case** を使用
- テスト ファイルは `_test.go` サフィックス
- 機能が明確に分かる名前

```
// Good
user_handler.go
user_service.go
user_repository.go
user_handler_test.go
email_validator.go

// Bad
userHandler.go          // camelCase
user-handler.go         // kebab-case
handler.go              // 不明確
userhandler.go          // 読みにくい
```

#### TypeScript/React ファイル
- **PascalCase** （コンポーネント）
- **camelCase** （その他）
- **kebab-case** （設定ファイル）

```
// Good
UserCard.tsx            // コンポーネント
UserCard.test.tsx       // テストファイル
userService.ts          // サービス
useUser.ts              // カスタムフック
apiConstants.ts         # 定数
types/User.ts           # 型定義

// Bad
usercard.tsx
user-card.tsx
UserService.ts          // PascalCase は不要
```

#### 設定・その他ファイル
- **kebab-case** または **snake_case**
- 用途が明確な名前

```
// Good
docker-compose.yml
.env.example
package.json
tsconfig.json
health-check.sh
dev-start.sh

// Bad
dockercompose.yml
envexample
package.JSON
```

### Git ブランチ命名規則
```
main                    # メインブランチ
develop                 # 開発ブランチ
feature/user-auth       # 機能ブランチ
feature/issue-123       # Issue番号ベース
bugfix/login-error      # バグ修正
hotfix/security-patch   # 緊急修正
release/v1.0.0          # リリースブランチ
```

### コミットメッセージ規約
```
feat: add user authentication
fix: resolve login validation error
docs: update API documentation
style: format code according to standards
refactor: restructure database connection
test: add unit tests for user service
chore: update dependencies
```

---

## 開発フロー

1. **Issue作成**: GitHub Issuesで作業内容を明確化
2. **ブランチ作成**: `git checkout -b feature/task-name`
3. **実装**: 本規約に従って実装
4. **テスト**: 動作確認・自動テスト実行
5. **コードレビュー**: Pull Request作成・レビュー
6. **マージ**: レビュー承認後にマージ

---

## 参考資料

- [Effective Go](https://golang.org/doc/effective_go.html)
- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- [PostgreSQL Naming Conventions](https://stackoverflow.com/questions/2878248/postgresql-naming-conventions)

---

**更新日**: 2024-06-22  
**バージョン**: 1.0.0