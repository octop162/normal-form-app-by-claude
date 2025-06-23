# 運用手順書

## 概要

会員登録webフォームアプリケーションの日常運用に関する手順書です。

## 運用体制

### 役割分担

| 役割 | 責任範囲 | 連絡先 |
|------|----------|--------|
| 運用責任者 | 全体統括、エスカレーション判断 | ops-manager@company.com |
| インフラエンジニア | AWS、ネットワーク、セキュリティ | infrastructure@company.com |
| アプリケーションエンジニア | アプリ不具合、デプロイ | developers@company.com |
| 監視担当 | 24/7監視、初期対応 | monitoring@company.com |

### 連絡体制

#### 営業時間内 (平日 9:00-18:00)
1. 監視担当が第一報
2. 担当エンジニアが詳細調査
3. 必要に応じて運用責任者にエスカレーション

#### 営業時間外・休日
1. 監視システムの自動アラート
2. オンコール担当者が対応
3. 重大インシデントは運用責任者に即座に連絡

## 日常監視項目

### システム監視

#### 1. アプリケーション監視

**確認項目**:
- ヘルスチェックエンドポイントの応答
- レスポンス時間 (95パーセンタイル < 2秒)
- エラー率 (< 1%)
- アクティブセッション数

**確認方法**:
```bash
# ヘルスチェック
curl -f https://normal-form-app.com/health

# CloudWatch メトリクス確認
aws cloudwatch get-metric-statistics \
  --namespace "Normal-Form-App/Application" \
  --metric-name "ErrorRate" \
  --start-time 2024-01-15T09:00:00Z \
  --end-time 2024-01-15T10:00:00Z \
  --period 300 \
  --statistics Average
```

**閾値**:
- レスポンス時間: 2秒以上で警告
- エラー率: 1%以上で警告、5%以上で重要
- CPU使用率: 80%以上で警告

#### 2. インフラ監視

**ECS監視**:
```bash
# ECSサービス状態確認
aws ecs describe-services \
  --cluster normal-form-app-cluster \
  --services normal-form-app-service

# タスク状態確認
aws ecs list-tasks \
  --cluster normal-form-app-cluster \
  --service-name normal-form-app-service
```

**RDS監視**:
```bash
# RDS状態確認
aws rds describe-db-clusters \
  --db-cluster-identifier normal-form-app-aurora-cluster

# 接続数確認
aws cloudwatch get-metric-statistics \
  --namespace "AWS/RDS" \
  --metric-name "DatabaseConnections" \
  --dimensions Name=DBClusterIdentifier,Value=normal-form-app-aurora-cluster \
  --start-time 2024-01-15T09:00:00Z \
  --end-time 2024-01-15T10:00:00Z \
  --period 300 \
  --statistics Average
```

#### 3. セキュリティ監視

**確認項目**:
- 不正アクセスの検知
- CSRFトークンエラー率
- レート制限の発動状況
- SSL証明書の有効期限

**確認方法**:
```bash
# セキュリティアラートの確認
aws logs filter-log-events \
  --log-group-name /ecs/normal-form-app/backend \
  --filter-pattern "SECURITY"

# SSL証明書の確認
echo | openssl s_client -servername normal-form-app.com \
  -connect normal-form-app.com:443 2>/dev/null | \
  openssl x509 -noout -dates
```

### 業務監視

#### 1. ユーザー登録状況

**確認項目**:
- 日次登録数
- 登録完了率
- エラー発生パターン

**確認方法**:
```sql
-- 日次登録数
SELECT DATE(created_at), COUNT(*) 
FROM users 
WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
GROUP BY DATE(created_at)
ORDER BY DATE(created_at);

-- 登録完了率（セッションベース）
SELECT 
  COUNT(CASE WHEN users.id IS NOT NULL THEN 1 END) as completed,
  COUNT(*) as total_sessions,
  ROUND(COUNT(CASE WHEN users.id IS NOT NULL THEN 1 END) * 100.0 / COUNT(*), 2) as completion_rate
FROM user_sessions 
LEFT JOIN users ON users.email = user_sessions.user_data->>'email'
WHERE user_sessions.created_at >= CURRENT_DATE - INTERVAL '1 day';
```

