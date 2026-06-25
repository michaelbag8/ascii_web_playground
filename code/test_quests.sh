#!/bin/bash

# Configuration
SERVER_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0;0m' # No Color

echo -e "${BLUE}=== Starting Go HTTP Exercise Verification Script ===${NC}\n"

# -----------------------------
# Exercise 1: Method Inspector
# -----------------------------
echo -e "${BLUE}[Exercise 1: /method-inspector]${NC}"

RESP_GET=$(curl -s "$SERVER_URL/method-inspector")
if [[ "$RESP_GET" == *"GET request"* ]]; then
    echo -e "${GREEN}✔ PASS: GET method detected${NC}"
else
    echo -e "${RED}✘ FAIL: Unexpected GET response '$RESP_GET'${NC}"
fi

RESP_POST=$(curl -s -X POST "$SERVER_URL/method-inspector")
if [[ "$RESP_POST" == *"POST request"* ]]; then
    echo -e "${GREEN}✔ PASS: POST method detected${NC}"
else
    echo -e "${RED}✘ FAIL: Unexpected POST response '$RESP_POST'${NC}"
fi

RESP_PUT=$(curl -s -X PUT "$SERVER_URL/method-inspector")
if [[ "$RESP_PUT" == *"PUT request"* ]]; then
    echo -e "${GREEN}✔ PASS: PUT method detected${NC}"
else
    echo -e "${RED}✘ FAIL: Unexpected PUT response '$RESP_PUT'${NC}"
fi
echo ""

# -----------------------------
# Exercise 9: Body Truth Validator
# -----------------------------
echo -e "${BLUE}[Exercise 9: /validate-body]${NC}"

# Test valid POST body
RESP_VALID_BODY=$(curl -s -X POST -d "Hello Go" "$SERVER_URL/validate-body")
if [[ "$RESP_VALID_BODY" == *"Valid body received: Hello Go"* ]]; then
    echo -e "${GREEN}✔ PASS: Valid body accepted and returned${NC}"
else
    echo -e "${RED}✘ FAIL: Expected valid body confirmation, got '$RESP_VALID_BODY'${NC}"
fi

# Test empty body
RESP_EMPTY_BODY=$(curl -s -X POST -d "" "$SERVER_URL/validate-body")
STATUS_EMPTY_BODY=$(curl -s -o /dev/null -w "%{http_code}" -X POST -d "" "$SERVER_URL/validate-body")

if [ "$STATUS_EMPTY_BODY" == "400" ] && [[ "$RESP_EMPTY_BODY" == *"body is required"* ]]; then
    echo -e "${GREEN}✔ PASS: Empty body rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Empty body should return 400 and 'body is required'. Got status $STATUS_EMPTY_BODY, body '$RESP_EMPTY_BODY'${NC}"
fi

# Test whitespace-only body
RESP_BLANK_BODY=$(curl -s -X POST -d "   " "$SERVER_URL/validate-body")
STATUS_BLANK_BODY=$(curl -s -o /dev/null -w "%{http_code}" -X POST -d "   " "$SERVER_URL/validate-body")

if [ "$STATUS_BLANK_BODY" == "400" ] && [[ "$RESP_BLANK_BODY" == *"body cannot be blank"* ]]; then
    echo -e "${GREEN}✔ PASS: Whitespace-only body rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Whitespace-only body should return 400 and 'body cannot be blank'. Got status $STATUS_BLANK_BODY, body '$RESP_BLANK_BODY'${NC}"
fi

# Test wrong method
STATUS_BODY_GET=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/validate-body")

if [ "$STATUS_BODY_GET" == "405" ]; then
    echo -e "${GREEN}✔ PASS: GET request rejected with 405${NC}"
else
    echo -e "${RED}✘ FAIL: Expected GET to return 405, got $STATUS_BODY_GET${NC}"
fi

echo ""


# -----------------------------
# Exercise 10: Query Parameter Firewall
# -----------------------------
echo -e "${BLUE}[Exercise 10: /query-firewall]${NC}"

# Test valid user and role
RESP_QUERY_VALID=$(curl -s "$SERVER_URL/query-firewall?user=John&role=admin")
if [[ "$RESP_QUERY_VALID" == *"User: John, Role: admin"* ]]; then
    echo -e "${GREEN}✔ PASS: Valid user and role accepted${NC}"
