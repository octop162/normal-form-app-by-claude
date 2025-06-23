# デプロイメントガイド

## 概要

会員登録webフォームアプリケーションのデプロイメント手順書です。

## 前提条件

### 必要なツール

- AWS CLI v2.x
- Docker v20.x以上
- Terraform v1.0以上
- kubectl v1.25以上（必要に応じて）
- Node.js v18.x
- Go v1.21.x

### AWS権限

デプロイに必要なAWS権限:

- ECS Full Access
- ECR Full Access
- VPC Read Access
- IAM Limited Access
- CloudWatch Full Access
- Secrets Manager Read/Write
- RDS Limited Access

## 環境別設定

### 本番環境 (production)

- **ECSクラスター**: `normal-form-app-cluster`
- **ECRリポジトリ**: 
  - `normal-form-app-backend`
  - `normal-form-app-frontend`
- **RDS**: Aurora PostgreSQL
- **ロードバランサー**: Application Load Balancer
- **ドメイン**: `https://normal-form-app.com`

### ステージング環境 (staging)

- **ECSクラスター**: `normal-form-app-cluster-staging`
- **ECRリポジトリ**: 共通（tagで区別）
- **RDS**: Aurora PostgreSQL (dev instance)
- **ドメイン**: `https://staging.normal-form-app.com`

## デプロイ手順

### 1. 事前準備

#### 1.1 環境変数設定

```bash
export AWS_REGION=ap-northeast-1
export AWS_ACCOUNT_ID=123456789012
export ENVIRONMENT=production  # or staging
```

#### 1.2 AWS認証

```bash
aws configure
# または
aws sso login --profile your-profile
```

#### 1.3 ECRログイン

```bash
aws ecr get-login-password --region $AWS_REGION | \
  docker login --username AWS --password-stdin \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com
```

### 2. Docker イメージビルド

#### 2.1 バックエンドイメージ

```bash
# プロジェクトルートで実行
docker build -t normal-form-app-backend:latest -f Dockerfile.backend .

# ECRにタグ付け
docker tag normal-form-app-backend:latest \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:latest

# ECRにプッシュ
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-backend:latest
```

#### 2.2 フロントエンドイメージ

```bash
# frontend ディレクトリで実行
cd frontend
docker build -t normal-form-app-frontend:latest -f Dockerfile --target production .

# ECRにタグ付け
docker tag normal-form-app-frontend:latest \
  $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:latest

# ECRにプッシュ
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/normal-form-app-frontend:latest
```

### 3. ECSデプロイ

#### 3.1 タスク定義更新

```bash
# タスク定義ファイルの環境変数を置換
envsubst < deployments/ecs-task-definition.json > task-definition-updated.json

# タスク定義を登録
aws ecs register-task-definition \
  --cli-input-json file://task-definition-updated.json \
  --region $AWS_REGION
```

#### 3.2 サービス更新

```bash
# サービスを更新（新しいタスク定義を使用）
aws ecs update-service \
  --cluster normal-form-app-cluster \
  --service normal-form-app-service \
  --force-new-deployment \
  --region $AWS_REGION

# デプロイ完了を待機
aws ecs wait services-stable \
  --cluster normal-form-app-cluster \
  --services normal-form-app-service \
  --region $AWS_REGION
```

### 4. データベースマイグレーション

#### 4.1 マイグレーション実行

```bash
# マイグレーションタスクを実行
aws ecs run-task \
  --cluster normal-form-app-cluster \
  --task-definition normal-form-app-migration \
  --launch-type FARGATE \
  --network-configuration '{
    "awsvpcConfiguration": {
      "subnets": ["subnet-xxx", "subnet-yyy"],
      "securityGroups": ["sg-xxx"],
      "assignPublicIp": "DISABLED"
    }
  }' \
  --region $AWS_REGION
```

### 5. デプロイ後確認

#### 5.1 ヘルスチェック

```bash
# ALB DNSを取得
ALB_DNS=$(aws elbv2 describe-load-balancers \
  --names normal-form-app-alb \
  --query 'LoadBalancers[0].DNSName' \
  --output text \
  --region $AWS_REGION)

# ヘルスチェック実行
curl -f "https://$ALB_DNS/health"
```

#### 5.2 ログ確認

```bash
# ECSタスクのログを確認
aws logs tail /ecs/normal-form-app/backend --follow --region $AWS_REGION
```

## 自動デプロイ (CI/CD)

### GitHub Actions

`.github/workflows/ci-cd.yml` ファイルで自動デプロイが設定されています。

#### 必要なSecrets

- `AWS_ACCESS_KEY_ID`: AWS アクセスキー
- `AWS_SECRET_ACCESS_KEY`: AWS シークレットキー
- `AWS_ACCOUNT_ID`: AWSアカウントID
- `SLACK_WEBHOOK_URL`: Slack通知用（オプション）

