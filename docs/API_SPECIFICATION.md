# API仕様書

## 概要

会員登録webフォームのREST API仕様書です。

### 基本情報

- **ベースURL**: `https://api.normal-form-app.com`
- **バージョン**: v1
- **認証**: CSRFトークン
- **データ形式**: JSON
- **文字エンコーディング**: UTF-8

### 共通レスポンス形式

```json
{
  "success": boolean,
  "data": object | null,
  "error": {
    "code": string,
    "message": string,
    "details": object | null
  } | null,
  "meta": {
    "request_id": string,
    "timestamp": string,
    "path": string,
    "method": string
  } | null
}
```

### エラーコード

| コード | 説明 |
|--------|------|
| `VALIDATION_FAILED` | 入力内容に不備があります |
| `REQUIRED_FIELD_MISSING` | 必須項目が入力されていません |
| `INVALID_FORMAT` | 入力形式が正しくありません |
| `USER_ALREADY_EXISTS` | 既に登録されているメールアドレスです |
| `SESSION_EXPIRED` | セッションが期限切れです |
| `CSRF_TOKEN_INVALID` | CSRFトークンが無効です |
| `RATE_LIMIT_EXCEEDED` | アクセス数が上限に達しました |
| `INTERNAL_SERVER_ERROR` | サーバーエラーが発生しました |

## エンドポイント

### ヘルスチェック

#### GET /health

サービスの稼働状況を確認します。