else
    echo -e "${RED}✘ FAIL: Expected valid user/role response, got '$RESP_QUERY_VALID'${NC}"
fi

# Test missing user
RESP_QUERY_NO_USER=$(curl -s "$SERVER_URL/query-firewall?role=admin")
STATUS_QUERY_NO_USER=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/query-firewall?role=admin")

if [ "$STATUS_QUERY_NO_USER" == "400" ] && [[ "$RESP_QUERY_NO_USER" == *"user parameter missing"* ]]; then
    echo -e "${GREEN}✔ PASS: Missing user rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 and 'user parameter missing', got status $STATUS_QUERY_NO_USER, body '$RESP_QUERY_NO_USER'${NC}"
fi

# Test empty user value
RESP_QUERY_EMPTY_USER=$(curl -s "$SERVER_URL/query-firewall?user=&role=admin")
STATUS_QUERY_EMPTY_USER=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/query-firewall?user=&role=admin")

if [ "$STATUS_QUERY_EMPTY_USER" == "400" ] && [[ "$RESP_QUERY_EMPTY_USER" == *"user parameter missing"* ]]; then
    echo -e "${GREEN}✔ PASS: Empty user value rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Empty user should return 400, got status $STATUS_QUERY_EMPTY_USER, body '$RESP_QUERY_EMPTY_USER'${NC}"
fi

# Test missing role
RESP_QUERY_NO_ROLE=$(curl -s "$SERVER_URL/query-firewall?user=John")
STATUS_QUERY_NO_ROLE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/query-firewall?user=John")

if [ "$STATUS_QUERY_NO_ROLE" == "400" ] && [[ "$RESP_QUERY_NO_ROLE" == *"role parameter missing"* ]]; then
    echo -e "${GREEN}✔ PASS: Missing role rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 and 'role parameter missing', got status $STATUS_QUERY_NO_ROLE, body '$RESP_QUERY_NO_ROLE'${NC}"
fi

# Test empty role value
RESP_QUERY_EMPTY_ROLE=$(curl -s "$SERVER_URL/query-firewall?user=John&role=")
STATUS_QUERY_EMPTY_ROLE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/query-firewall?user=John&role=")

if [ "$STATUS_QUERY_EMPTY_ROLE" == "400" ] && [[ "$RESP_QUERY_EMPTY_ROLE" == *"role parameter missing"* ]]; then
    echo -e "${GREEN}✔ PASS: Empty role value rejected with 400${NC}"
else
    echo -e "${RED}✘ FAIL: Empty role should return 400, got status $STATUS_QUERY_EMPTY_ROLE, body '$RESP_QUERY_EMPTY_ROLE'${NC}"
fi

# Stretch test: invalid role should be forbidden
STATUS_QUERY_BAD_ROLE=$(curl -s -o /dev/null -w "%{http_code}" \
    "$SERVER_URL/query-firewall?user=John&role=superadmin")

if [ "$STATUS_QUERY_BAD_ROLE" == "403" ]; then
    echo -e "${GREEN}✔ PASS: Unsupported role rejected with 403 Forbidden${NC}"
else
    echo -e "${RED}✘ FAIL: Expected invalid role to return 403, got $STATUS_QUERY_BAD_ROLE${NC}"
fi

# Edge case: duplicate user parameters.
# r.URL.Query().Get("user") normally uses the first value: John.
RESP_QUERY_DUPLICATE=$(curl -s "$SERVER_URL/query-firewall?user=John&user=Mike&role=user")
if [[ "$RESP_QUERY_DUPLICATE" == *"User: John, Role: user"* ]]; then
    echo -e "${GREEN}✔ PASS: Duplicate user parameters use the first value${NC}"
else
    echo -e "${RED}✘ FAIL: Duplicate parameter behavior was unexpected: '$RESP_QUERY_DUPLICATE'${NC}"
fi

echo ""
# -----------------------------
# Exercise 4: Form Decoder
# -----------------------------
echo -e "${BLUE}[Exercise 4: /form]${NC}"

RESP_FORM=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d '{"username":"Alice","language":"Go"}' \
    "$SERVER_URL/form")

