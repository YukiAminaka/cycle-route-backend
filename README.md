# 1. スキーマ編集

vim schema.sql

# 2. sqlc でコード生成

sqlc generate

# 3. Docker 上の DB に反映

docker compose exec -T postgres psql -U postgres -d postgres_db < schema.sql