#### 2. 外部API連携状況

**確認項目**:
- 在庫API応答時間・成功率
- 住所検索API応答時間・成功率
- 地域制限API応答時間・成功率

**確認方法**:
```bash
# API応答時間の確認
aws cloudwatch get-metric-statistics \
  --namespace "Normal-Form-App/External-API" \
  --metric-name "ResponseTime" \
  --dimensions Name=APIName,Value=InventoryAPI \
  --start-time 2024-01-15T09:00:00Z \
  --end-time 2024-01-15T10:00:00Z \
  --period 300 \
  --statistics Average
```

## 障害対応手順

### 障害レベル定義

| レベル | 定義 | 対応時間 | エスカレーション |
|--------|------|----------|------------------|
| P1 (緊急) | サービス全停止 | 15分以内 | 即座に運用責任者 |
| P2 (重要) | 一部機能停止、性能劣化 | 1時間以内 | 30分後に運用責任者 |
| P3 (軽微) | 軽微な不具合 | 4時間以内 | 定期報告 |

### P1障害対応手順

#### 1. 初期対応 (0-15分)

1. **障害確認**
   ```bash
   # サービス状態確認
   curl -f https://normal-form-app.com/health
   
   # ECS状態確認
   aws ecs describe-services \
     --cluster normal-form-app-cluster \
     --services normal-form-app-service
   ```

2. **ステークホルダー通知**
   - 運用責任者に即座に連絡
   - 社内チャットで障害発生を通知
   - 顧客向けお知らせページを更新（必要に応じて）

3. **影響範囲特定**
   ```bash
   # エラーログ確認
   aws logs filter-log-events \
     --log-group-name /ecs/normal-form-app/backend \
     --filter-pattern "ERROR" \
     --start-time $(date -d '15 minutes ago' +%s)000
   ```

#### 2. 調査・復旧作業 (15分-1時間)

1. **ログ分析**
   ```bash
   # 詳細ログ確認
   aws logs get-log-events \
     --log-group-name /ecs/normal-form-app/backend \
     --log-stream-name STREAM_NAME \
     --start-time $(date -d '1 hour ago' +%s)000
   ```

2. **リソース状況確認**
   ```bash
   # CPU/メモリ使用率確認
   aws cloudwatch get-metric-statistics \
     --namespace "AWS/ECS" \
     --metric-name "CPUUtilization" \
     --dimensions Name=ServiceName,Value=normal-form-app-service Name=ClusterName,Value=normal-form-app-cluster \
     --start-time $(date -d '1 hour ago' -u +%Y-%m-%dT%H:%M:%S) \
     --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
     --period 300 \
     --statistics Average
   ```

3. **復旧作業**
   
   **アプリケーション再起動**:
   ```bash
   # ECSサービス再起動
   aws ecs update-service \
     --cluster normal-form-app-cluster \
     --service normal-form-app-service \
     --force-new-deployment
   ```
   
   **ロールバック**:
   ```bash
   # 前バージョンにロールバック
   ./scripts/deploy.sh production rollback
   ```

#### 3. 復旧確認・報告 (1時間以降)

1. **サービス復旧確認**
   ```bash
   # ヘルスチェック
   for i in {1..5}; do
     curl -f https://normal-form-app.com/health && echo " - OK" || echo " - FAIL"
     sleep 10
   done
   ```

2. **障害報告書作成**
   - 発生時刻、復旧時刻
   - 原因分析
   - 影響範囲
   - 再発防止策

### P2障害対応手順

#### 1. 性能劣化対応

1. **現状把握**
   ```bash
   # レスポンス時間確認
   aws cloudwatch get-metric-statistics \
     --namespace "AWS/ApplicationELB" \
     --metric-name "TargetResponseTime" \
     --dimensions Name=LoadBalancer,Value=ALB_NAME \
     --start-time $(date -d '30 minutes ago' -u +%Y-%m-%dT%H:%M:%S) \
     --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
     --period 300 \
     --statistics Average
   ```