if [[ "$RESP_FORM" == *"Alice"* && "$RESP_FORM" == *"Go"* ]]; then
    echo -e "${GREEN}✔ PASS: Form parsed successfully${NC}"
else
    echo -e "${RED}✘ FAIL: Form parsing failed '$RESP_FORM'${NC}"
fi

STATUS_FORM_CT=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
    -H "Content-Type: text/plain" -d "username=Alice" "$SERVER_URL/form")

if [ "$STATUS_FORM_CT" == "415" ]; then
    echo -e "${GREEN}✔ PASS: Unsupported Content-Type rejected${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 415, got $STATUS_FORM_CT${NC}"
fi
echo ""

# -----------------------------
# Exercise 5: Status Code Factory
# -----------------------------
echo -e "${BLUE}[Exercise 5: /status]${NC}"

RESP_STATUS=$(curl -s "$SERVER_URL/status?code=404")
if [[ "$RESP_STATUS" == *"404"* && "$RESP_STATUS" == *"Not Found"* ]]; then
    echo -e "${GREEN}✔ PASS: Status 404 handled correctly${NC}"
else
    echo -e "${RED}✘ FAIL: Unexpected response '$RESP_STATUS'${NC}"
fi

STATUS_INVALID=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/status?code=abc")
if [ "$STATUS_INVALID" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Invalid code rejected${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400, got $STATUS_INVALID${NC}"
fi
echo ""

# -----------------------------
# Exercise 13: Multi-Stage Router System
# -----------------------------
echo -e "${BLUE}[Exercise 13: /api/v1/auth and /api/v1/data]${NC}"

# Test login route
RESP_LOGIN=$(curl -s "$SERVER_URL/api/v1/auth/login")
if [[ "$RESP_LOGIN" == *"login endpoint reached"* ]]; then
    echo -e "${GREEN}✔ PASS: Login route reached through all router layers${NC}"
else
    echo -e "${RED}✘ FAIL: Expected login response, got '$RESP_LOGIN'${NC}"
fi

# Test logout route
RESP_LOGOUT=$(curl -s "$SERVER_URL/api/v1/auth/logout")
if [[ "$RESP_LOGOUT" == *"logout endpoint reached"* ]]; then
    echo -e "${GREEN}✔ PASS: Logout route reached through all router layers${NC}"
else
    echo -e "${RED}✘ FAIL: Expected logout response, got '$RESP_LOGOUT'${NC}"
fi

# Test data info route
RESP_INFO=$(curl -s "$SERVER_URL/api/v1/data/info")
if [[ "$RESP_INFO" == *"data info endpoint reached"* ]]; then
    echo -e "${GREEN}✔ PASS: Data info route reached through all router layers${NC}"
else
    echo -e "${RED}✘ FAIL: Expected data info response, got '$RESP_INFO'${NC}"
fi

# Test unknown auth route
STATUS_UNKNOWN_AUTH=$(curl -s -o /dev/null -w "%{http_code}" \
    "$SERVER_URL/api/v1/auth/register")

if [ "$STATUS_UNKNOWN_AUTH" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Unknown auth route returned 404${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 for unknown auth route, got $STATUS_UNKNOWN_AUTH${NC}"
fi

# Test unknown data route
STATUS_UNKNOWN_DATA=$(curl -s -o /dev/null -w "%{http_code}" \
    "$SERVER_URL/api/v1/data/delete")

if [ "$STATUS_UNKNOWN_DATA" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Unknown data route returned 404${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 for unknown data route, got $STATUS_UNKNOWN_DATA${NC}"
fi

# Test missing /api prefix
STATUS_NO_API=$(curl -s -o /dev/null -w "%{http_code}" \
    "$SERVER_URL/v1/auth/login")

if [ "$STATUS_NO_API" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Missing /api prefix returned 404${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 without /api prefix, got $STATUS_NO_API${NC}"
fi

# Test wrong API version
STATUS_V2=$(curl -s -o /dev/null -w "%{http_code}" \
    "$SERVER_URL/api/v2/auth/login")

if [ "$STATUS_V2" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Unsupported API version returned 404${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 for /api/v2 route, got $STATUS_V2${NC}"
fi

echo ""

# -----------------------------
# Final Summary
# -----------------------------
echo -e "${BLUE}=== Testing Complete ===${NC}"