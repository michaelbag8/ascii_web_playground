#!/bin/bash
SERVER_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0;0m'

echo -e "${BLUE}=== HTTP Resurgence Verification ===${NC}\n"

# Exercise 1: /method-inspector
echo -e "${BLUE}[Exercise 1: /method-inspector]${NC}"
R1=$(curl -s -X GET "$SERVER_URL/method-inspector")
if [[ "$R1" == *"GET"* ]]; then echo -e "${GREEN}✔ PASS: GET detected${NC}"; else echo -e "${RED}✘ FAIL: got '$R1'${NC}"; fi
R1P=$(curl -s -X POST "$SERVER_URL/method-inspector")
if [[ "$R1P" == *"POST"* ]]; then echo -e "${GREEN}✔ PASS: POST detected${NC}"; else echo -e "${RED}✘ FAIL: got '$R1P'${NC}"; fi


# Exercise 2: /echo
echo -e "\n${BLUE}[Exercise 2: /echo]${NC}"
R2=$(curl -s -X POST -d "Hello Go" "$SERVER_URL/echo")
if [[ "$R2" == *"Hello Go"* ]]; then echo -e "${GREEN}✔ PASS: body echoed${NC}"; else echo -e "${RED}✘ FAIL: got '$R2'${NC}"; fi
R2G=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$SERVER_URL/echo")
if [ "$R2G" == "405" ]; then echo -e "${GREEN}✔ PASS: GET blocked with 405${NC}"; else echo -e "${RED}✘ FAIL: expected 405 got $R2G${NC}"; fi

# Exercise 3: /headers
echo -e "\n${BLUE}[Exercise 3: /headers]${NC}"
R3=$(curl -s -H "X-Custom-Token: abc123" "$SERVER_URL/headers")
if [[ "$R3" == *"abc123"* ]]; then echo -e "${GREEN}✔ PASS: header echoed${NC}"; else echo -e "${RED}✘ FAIL: got '$R3'${NC}"; fi
R3E=$(curl -s "$SERVER_URL/headers")
if [[ "$R3E" == *"missing"* || "$R3E" == *"Missing"* ]]; then echo -e "${GREEN}✔ PASS: missing header handled${NC}"; else echo -e "${RED}✘ FAIL: got '$R3E'${NC}"; fi

# Exercise 4: /form
echo -e "\n${BLUE}[Exercise 4: /form]${NC}"
R4=$(curl -s -X POST -d "username=Ada&language=Go" "$SERVER_URL/form")
if [[ "$R4" == *"Ada"* && "$R4" == *"Go"* ]]; then echo -e "${GREEN}✔ PASS: form parsed${NC}"; else echo -e "${RED}✘ FAIL: got '$R4'${NC}"; fi
R4E=$(curl -s -o /dev/null -w "%{http_code}" -X POST -d "username=&language=Go" "$SERVER_URL/form")
if [ "$R4E" == "400" ]; then echo -e "${GREEN}✔ PASS: empty field returns 400${NC}"; else echo -e "${RED}✘ FAIL: expected 400 got $R4E${NC}"; fi

# Exercise 5: /status
echo -e "\n${BLUE}[Exercise 5: /status]${NC}"
R5=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/status?code=404")
if [ "$R5" == "404" ]; then echo -e "${GREEN}✔ PASS: 404 returned${NC}"; else echo -e "${RED}✘ FAIL: expected 404 got $R5${NC}"; fi
R5B=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/status?code=banana")
if [ "$R5B" == "404" ]; then echo -e "${GREEN}✔ PASS: bad code returns 404${NC}"; else echo -e "${RED}✘ FAIL: expected 404 got $R5B${NC}"; fi

# Exercise 6: /api/v1/greet (ServeMux subtree)
echo -e "\n${BLUE}[Exercise 6: /api/v1/greet]${NC}"
R6=$(curl -s "$SERVER_URL/api/v1/greet?name=Zion")
if [[ "$R6" == *"Zion"* ]]; then echo -e "${GREEN}✔ PASS: subtree route greet works${NC}"; else echo -e "${RED}✘ FAIL: got '$R6'${NC}"; fi
R6P=$(curl -s "$SERVER_URL/api/v1/ping")
if [[ "$R6P" == *"pong"* ]]; then echo -e "${GREEN}✔ PASS: subtree route ping works${NC}"; else echo -e "${RED}✘ FAIL: got '$R6P'${NC}"; fi

# Exercise 7: /render (template)
echo -e "\n${BLUE}[Exercise 7: /render]${NC}"
R7=$(curl -s "$SERVER_URL/render?title=SENTINEL&body=Online")
if [[ "$R7" == *"SENTINEL"* && "$R7" == *"Online"* ]]; then echo -e "${GREEN}✔ PASS: template rendered${NC}"; else echo -e "${RED}✘ FAIL: got '$R7'${NC}"; fi
R7E=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/render")
if [ "$R7E" == "400" ]; then echo -e "${GREEN}✔ PASS: missing params returns 400${NC}"; else echo -e "${RED}✘ FAIL: expected 400 got $R7E${NC}"; fi

echo -e "\n${BLUE}=== Verification Complete ===${NC}"