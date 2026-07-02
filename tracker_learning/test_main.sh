
#!/bin/bash

# Configuration
SERVER_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0;0m' # No Color

echo -e "${BLUE}=== Starting Groupie Tracker Verification Script ===${NC}\n"

# Exercise 1: The Data Fetcher
echo -e "${BLUE}[Exercise 1: Data Fetcher]${NC}"
RESP_ARTISTS=$(curl -s "https://groupietrackers.herokuapp.com/api/artists")
if [[ "$RESP_ARTISTS" == *"\"id\""* && "$RESP_ARTISTS" == *"\"name\""* ]]; then
    echo -e "${GREEN}✔ PASS: API reachable and returns artist data${NC}"
else
    echo -e "${RED}✘ FAIL: Could not reach API or unexpected response${NC}"
fi

RESP_LOCATIONS=$(curl -s "https://groupietrackers.herokuapp.com/api/locations")
if [[ "$RESP_LOCATIONS" == *"\"index\""* ]]; then
    echo -e "${GREEN}✔ PASS: Locations endpoint reachable${NC}"
else
    echo -e "${RED}✘ FAIL: Locations endpoint failed. Got '$RESP_LOCATIONS'${NC}"
fi

RESP_DATES=$(curl -s "https://groupietrackers.herokuapp.com/api/dates")
if [[ "$RESP_DATES" == *"\"index\""* ]]; then
    echo -e "${GREEN}✔ PASS: Dates endpoint reachable${NC}"
else
    echo -e "${RED}✘ FAIL: Dates endpoint failed. Got '$RESP_DATES'${NC}"
fi

RESP_RELATION=$(curl -s "https://groupietrackers.herokuapp.com/api/relation")
if [[ "$RESP_RELATION" == *"\"index\""* ]]; then
    echo -e "${GREEN}✔ PASS: Relation endpoint reachable${NC}"
else
    echo -e "${RED}✘ FAIL: Relation endpoint failed. Got '$RESP_RELATION'${NC}"
fi
echo ""

# Exercise 2: The Struct Decoder
echo -e "${BLUE}[Exercise 2: Struct Decoder]${NC}"
RESP_DECODE=$(curl -s "https://groupietrackers.herokuapp.com/api/artists")

if [[ "$RESP_DECODE" == *"\"members\""* ]]; then
    echo -e "${GREEN}✔ PASS: Members field present in API response${NC}"
else
    echo -e "${RED}✘ FAIL: Members field missing from API response${NC}"
fi

if [[ "$RESP_DECODE" == *"\"creationDate\""* ]]; then
    echo -e "${GREEN}✔ PASS: creationDate field present in API response${NC}"
else
    echo -e "${RED}✘ FAIL: creationDate field missing — check your JSON tags${NC}"
fi

if [[ "$RESP_DECODE" == *"\"firstAlbum\""* ]]; then
    echo -e "${GREEN}✔ PASS: firstAlbum field present in API response${NC}"
else
    echo -e "${RED}✘ FAIL: firstAlbum field missing — check your JSON tags${NC}"
fi
echo ""

# Exercise 3: The Home Page Handler
echo -e "${BLUE}[Exercise 3: Home Page Handler /]${NC}"
RESP_HOME=$(curl -s "$SERVER_URL/")
if [[ "$RESP_HOME" == *"<html"* ]]; then
    echo -e "${GREEN}✔ PASS: Home page returns HTML${NC}"
else
    echo -e "${RED}✘ FAIL: Home page did not return HTML. Got '$RESP_HOME'${NC}"
fi

if [[ "$RESP_HOME" == *"artist"* || "$RESP_HOME" == *"Artist"* ]]; then
    echo -e "${GREEN}✔ PASS: Home page contains artist data${NC}"
else
    echo -e "${RED}✘ FAIL: Home page does not contain artist data${NC}"
fi