2. **リソーススケーリング**
   ```bash
   # ECSタスク数増加
   aws ecs update-service \
     --cluster normal-form-app-cluster \
     --service normal-form-app-service \
     --desired-count 4
   ```

3. **データベース最適化**
   ```sql
   -- 長時間実行中のクエリ確認
   SELECT pid, now() - pg_stat_activity.query_start AS duration, query 
   FROM pg_stat_activity 
   WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
   
   -- 必要に応じてクエリ停止
   SELECT pg_terminate_backend(pid);
   ```

#### 2. 部分機能停止対応

1. **機能別影響確認**
   ```bash
   # API別エラー率確認
   aws logs filter-log-events \
     --log-group-name /ecs/normal-form-app/backend \
     --filter-pattern "[timestamp, level=ERROR, message, path]"
   ```

2. **外部API連携確認**
   ```bash
   # 外部API疎通確認
   curl -f -m 10 "$INVENTORY_API_URL/health"
   curl -f -m 10 "$ADDRESS_API_URL/health"
   ```

3. **フェイルオーバー実行**
   - 外部API障害時は代替手段（手動入力）を提供
   - データベース障害時は読み取り専用レプリカに切り替え

## 定期作業

### 日次作業

#### 1. 監視ダッシュボード確認

**確認項目**:
- サービス稼働状況
- 主要メトリクス
- エラーログサマリー
- 前日の登録数

**実行時間**: 毎朝9:00

**チェックリスト**:
- [ ] Grafanaダッシュボードでメトリクス確認
- [ ] CloudWatchアラーム状況確認
- [ ] エラーログで新しい問題がないか確認
- [ ] ユーザー登録数が正常範囲内か確認

#### 2. バックアップ状況確認

```bash
# AWS Backup状況確認
aws backup list-backup-jobs \
  --by-creation-date-after $(date -d '1 day ago' +%Y-%m-%d) \
  --by-creation-date-before $(date +%Y-%m-%d)

# RDSスナップショット確認
aws rds describe-db-cluster-snapshots \
  --db-cluster-identifier normal-form-app-aurora-cluster \
  --snapshot-type automated
```

### 週次作業

#### 1. セキュリティ状況確認

**実行時間**: 毎週月曜日10:00

**作業内容**:
```bash
# セキュリティアラート確認
aws logs filter-log-events \
  --log-group-name /ecs/normal-form-app/backend \
  --filter-pattern "SECURITY" \
  --start-time $(date -d '7 days ago' +%s)000

# SSL証明書有効期限確認
echo | openssl s_client -servername normal-form-app.com \
  -connect normal-form-app.com:443 2>/dev/null | \
  openssl x509 -noout -dates

# 不正アクセス試行の確認
aws logs filter-log-events \
  --log-group-name /aws/applicationloadbalancer/normal-form-app-alb \
  --filter-pattern "[..., status_code=4*, ...]" \
  --start-time $(date -d '7 days ago' +%s)000
```

#### 2. 性能監視レポート作成

**実行時間**: 毎週金曜日16:00

**レポート項目**:
- 週間リクエスト数・レスポンス時間
- エラー率推移
- リソース使用率
- ユーザー登録数推移

### 月次作業

#### 1. 容量計画見直し

**実行時間**: 毎月第1営業日

**確認項目**:
- ECSタスク数・リソース使用率
- RDS容量・性能
- ログ保存容量
- S3ストレージ使用量

```bash
# RDS容量確認
aws cloudwatch get-metric-statistics \
  --namespace "AWS/RDS" \
  --metric-name "FreeStorageSpace" \
  --dimensions Name=DBClusterIdentifier,Value=normal-form-app-aurora-cluster \
  --start-time $(date -d '30 days ago' -u +%Y-%m-%dT%H:%M:%S) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%S) \
  --period 86400 \
  --statistics Average
```

#### 2. セキュリティパッチ適用

**実行時間**: 毎月第2土曜日（メンテナンス時間）

