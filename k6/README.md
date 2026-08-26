## テスト実施方法

### 使用するファイル

- k6/seed-users.sh

Ory Kratos の PATCH /admin/identities でユーザーを一括登録
登録したユーザーを user.csv に書き出す
200件ずつバッチ処理（プレーンテキストPWの上限に合わせる）

シードスクリプトの処理フロー

```
seed-users.sh
  │
  ├─ [Step 1] PATCH /admin/identities  (Kratos Admin API)
  │    └─ 最大200件ずつバッチ登録
  │    └─ レスポンスに kratos_id が含まれる
  │
  └─ [Step 2] POST /api/v1/users  (アプリAPI)
       └─ kratos_id + name + email を渡して
          アプリDBにユーザーを作成
       └─ 成功したユーザーのみ user.csv に追記
```

patch_id を使った ID 連携の仕組み
Kratosのバッチレスポンスは patch_id と identity(Kratos ID) を返します。シードスクリプトでは patch_id の先頭8桁にユーザーの連番を埋め込んでいるので、レスポンス処理時にメールアドレスとパスワードを再導出できます。

```
patch_id: "00000042-0000-0000-0000-000000000042"
              ↑
         int("00000042") → i=42
         → email: loadtest0042@example.com
         → password: LoadTest!0042
```

- k6/load-test.js

SharedArray + papaparse で user.csv を読み込む
login() 関数が users[__VU % users.length] でVUごとにユーザーを割り当てる

### 作業ステップ

1. docker compose upして毎回クリーンな環境でテストする
2. 環境変数を設定
   ```
   KRATOS_ADMIN_URL=http://localhost:4434 \
   API_BASE_URL=http://localhost:8080 \
   USER_COUNT=100 \
   bash k6/seed-users.sh
   ```
3. seed-user.shを実行

   ```
   ./seed-users.sh
   ```

4. 負荷テスト実行
   ```
   k6 run \
    -e API_BASE_URL=http://localhost:8080 \
    -e KRATOS_PUBLIC_URL=http://localhost:4433 \
    k6/load-test.js
   ```

### 終了時

```
docker compose down -v
```

### 補足