STATUS_WRONG=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/wrongpath")
if [ "$STATUS_WRONG" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Wrong path returns 404 Not Found${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 for wrong path, got $STATUS_WRONG${NC}"
fi

STATUS_HOME=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/")
if [ "$STATUS_HOME" == "200" ]; then
    echo -e "${GREEN}✔ PASS: Home page returns 200 OK${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 200 for home page, got $STATUS_HOME${NC}"
fi
echo ""

# Exercise 4: The Artist Detail Page
echo -e "${BLUE}[Exercise 4: Artist Detail Page /artist]${NC}"
RESP_ARTIST=$(curl -s "$SERVER_URL/artist?id=1")
if [[ "$RESP_ARTIST" == *"<html"* ]]; then
    echo -e "${GREEN}✔ PASS: Artist page returns HTML${NC}"
else
    echo -e "${RED}✘ FAIL: Artist page did not return HTML. Got '$RESP_ARTIST'${NC}"
fi

if [[ "$RESP_ARTIST" == *"Queen"* || "$RESP_ARTIST" == *"members"* || "$RESP_ARTIST" == *"Members"* ]]; then
    echo -e "${GREEN}✔ PASS: Artist page contains artist details${NC}"
else
    echo -e "${RED}✘ FAIL: Artist page does not contain expected artist data${NC}"
fi

STATUS_INVALID_ID=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/artist?id=abc")
if [ "$STATUS_INVALID_ID" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Invalid ID returns 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 for invalid ID, got $STATUS_INVALID_ID${NC}"
fi

STATUS_MISSING_ID=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/artist")
if [ "$STATUS_MISSING_ID" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Missing ID returns 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 for missing ID, got $STATUS_MISSING_ID${NC}"
fi

STATUS_NOT_FOUND=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/artist?id=999")
if [ "$STATUS_NOT_FOUND" == "404" ]; then
    echo -e "${GREEN}✔ PASS: Unknown artist ID returns 404 Not Found${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 404 for unknown artist, got $STATUS_NOT_FOUND${NC}"
fi
echo ""

# Exercise 5: The Search Handler
echo -e "${BLUE}[Exercise 5: Search Handler /search]${NC}"
RESP_SEARCH=$(curl -s "$SERVER_URL/search?q=queen")
if [[ "$RESP_SEARCH" == *"Queen"* ]]; then
    echo -e "${GREEN}✔ PASS: Search for 'queen' returns Queen${NC}"
else
    echo -e "${RED}✘ FAIL: Search for 'queen' did not return Queen. Got '$RESP_SEARCH'${NC}"
fi

RESP_SEARCH_CASE=$(curl -s "$SERVER_URL/search?q=QUEEN")
if [[ "$RESP_SEARCH_CASE" == *"Queen"* ]]; then
    echo -e "${GREEN}✔ PASS: Search is case insensitive${NC}"
else
    echo -e "${RED}✘ FAIL: Search is case sensitive — fix strings.ToLower()${NC}"
fi

RESP_NO_RESULTS=$(curl -s "$SERVER_URL/search?q=zzzzzzzzz")
if [[ "$RESP_NO_RESULTS" == *"No artists found"* || "$RESP_NO_RESULTS" == *"no artists"* ]]; then
    echo -e "${GREEN}✔ PASS: Empty results shows no artists message${NC}"
else
    echo -e "${RED}✘ FAIL: No message shown for empty results. Got '$RESP_NO_RESULTS'${NC}"
fi

STATUS_EMPTY_QUERY=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/search?q=")
if [ "$STATUS_EMPTY_QUERY" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Empty search query returns 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 for empty query, got $STATUS_EMPTY_QUERY${NC}"
fi

STATUS_POST_SEARCH=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SERVER_URL/search")
if [ "$STATUS_POST_SEARCH" == "405" ]; then
    echo -e "${GREEN}✔ PASS: POST to /search returns 405 Method Not Allowed${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 405 for POST, got $STATUS_POST_SEARCH${NC}"
fi
echo ""

# Exercise 6: The Filter Handler
echo -e "${BLUE}[Exercise 6: Filter Handler /filter]${NC}"
RESP_FILTER_DATE=$(curl -s "$SERVER_URL/filter?creationDate=1970")
if [[ "$RESP_FILTER_DATE" == *"<html"* ]]; then
    echo -e "${GREEN}✔ PASS: Filter by creationDate returns HTML${NC}"
else
    echo -e "${RED}✘ FAIL: Filter by creationDate failed. Got '$RESP_FILTER_DATE'${NC}"
fi

RESP_FILTER_MEMBERS=$(curl -s "$SERVER_URL/filter?memberCount=4")
if [[ "$RESP_FILTER_MEMBERS" == *"<html"* ]]; then
    echo -e "${GREEN}✔ PASS: Filter by memberCount returns HTML${NC}"
else
    echo -e "${RED}✘ FAIL: Filter by memberCount failed. Got '$RESP_FILTER_MEMBERS'${NC}"
fi

RESP_FILTER_BOTH=$(curl -s "$SERVER_URL/filter?creationDate=1970&memberCount=4")
if [[ "$RESP_FILTER_BOTH" == *"<html"* ]]; then
    echo -e "${GREEN}✔ PASS: Filter by both parameters returns HTML${NC}"
else
    echo -e "${RED}✘ FAIL: Combined filter failed. Got '$RESP_FILTER_BOTH'${NC}"
fi

STATUS_NO_FILTER=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/filter")
if [ "$STATUS_NO_FILTER" == "400" ]; then
    echo -e "${GREEN}✔ PASS: No filter parameters returns 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 for no filters, got $STATUS_NO_FILTER${NC}"
fi

STATUS_INVALID_DATE=$(curl -s -o /dev/null -w "%{http_code}" "$SERVER_URL/filter?creationDate=abc")
if [ "$STATUS_INVALID_DATE" == "400" ]; then
    echo -e "${GREEN}✔ PASS: Invalid creationDate returns 400 Bad Request${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 400 for invalid date, got $STATUS_INVALID_DATE${NC}"
fi

STATUS_POST_FILTER=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$SERVER_URL/filter")
if [ "$STATUS_POST_FILTER" == "405" ]; then
    echo -e "${GREEN}✔ PASS: POST to /filter returns 405 Method Not Allowed${NC}"
else
    echo -e "${RED}✘ FAIL: Expected 405 for POST, got $STATUS_POST_FILTER${NC}"
fi

RESP_NO_MATCH=$(curl -s "$SERVER_URL/filter?creationDate=1800")
if [[ "$RESP_NO_MATCH" == *"No artists"* || "$RESP_NO_MATCH" == *"no artists"* ]]; then
    echo -e "${GREEN}✔ PASS: No matching filter shows empty results message${NC}"
else
    echo -e "${RED}✘ FAIL: No message for empty filter results. Got '$RESP_NO_MATCH'${NC}"
fi
echo ""

# Exercise 7: Unit Tests
echo -e "${BLUE}[Exercise 7: Unit Tests]${NC}"
if [ -f "handlers_test.go" ]; then
    echo -e "${GREEN}✔ PASS: handlers_test.go file exists${NC}"
else
    echo -e "${RED}✘ FAIL: handlers_test.go not found${NC}"
fi

if go test ./... > /dev/null 2>&1; then
    echo -e "${GREEN}✔ PASS: All unit tests pass${NC}"
else
    echo -e "${RED}✘ FAIL: Some unit tests failed — run 'go test ./...' to see details${NC}"
fi

TEST_COUNT=$(go test -v ./... 2>/dev/null | grep -c "^--- PASS")
if [ "$TEST_COUNT" -ge 5 ]; then
    echo -e "${GREEN}✔ PASS: At least 5 unit tests found and passing ($TEST_COUNT tests)${NC}"
else
    echo -e "${RED}✘ FAIL: Expected at least 5 tests, found $TEST_COUNT passing${NC}"
fi
echo ""

echo -e "${BLUE}=== Groupie Tracker Verification Complete ===${NC}"