**作業手順**:
1. Docker基盤イメージの更新
2. 依存ライブラリの更新
3. セキュリティテスト実行
4. ステージング環境でのテスト
5. 本番環境への適用

#### 3. 災害復旧テスト

**実行時間**: 毎月第3日曜日

**テスト内容**:
- バックアップからの復旧テスト
- 代替リージョンでの起動テスト
- データ整合性確認

## アラート対応

### CloudWatchアラート

#### 高CPU使用率アラート

**対応手順**:
1. 現在のタスク数確認
2. スケールアウト実行
3. 原因調査（重いクエリ、外部API遅延等）

```bash
# スケールアウト
aws ecs update-service \
  --cluster normal-form-app-cluster \
  --service normal-form-app-service \
  --desired-count $((CURRENT_COUNT + 2))
```

#### 高エラー率アラート

**対応手順**:
1. エラーログ確認
2. 外部API状況確認
3. 必要に応じてフェイルオーバー実行

```bash
# エラー分析
aws logs filter-log-events \
  --log-group-name /ecs/normal-form-app/backend \
  --filter-pattern "ERROR" \
  --start-time $(date -d '10 minutes ago' +%s)000 | \
  jq '.events[] | .message' | \
  sort | uniq -c | sort -nr
```

### Grafanaアラート

#### データベース接続数上限

**対応手順**:
1. 接続プール設定確認
2. 長時間実行クエリの確認
3. 接続リーク調査

```sql
-- アクティブ接続確認
SELECT state, count(*) 
FROM pg_stat_activity 
GROUP BY state;

-- 長時間接続確認
SELECT pid, usename, application_name, state, 
       now() - state_change as duration
FROM pg_stat_activity 
WHERE state != 'idle'
ORDER BY duration DESC;
```

### 外部通知連携

#### Slack通知

アラート発生時の自動通知内容:
- アラート名・重要度
- 発生時刻
- 影響サービス
- 推奨対応アクション

#### メール通知

重大アラート（P1/P2）時の通知先:
- 運用責任者
- 当直エンジニア
- インフラチーム

## 業務継続計画 (BCP)

### 災害レベル定義

| レベル | 状況 | 対応策 |
|--------|------|--------|
| L1 | 単一AZ障害 | 自動フェイルオーバー |
| L2 | リージョン障害 | 手動クロスリージョン復旧 |
| L3 | AWS全体障害 | 代替クラウド起動 |

### 復旧目標

- **RTO (復旧時間目標)**: 1時間以内
- **RPO (復旧ポイント目標)**: 15分以内

### 代替運用手順

#### 外部API障害時

1. **在庫確認API障害**:
   - 在庫状況を「要確認」として表示
   - 管理者による手動確認フローに切り替え

2. **住所検索API障害**:
   - 郵便番号入力を無効化
   - 住所の手動入力を必須化

3. **地域制限API障害**:
   - 全地域で申し込み受付
   - 後処理で地域制限を確認

## 連絡先・エスカレーション

### 緊急連絡先

| 役割 | 氏名 | 電話番号 | メール | 対応時間 |
|------|------|----------|--------|----------|
| 運用責任者 | 田中太郎 | 090-1234-5678 | tanaka@company.com | 24/7 |
| オンコール（平日） | 佐藤花子 | 080-9876-5432 | sato@company.com | 平日9-18時 |
| オンコール（休日） | 鈴木一郎 | 070-1111-2222 | suzuki@company.com | 休日・夜間 |

### エスカレーション基準

1. **15分以内に初期対応完了しない場合**
   → 運用責任者にエスカレーション

2. **1時間以内に復旧見込みが立たない場合**
   → 経営陣に報告

3. **データ損失の可能性がある場合**
   → 即座に最高責任者に報告

### 外部ベンダー連絡先

| サービス | 担当者 | 連絡先 | 備考 |
|----------|--------|--------|------|
| AWS Enterprise Support | - | +81-3-4578-4037 | 24/7サポート |
| 外部API提供者A | 山田 | api-support@vendor-a.com | 平日10-17時 |
| 外部API提供者B | 田村 | emergency@vendor-b.com | 24/7 |