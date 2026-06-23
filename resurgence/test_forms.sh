#!/bin/bash

URL="http://localhost:8080/form"

BLUE='\033[0;34m'
NC='\033[0;0m'

echo "Running /form endpoint tests..."
echo "--------------------------------------"

pass=0
fail=0

run_test () {
  description=$1
  expected_code=$2
  response=$(curl -s -o /tmp/resp.txt -w "%{http_code}" "${@:3}")
  body=$(cat /tmp/resp.txt)

  echo ""
  echo "Test: $description"
  echo "HTTP Code: $response"

  if [ "$response" -eq "$expected_code" ]; then
    echo "✔ Status OK"
    pass=$((pass + 1))
  else
    echo "✘ Expected $expected_code"
    fail=$((fail + 1))
  fi

  echo "Response Body: $body"
}

# -----------------------------
# TEST CASES
# -----------------------------

# 1. Valid request
run_test "Valid request" 200 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Michael&language=Go" "$URL"

# 2. Wrong method
run_test "GET method not allowed" 405 \
  "$URL"

# 3. Missing username
run_test "Missing username" 400 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "language=Go" "$URL"

# 4. Missing language
run_test "Missing language" 400 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Michael" "$URL"

# 5. Both missing
run_test "Both fields missing" 400 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "" "$URL"

# 6. Wrong content type
run_test "Unsupported Content-Type" 415 \
  -X POST -H "Content-Type: text/plain" \
  -d "username=Michael&language=Go" "$URL"

# 7. Empty username
run_test "Empty username" 400 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=&language=Go" "$URL"

# 8. Empty language
run_test "Empty language" 400 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Michael&language=" "$URL"

# 9. Extra field ignored
run_test "Extra field ignored" 200 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=Michael&language=Go&country=Nigeria" "$URL"

# 10. URL encoded values
run_test "URL encoded values" 200 \
  -X POST -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "username=Michael Bulus" \
  --data-urlencode "language=Go Lang" \
  "$URL"

# -----------------------------
# SUMMARY
# -----------------------------

echo ""
echo "--------------------------------------"
echo "Tests passed: $pass"
echo "Tests failed: $fail"
echo "--------------------------------------"

if [ "$fail" -eq 0 ]; then
  echo "🎉 All tests passed!"
else
  echo "❌ Some tests failed"
fi

echo -e "\n${BLUE}=== Verification Complete ===${NC}"