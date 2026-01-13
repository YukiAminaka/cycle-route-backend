#!/bin/bash

# Kratos Admin APIを使用してテストユーザーを作成し、セッションクッキーを取得するスクリプト

set -e

# 色付きの出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# デフォルト値
KRATOS_ADMIN_URL="${KRATOS_ADMIN_URL:-http://127.0.0.1:4434}"
KRATOS_PUBLIC_URL="${KRATOS_PUBLIC_URL:-http://127.0.0.1:4433}"
EMAIL="${1:-testuser$(date +%s)@example.com}"
PASSWORD="${2:-SecureP@ssw0rd!2024}"

echo -e "${YELLOW}=== Kratos Test User Creator (Admin API) ===${NC}"
echo "Email: $EMAIL"
echo "Password: $PASSWORD"
echo "Kratos Admin URL: $KRATOS_ADMIN_URL"
echo ""

# jqがインストールされているか確認
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed. Please install jq first.${NC}"
    echo "  Ubuntu/Debian: sudo apt-get install jq"
    echo "  macOS: brew install jq"
    exit 1
fi

# 一時的なクッキーファイルを作成（スクリプト終了時に自動削除）
COOKIE_FILE=$(mktemp)
trap "rm -f $COOKIE_FILE" EXIT

# Step 1: Admin APIでIdentityを作成
echo -e "${YELLOW}[1/3] Creating identity via Admin API...${NC}"

USERNAME=$(echo $EMAIL | cut -d'@' -f1)

# パスワードハッシュを作成（Admin APIではプレーンテキストでパスワードを設定可能）
IDENTITY_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$KRATOS_ADMIN_URL/admin/identities" \
    -H 'Content-Type: application/json' \
    -d "{
        \"schema_id\": \"default\",
        \"traits\": {
            \"email\": \"$EMAIL\",
            \"username\": \"$USERNAME\",
            \"name\": {
                \"first\": \"Test\",
                \"last\": \"User\"
            }
        },
        \"credentials\": {
            \"password\": {
                \"config\": {
                    \"password\": \"$PASSWORD\"
                }
            }
        },
        \"state\": \"active\"
    }")

HTTP_CODE=$(echo "$IDENTITY_RESPONSE" | tail -n1)
RESPONSE_BODY=$(echo "$IDENTITY_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" != "201" ]; then
    # ユーザーが既に存在する可能性がある
    ERROR_MSG=$(echo "$RESPONSE_BODY" | jq -r '.error.message // "Unknown error"')

    if echo "$RESPONSE_BODY" | grep -q "unique"; then
        echo -e "${YELLOW}⚠ User already exists. Trying to login...${NC}"
    else
        echo -e "${RED}✗ Failed to create identity: $ERROR_MSG${NC}"
        echo "Response: $RESPONSE_BODY" | jq '.'
        exit 1
    fi
else
    IDENTITY_ID=$(echo "$RESPONSE_BODY" | jq -r '.id')
    echo -e "${GREEN}✓ Identity created successfully (ID: $IDENTITY_ID)${NC}"
fi

# Step 2: ログインしてセッションを取得
echo ""
echo -e "${YELLOW}[2/3] Logging in to get session cookie...${NC}"

LOGIN_FLOW=$(curl -s "$KRATOS_PUBLIC_URL/self-service/login/api" \
    -H 'Accept: application/json' \
    -c $COOKIE_FILE)

LOGIN_FLOW_ID=$(echo $LOGIN_FLOW | jq -r '.id')
LOGIN_CSRF=$(echo $LOGIN_FLOW | jq -r '.ui.nodes[] | select(.attributes.name=="csrf_token") | .attributes.value')

if [ "$LOGIN_FLOW_ID" = "null" ] || [ -z "$LOGIN_FLOW_ID" ]; then
    echo -e "${RED}Error: Failed to initialize login flow${NC}"
    exit 1
fi

sleep 1  # Kratosの処理を待つ

LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$KRATOS_PUBLIC_URL/self-service/login?flow=$LOGIN_FLOW_ID" \
    -H 'Content-Type: application/json' \
    -H 'Accept: application/json' \
    -b $COOKIE_FILE \
    -c $COOKIE_FILE \
    -d "{
        \"method\": \"password\",
        \"csrf_token\": \"$LOGIN_CSRF\",
        \"identifier\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

LOGIN_HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
LOGIN_RESPONSE_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

# ログイン成功の確認
if [ "$LOGIN_HTTP_CODE" = "200" ] || [ "$LOGIN_HTTP_CODE" = "422" ]; then
    # APIフローの場合、session_tokenが返される
    SESSION_TOKEN=$(echo "$LOGIN_RESPONSE_BODY" | jq -r '.session_token // empty')
    IDENTITY_ID=$(echo "$LOGIN_RESPONSE_BODY" | jq -r '.session.identity.id // empty')

    if [ -n "$SESSION_TOKEN" ]; then
        echo -e "${GREEN}✓ Login successful!${NC}"

        # Step 3: セッション情報を表示
        echo ""
        echo -e "${YELLOW}[3/3] Session information obtained${NC}"

        echo ""
        echo -e "${GREEN}=== Test User Details ===${NC}"
        echo "Email: $EMAIL"
        echo "Password: $PASSWORD"
        echo "Identity ID: $IDENTITY_ID"
        echo ""
        echo -e "${GREEN}Session Token:${NC}"
        echo "$SESSION_TOKEN"
        echo ""
        echo -e "${YELLOW}Usage examples:${NC}"
        echo ""
        echo "# 1. Using session token with curl:"
        echo "  curl -H 'Cookie: ory_kratos_session=$SESSION_TOKEN' \\"
        echo "    http://localhost:8080/api/v1/users/$IDENTITY_ID"
        echo ""
        echo "# 2. Create a route:"
        echo "  curl -H 'Cookie: ory_kratos_session=$SESSION_TOKEN' \\"
        echo "    http://localhost:8080/api/v1/routes \\"
        echo "    -X POST -H 'Content-Type: application/json' \\"
        echo "    -d '{\"name\":\"Test Route\",...}'"
        echo ""
        echo "# 3. For Swagger UI:"
        echo "  - Open: http://localhost:8080/api/v1/swagger/index.html"
        echo "  - Click: 'Authorize' button (padlock icon)"
        echo "  - In the CookieAuth field, paste the session token above"
        echo "  - Click 'Authorize' and then 'Close'"
        echo ""
        echo -e "${GREEN}✓ You can now use this session token in Swagger UI${NC}"
        exit 0
    fi

    # クッキーでもチェック（念のため）
    if grep -q "ory_kratos_session" $COOKIE_FILE 2>/dev/null; then
        echo -e "${GREEN}✓ Login successful!${NC}"
        COOKIE_VALUE=$(grep ory_kratos_session $COOKIE_FILE | awk '{print $7}')

        echo ""
        echo -e "${GREEN}=== Test User Details ===${NC}"
        echo "Email: $EMAIL"
        echo "Password: $PASSWORD"
        echo ""
        echo -e "${GREEN}Session Token:${NC}"
        echo "$COOKIE_VALUE"
        echo ""
        echo -e "${YELLOW}Usage examples:${NC}"
        echo ""
        echo "# 1. Using session token with curl:"
        echo "  curl -H 'Cookie: ory_kratos_session=$COOKIE_VALUE' \\"
        echo "    http://localhost:8080/api/v1/users/USER_ID"
        echo ""
        echo "# 2. Create a route:"
        echo "  curl -H 'Cookie: ory_kratos_session=$COOKIE_VALUE' \\"
        echo "    http://localhost:8080/api/v1/routes \\"
        echo "    -X POST -H 'Content-Type: application/json' \\"
        echo "    -d '{\"name\":\"Test Route\",...}'"
        echo ""
        echo "# 3. For Swagger UI:"
        echo "  - Open: http://localhost:8080/api/v1/swagger/index.html"
        echo "  - Click: 'Authorize' button"
        echo "  - In the CookieAuth field, paste the session token above"
        echo ""
        echo -e "${GREEN}✓ You can now use this session token in Swagger UI${NC}"
        exit 0
    fi
fi

# ログイン失敗
LOGIN_ERROR=$(echo "$LOGIN_RESPONSE_BODY" | jq -r '.ui.messages[]? | select(.type=="error") | .text' 2>/dev/null | head -1)
if [ -n "$LOGIN_ERROR" ]; then
    echo -e "${RED}✗ Login failed: $LOGIN_ERROR${NC}"
else
    echo -e "${RED}✗ Login failed (HTTP $LOGIN_HTTP_CODE)${NC}"
fi

echo ""
echo "Login response:"
echo "$LOGIN_RESPONSE_BODY" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE_BODY"
echo ""
echo -e "${YELLOW}Tip: Check if Kratos is running properly${NC}"
echo "  - Admin API: $KRATOS_ADMIN_URL/health/ready"
echo "  - Public API: $KRATOS_PUBLIC_URL/health/ready"
exit 1
