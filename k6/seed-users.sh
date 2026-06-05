#!/bin/bash
# Ory Kratos にk6負荷テスト用ユーザーを一括登録し、アプリDBにも作成する
set -euo pipefail

export KRATOS_ADMIN_URL="${KRATOS_ADMIN_URL:-http://localhost:4434}"
export API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
USER_COUNT="${USER_COUNT:-100}"
BATCH_SIZE=200  # プレーンテキストPWの上限
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
export CSV_FILE="$SCRIPT_DIR/user.csv"

echo "identifier,password" > "$CSV_FILE"

for batch_start in $(seq 1 $BATCH_SIZE $USER_COUNT); do
  batch_end=$(( batch_start + BATCH_SIZE - 1 ))
  [ $batch_end -gt $USER_COUNT ] && batch_end=$USER_COUNT

  identities_json=""
  for i in $(seq $batch_start $batch_end); do
    email="loadtest$(printf '%04d' $i)@example.com"
    password="LoadTest!$(printf '%04d' $i)"
    # patch_idの先頭8桁にiを埋め込み、Python側でindexを復元する
    patch_id="$(printf '%08d-0000-0000-0000-%012d' $i $i)"

    entry=$(printf '{
      "patch_id": "%s",
      "create": {
        "schema_id": "default",
        "state": "active",
        "traits": { "email": "%s" },
        "credentials": {
          "password": { "config": { "password": "%s" } }
        }
      }
    }' "$patch_id" "$email" "$password")

    identities_json="${identities_json:+$identities_json,}$entry"
  done

  # Step 1: Kratosにユーザーを一括登録
  echo "Registering users $batch_start-$batch_end in Kratos..."
  response=$(curl -s -w "\n%{http_code}" -X PATCH "$KRATOS_ADMIN_URL/admin/identities" \
    -H 'Content-Type: application/json' \
    -d "{\"identities\": [$identities_json]}")

  http_code=$(echo "$response" | tail -1)
  body=$(echo "$response" | head -n -1)

  if [ "$http_code" != "200" ]; then
    echo "Kratos error: HTTP $http_code"
    echo "$body"
    exit 1
  fi

  # Step 2: 取得したKratos IDでアプリDBにユーザーを作成しCSVに追記
  echo "$body" | python3 <<'PYEOF'
import sys, json, os, urllib.request, urllib.error

body = json.load(sys.stdin)
api_base = os.environ["API_BASE_URL"]
csv_file = os.environ["CSV_FILE"]

created = 0
for item in body.get("identities", []):
    if item.get("action") != "create":
        continue

    kratos_id = item["identity"]
    patch_id = item["patch_id"]

    # patch_idの先頭8桁からindexを復元
    i = int(patch_id.split("-")[0])
    email = f"loadtest{i:04d}@example.com"
    password = f"LoadTest!{i:04d}"
    name = f"LoadTestUser{i:04d}"

    payload = json.dumps({"kratos_id": kratos_id, "name": name, "email": email}).encode()
    req = urllib.request.Request(
        f"{api_base}/api/v1/users",
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req) as res:
            if res.status == 201:
                with open(csv_file, "a") as f:
                    f.write(f"{email},{password}\n")
                created += 1
            else:
                print(f"  App DB error for {email}: HTTP {res.status}", file=sys.stderr)
    except urllib.error.HTTPError as e:
        print(f"  App DB error for {email}: {e.code} {e.read().decode()}", file=sys.stderr)

print(f"  {created} users created in app DB")
PYEOF

done

echo ""
echo "Done. Seeded up to $USER_COUNT users."
echo "CSV written to $CSV_FILE"