#### デプロイトリガー

- `main` ブランチへのプッシュ
- プルリクエストのマージ

### 手動デプロイスクリプト

```bash
# スクリプトを使用したデプロイ
./scripts/deploy.sh production deploy

# ロールバック
./scripts/deploy.sh production rollback
```

## トラブルシューティング

### よくある問題

#### 1. ECSタスクが起動しない

**原因**: 
- Docker イメージが見つからない
- 環境変数の設定ミス
- セキュリティグループの設定

**対処法**:
```bash
# ECSタスクの詳細を確認
aws ecs describe-tasks \
  --cluster normal-form-app-cluster \
  --tasks TASK_ARN \
  --region $AWS_REGION

# CloudWatch Logs を確認
aws logs get-log-events \
  --log-group-name /ecs/normal-form-app/backend \
  --log-stream-name ecs/backend/TASK_ID \
  --region $AWS_REGION
```

#### 2. ヘルスチェックが失敗する

**原因**:
- アプリケーションの起動が完了していない
- データベース接続エラー
- セキュリティグループの制限

**対処法**:
```bash
# ターゲットグループの状態確認
aws elbv2 describe-target-health \
  --target-group-arn TARGET_GROUP_ARN \
  --region $AWS_REGION

# アプリケーションログ確認
aws logs filter-log-events \
  --log-group-name /ecs/normal-form-app/backend \
  --filter-pattern "ERROR" \
  --region $AWS_REGION
```

#### 3. データベース接続エラー

**原因**:
- Secrets Manager の設定ミス
- VPC/セキュリティグループの設定
- RDS インスタンスの状態

**対処法**:
```bash
# Secrets Manager の確認
aws secretsmanager get-secret-value \
  --secret-id normal-form-app/db-password \
  --region $AWS_REGION

# RDS の状態確認
aws rds describe-db-clusters \
  --db-cluster-identifier normal-form-app-aurora-cluster \
  --region $AWS_REGION
```

### ロールバック手順

#### 1. 自動ロールバック

```bash
./scripts/deploy.sh production rollback
```

#### 2. 手動ロールバック

```bash
# 前のタスク定義リビジョンを確認
aws ecs list-task-definitions \
  --family-prefix normal-form-app \
  --region $AWS_REGION

# 前のリビジョンでサービス更新
aws ecs update-service \
  --cluster normal-form-app-cluster \
  --service normal-form-app-service \
  --task-definition normal-form-app:PREVIOUS_REVISION \
  --region $AWS_REGION
```

## セキュリティ考慮事項

### 1. Secrets管理

- データベース認証情報は AWS Secrets Manager を使用
- API キーは環境変数として設定
- 本番環境では IAM ロールベースの認証を使用

### 2. ネットワークセキュリティ

- ECS タスクはプライベートサブネットで実行
- データベースはプライベートサブネットに配置
- ALB のみパブリックサブネットに配置

### 3. 監査ログ

- CloudTrail でAPI呼び出しをログ記録
- VPC Flow Logs でネットワーク通信を監視
- CloudWatch Logs でアプリケーションログを集約

## モニタリング・アラート

### CloudWatch メトリクス

- ECS CPU/メモリ使用率
- ALB レスポンス時間・エラー率
- RDS 接続数・CPU使用率

### アラート設定

- 高CPU使用率 (80%以上)
- 高エラー率 (5%以上)
- ヘルスチェック失敗
- データベース接続数上限

### ダッシュボード

Grafana ダッシュボードで以下を監視:

- リクエスト数・レスポンス時間
- エラー率・成功率
- インフラリソース使用状況
- ビジネスメトリクス

## バックアップ・復旧

### 自動バックアップ

- RDS 自動バックアップ: 7日間保持
- AWS Backup: 日次バックアップ
- クロスリージョンレプリケーション

### 復旧手順

```bash
# データベース復旧
aws rds restore-db-cluster-from-snapshot \
  --db-cluster-identifier restored-cluster \
  --snapshot-identifier SNAPSHOT_ID \
  --region $AWS_REGION

# アプリケーション復旧
./scripts/deploy.sh production deploy
```

## パフォーマンス最適化

### ECS設定

- CPU: 1024 (1 vCPU)
- メモリ: 2048 MB
- Auto Scaling: CPU 70%でスケールアウト

### データベース最適化

- 接続プール: 最大25接続
- クエリ最適化: インデックス設定
- 定期メンテナンス: 週次

### CDN設定

- CloudFront でフロントエンドアセットをキャッシュ
- 画像・CSS・JavaScript を最適化
- Gzip 圧縮を有効化

## 定期メンテナンス

### 月次作業

- セキュリティパッチ適用
- パフォーマンス監視レポート作成
- バックアップ復旧テスト実施

### 四半期作業

- インフラコスト最適化レビュー
- セキュリティ監査実施
- 災害復旧訓練実施