**レスポンス**

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00Z",
    "version": "1.0.0",
    "services": {
      "database": "healthy",
      "external_apis": "healthy"
    }
  }
}
```

#### GET /health/live

Kubernetes Liveness Probe用のエンドポイントです。

#### GET /health/ready

Kubernetes Readiness Probe用のエンドポイントです。

### セキュリティ

#### GET /api/v1/csrf-token

CSRFトークンを取得します。

**レスポンス**

```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9..."
  }
}
```

### ユーザー管理

#### POST /api/v1/users

新しいユーザーを登録します。

**リクエストヘッダー**

```
Content-Type: application/json
X-CSRF-Token: {token}
```

**リクエストボディ**

```json
{
  "last_name": "田中",
  "first_name": "太郎",
  "last_name_kana": "タナカ",
  "first_name_kana": "タロウ",
  "phone1": "03",
  "phone2": "1234",
  "phone3": "5678",
  "postal_code1": "100",
  "postal_code2": "0001",
  "prefecture": "東京都",
  "city": "千代田区",
  "town": "丸の内",
  "chome": "1",
  "banchi": "1-1",
  "go": "101",
  "building": "東京ビル",
  "room": "1001",
  "email": "taro@example.com",
  "plan_type": "A",
  "option_types": ["AA", "AB"]
}
```

**レスポンス**

```json
{
  "success": true,
  "data": {
    "id": 123,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### POST /api/v1/users/validate

ユーザーデータのバリデーションを実行します。

**リクエスト形式**: `/api/v1/users`と同じ

**レスポンス**

```json
{
  "success": true,
  "data": {
    "valid": true,
    "errors": {}
  }
}
```

エラーがある場合:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "入力内容に不備があります",
    "details": {
      "email": "メールアドレスの形式が正しくありません",
      "phone1": "市外局番は2桁以上で入力してください"
    }
  }
}
```

### セッション管理

#### POST /api/v1/sessions

一時保存用のセッションを作成します。

**リクエストボディ**

```json
{
  "form_data": {
    "last_name": "田中",
    "first_name": "太郎",
    // ... その他のフォームデータ
  },
  "current_step": "input",
  "expires_at": "2024-01-15T14:30:00Z"
}
```

**レスポンス**

```json
{
  "success": true,
  "data": {
    "session_id": "sess_abc123def456",
    "expires_at": "2024-01-15T14:30:00Z"
  }
}
```

#### GET /api/v1/sessions/{session_id}

セッションデータを取得します。

**レスポンス**

```json
{
  "success": true,
  "data": {
    "session_id": "sess_abc123def456",
    "form_data": {
      "last_name": "田中",
      "first_name": "太郎"
      // ... その他のフォームデータ
    },
    "current_step": "input",
    "created_at": "2024-01-15T10:30:00Z",
    "expires_at": "2024-01-15T14:30:00Z"
  }
}
```

#### PUT /api/v1/sessions/{session_id}

セッションデータを更新します。

**リクエスト形式**: `/api/v1/sessions`と同じ

#### DELETE /api/v1/sessions/{session_id}

セッションを削除します。

**レスポンス**

```json
{
  "success": true
}
```

### マスターデータ

#### GET /api/v1/prefectures

都道府県一覧を取得します。

**レスポンス**

```json
{
  "success": true,
  "data": {
    "prefectures": [
      {
        "code": "01",
        "name": "北海道"
      },
      {
        "code": "13",
        "name": "東京都"
      }
      // ... 他の都道府県
    ]
  }
}
```

#### GET /api/v1/plans

プラン一覧を取得します。

**レスポンス**

```json
{
  "success": true,
  "data": {
    "plans": [
      {
        "type": "A",
        "name": "Aプラン",
        "description": "基本プラン",
        "price": 1000
      },
      {
        "type": "B", 
        "name": "Bプラン",
        "description": "プレミアムプラン",
        "price": 2000
      }
    ]
  }
}
```

#### GET /api/v1/options

オプション一覧を取得します。

**クエリパラメータ**

- `plan_type`: プランタイプ（A または B）

**レスポンス**

```json
{
  "success": true,
  "data": {
    "options": [
      {
        "type": "AA",
        "name": "AAオプション",
        "description": "Aプラン専用オプション",
        "price": 500,
        "available_plans": ["A"]
      },
      {
        "type": "AB",
        "name": "ABオプション", 
        "description": "共通オプション",
        "price": 300,
        "available_plans": ["A", "B"]
      }
    ]
  }
}
```

### 外部API連携

#### POST /api/v1/options/check-inventory

在庫状況を確認します。

**リクエストボディ**

```json
{
  "option_types": ["AA", "BB", "AB"]
}
```

**レスポンス**

```json
{
  "success": true,
  "data": {
    "inventory": {
      "AA": 10,
      "BB": 0,
      "AB": 5
    },
    "checked_at": "2024-01-15T10:30:00Z"
  }
}
```

#### GET /api/v1/address/search

郵便番号から住所を検索します。

**クエリパラメータ**

- `postal_code`: 郵便番号（ハイフンなし 7桁）

**レスポンス**

```json
{
  "success": true,
  "data": {
    "postal_code": "1000001",
    "prefecture": "東京都",
    "city": "千代田区",
    "town": "丸の内"
  }
}
```

住所が見つからない場合:

```json
{
  "success": false,
  "error": {
    "code": "ADDRESS_NOT_FOUND",
    "message": "入力された郵便番号では住所が見つかりません"
  }
}
```

#### POST /api/v1/region/check

地域制限を確認します。

**リクエストボディ**

```json
{
  "prefecture": "東京都",
  "city": "千代田区",
  "option_types": ["AA", "BB"]
}
```

**レスポンス**

```json
{
  "success": true,
  "data": {
    "availability": {
      "AA": true,
      "BB": false
    },
    "checked_at": "2024-01-15T10:30:00Z"
  }
}
```

## レート制限

- **制限**: 100リクエスト/分/IP
- **制限時のレスポンス**: HTTP 429 Too Many Requests
- **ヘッダー**:
  - `X-RateLimit-Limit`: 制限値
  - `X-RateLimit-Window`: 時間窓（秒）
  - `Retry-After`: 再試行可能時間（秒）

## セキュリティ

### CSRF保護

- すべてのPOST、PUT、DELETEリクエストでCSRFトークンが必要
- トークンは`X-CSRF-Token`ヘッダーで送信
- トークンの有効期限は4時間

### セキュリティヘッダー

- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy: default-src 'self'`
- `Strict-Transport-Security: max-age=31536000`

### 入力バリデーション

- すべての入力データは型チェック、長さチェック、形式チェックを実施
- 疑わしいパターン（スクリプトタグ等）は拒否
- SQLインジェクション、XSS対策を実装

## パフォーマンス

### レスポンス時間

- 通常のAPIコール: 500ms以下
- 外部API連携: 2秒以下

### キャッシュ

- マスターデータ（都道府県、プラン）: 1時間
- セッションデータ: 4時間

## 監視・ログ

### メトリクス

- リクエスト数、レスポンス時間、エラー率
- データベース接続数、クエリ実行時間
- 外部API連携の成功率、レスポンス時間

### ログ形式

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "message": "Request processed",
  "request_id": "req_abc123",
  "method": "POST",
  "path": "/api/v1/users",
  "status_code": 201,
  "response_time_ms": 245,
  "user_agent": "Mozilla/5.0...",
  "ip": "192.168.1.100"
}
```

## 環境固有設定

### 本番環境

- ベースURL: `https://api.normal-form-app.com`
- ログレベル: `warn`
- レート制限: 100req/min

### 開発環境

- ベースURL: `https://dev-api.normal-form-app.com`
- ログレベル: `info`
- レート制限: 1000req/min

### テスト環境

- ベースURL: `https://test-api.normal-form-app.com`
- ログレベル: `debug`
- レート制限: 500req/min