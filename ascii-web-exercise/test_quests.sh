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
# Exercise 2: Echo Chamber
# -----------------------------
echo -e "${BLUE}[Exercise 2: /echo]${NC}"

RESP_ECHO=$(curl -s -X POST -d "hello world" "$SERVER_URL/echo")
if [[ "$RESP_ECHO" == "hello world" ]]; then
    echo -e "${GREEN}✔ PASS: Echo returned exact body${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 'hello world', got '$RESP_ECHO'${NC}"
fi

STATUS_ECHO_GET=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/echo")
if [ "$STATUS_ECHO_GET" == "405" ]; then
    echo -e "${GREEN}✔ PASS: GET blocked with 405${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 405, got $STATUS_ECHO_GET${NC}"
fi
echo ""

# -----------------------------
# Exercise 3: Header Detective
# -----------------------------
echo -e "${BLUE}[Exercise 3: /headers]${NC}"

STATUS_HEADER_MISSING=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/headers")
if [ "$STATUS_HEADER_MISSING" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Missing X-Custom-Token rejected${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400, got $STATUS_HEADER_MISSING${NC}"
fi

RESP_HEADER=$(curl -s -H "X-Custom-Token: abc123" -H "Content-Type: application/json" "$SERVER_URL/headers")
if [[ "$RESP_HEADER" == *"abc123"* && "$RESP_HEADER" == *"Content-Type"* ]]; then
    echo -e "${GREEN}✔ PASS: Headers correctly parsed${NC}"
else
    echo -e "${RED}✘ FAIL: Header parsing failed: '$RESP_HEADER'${NC}"
fi
echo ""

# -----------------------------
# Exercise 4: Form Decoder
# -----------------------------
echo -e "${BLUE}[Exercise 4: /form]${NC}"

RESP_FORM=$(curl -s -X POST -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=Alice&language=Go" "$SERVER_URL/form")

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
# Exercise 6: API Subtree
# -----------------------------
echo -e "${BLUE}[Exercise 6: /api/v1]${NC}"

RESP_PING=$(curl -s "$SERVER_URL/api/v1/ping")
if [[ "$RESP_PING" == "pong" ]]; then
    echo -e "${GREEN}✔ PASS: /ping works${NC}"
else
    echo -e "${RED}✘ FAIL: /ping failed '$RESP_PING'${NC}"
fi

RESP_GREET=$(curl -s "$SERVER_URL/api/v1/greet?name=Zion")
if [[ "$RESP_GREET" == *"Zion"* ]]; then
    echo -e "${GREEN}✔ PASS: /greet works with name${NC}"
else
    echo -e "${RED}✘ FAIL: /greet failed '$RESP_GREET'${NC}"
fi

RESP_GREET_EMPTY=$(curl -s "$SERVER_URL/api/v1/greet")
if [[ "$RESP_GREET_EMPTY" == *"Stranger"* ]]; then
    echo -e "${GREEN}✔ PASS: /greet fallback works${NC}"
else
    echo -e "${RED}✘ FAIL: Missing fallback '$RESP_GREET_EMPTY'${NC}"
fi

echo ""

# -----------------------------
# Final Summary
# -----------------------------
echo -e "${BLUE}=== Testing Complete ===${NC}"