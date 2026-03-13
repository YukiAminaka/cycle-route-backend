# Terraform デプロイ手順

## 前提条件

- AWS CLI がインストール・設定済み
- Terraform >= 1.0 がインストール済み
- Docker がインストール済み（コンテナイメージのビルド用）

## 初期セットアップ

### 1. AWS 認証情報の設定

```bash
aws configure
```

### 2. Terraform の初期化

```bash
cd terraform/prod
terraform init
```

### 3. 変数の確認・カスタマイズ（必要に応じて）

`terraform/prod/variables.tf` を編集してリージョンやプロジェクト名を変更できます。

## デプロイ手順

### 1. インフラのプロビジョニング

```bash
cd terraform/prod

# 実行計画の確認
terraform plan

# インフラの作成
terraform apply
```

デプロイには約10-15分かかります。

### 2. ECR リポジトリ URL の取得

```bash
terraform output ecr_repository_urls
```

### 3. Docker イメージのビルドとプッシュ

#### ECR にログイン

```bash
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin <AWS_ACCOUNT_ID>.dkr.ecr.ap-northeast-1.amazonaws.com
```

#### API イメージのビルドとプッシュ

```bash
cd ../../  # プロジェクトルートへ

# イメージのビルド
docker build -t cycle-route-api .

# タグ付け
docker tag cycle-route-api:latest <ECR_API_URL>:latest

# プッシュ
docker push <ECR_API_URL>:latest
```

#### Kratos イメージのプッシュ

```bash
# Kratos の公式イメージを取得してプッシュ
docker pull oryd/kratos:v25.4.0
docker tag oryd/kratos:v25.4.0 <ECR_KRATOS_URL>:latest
docker push <ECR_KRATOS_URL>:latest
```

#### Frontend イメージのビルドとプッシュ

```bash
# フロントエンドのディレクトリで実行
cd ../frontend  # フロントエンドのパスに応じて調整

docker build -t cycle-route-frontend .
docker tag cycle-route-frontend:latest <ECR_FRONTEND_URL>:latest
docker push <ECR_FRONTEND_URL>:latest
```

### 4. ECS サービスの更新

イメージをプッシュ後、ECS サービスが自動的に新しいイメージを取得します。
手動で更新する場合：

```bash
aws ecs update-service \
  --cluster cycle-route-prod-cluster \
  --service cycle-route-prod-api \
  --force-new-deployment \
  --region ap-northeast-1
```

### 5. データベースのマイグレーション

ECS タスクに接続してマイグレーションを実行：

```bash
# タスク ID を取得
aws ecs list-tasks \
  --cluster cycle-route-prod-cluster \
  --service-name cycle-route-prod-api \
  --region ap-northeast-1

# ECS Exec で接続
aws ecs execute-command \
  --cluster cycle-route-prod-cluster \
  --task <TASK_ID> \
  --container api \
  --interactive \
  --command "/bin/sh"

# コンテナ内でマイグレーション実行
atlas migrate apply --env prod
```

### 6. アプリケーションへのアクセス

```bash
# ALB の DNS 名を取得
terraform output alb_dns_name
```

ブラウザで `http://<ALB_DNS_NAME>` にアクセス

## 主要な出力値

```bash
# すべての出力値を表示
terraform output

# 特定の出力値
terraform output alb_dns_name
terraform output ecr_repository_urls
terraform output db_endpoint
```

## リソースの削除

```bash
cd terraform/prod
terraform destroy
```

## トラブルシューティング

### ECS タスクが起動しない

```bash
# ログの確認
aws logs tail /ecs/cycle-route-prod/api --follow
```

### データベース接続エラー

- セキュリティグループの設定を確認
- RDS エンドポイントが正しいか確認
- Secrets Manager のパスワードが正しいか確認

### イメージのプル失敗

- ECR にイメージがプッシュされているか確認
- ECS タスク実行ロールに ECR へのアクセス権限があるか確認

## 本番環境の推奨設定

現在の構成は最小限のリソースです。本番環境では以下を検討してください：

1. **HTTPS の有効化**: ACM で証明書を取得し、ALB に設定
2. **RDS のスケールアップ**: `db.t3.micro` → `db.t3.small` 以上
3. **ECS タスク数の増加**: `desired_count = 1` → `2` 以上
4. **Auto Scaling の設定**: CPU/メモリ使用率に基づく自動スケーリング
5. **CloudFront の追加**: 静的コンテンツのキャッシュ
6. **Route53 の設定**: カスタムドメインの設定
7. **バックアップの設定**: RDS の自動バックアップ保持期間の延長
8. **監視・アラート**: CloudWatch Alarms の設定
9. **WAF の追加**: セキュリティ強化
10. **Terraform State の S3 バックエンド化**: チーム開発用

## 注意事項

- **Kratos の DSN**: 本番環境では SQLite ではなく PostgreSQL を使用するよう設定済み
- **Secrets**: パスワードは Secrets Manager で自動生成・管理
- **コスト**: NAT Gateway が最もコストがかかります（約$32/月 × 2）
- **PostGIS**: RDS PostgreSQL で PostGIS 拡張を有効化する必要があります

```sql
-- RDS に接続後
CREATE EXTENSION IF NOT EXISTS postgis;
```
