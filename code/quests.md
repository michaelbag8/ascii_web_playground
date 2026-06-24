1.
---

# 🧩 Exercise: Admin API Subtree

## 🎯 Goal

Build a small **admin API system** under the `/admin/v1/` prefix using a separate `http.ServeMux`.

You will mount it onto a main mux using `http.StripPrefix`.

---

# 🧠 Concept you are practicing

* Sub-router (`ServeMux` inside another `ServeMux`)
* URL prefix mounting (`/admin/`)
* Request routing separation
* Query parameters
* Basic validation

---

# 📦 Requirements

## 1. Create two muxes

* `mainMux := http.NewServeMux()`
* `adminMux := http.NewServeMux()`

---

## 2. Admin routes (inside `adminMux` only)

### 🔹 `/v1/status`

* Method: `GET`
* Response:

  ```
  Admin system is running
  ```

---

### 🔹 `/v1/user`

* Method: `GET`
* Reads query parameter:

```
id
```

### Behavior:

| Condition      | Response           |
| -------------- | ------------------ |
| id is provided | `User ID: [id]`    |
| id is missing  | `User ID: unknown` |

---

## 3. Mount adminMux on mainMux

* All admin routes must be accessible under:

```text
/admin/
```

So:

```text
/admin/v1/status
/admin/v1/user?id=10
```

---

## 4. Server setup

Start server on:

```text
:8080
```

Using:

```go
http.ListenAndServe(":8080", mainMux)
```

---

# 🧭 Expected request flow

## Example 1

```text
GET /admin/v1/status
```

Flow:

```
mainMux → matches /admin/
Strip "/admin"
adminMux receives → /v1/status
```

Response:

```
Admin system is running
```

---

## Example 2

```text
GET /admin/v1/user?id=42
```

Response:

```
User ID: 42
```

---

## Example 3

```text
GET /admin/v1/user
```

Response:

```
User ID: unknown
```

---

# 🧪 What you are being tested on

You must understand:

### ✔ Sub-mux routing

adminMux handles only `/v1/...`

### ✔ Prefix mounting

mainMux handles `/admin/...`

### ✔ StripPrefix concept

removes `/admin` before passing request

### ✔ Query parameters

`r.URL.Query().Get("id")`

---

# 🔥 Challenge rules (important)

* Do NOT register `/admin/v1/...` in mainMux
* adminMux must NOT know about `/admin`
* mainMux must NOT know about `/v1`
* Each mux has a clear responsibility

---

# 💡 Bonus thinking question

Before coding, ask yourself:

> “At what stage does `/admin` disappear, and who is responsible for `/v1` routing?”

---



# Here
Perfect — I’ll now generate a **new set of similar exercises with deeper edge cases, trickier HTTP behavior, and more realistic backend patterns**, while keeping the same learning progression.

---

# 🧩 Exercise 8: Method Gatekeeper (Strict Version)

## 🎯 Goal

Build a `/secure-method` endpoint that **only allows GET and POST**, and properly rejects everything else.

---

## Key Tasks

* Register `/secure-method`
* Read `r.Method`
* Behavior:

| Method        | Response                             |
| ------------- | ------------------------------------ |
| GET           | `"GET method accepted"`              |
| POST          | `"POST method accepted"`             |
| Anything else | `"Method [METHOD] is not supported"` |

---

## Edge cases to think about

* What happens if method is `DELETE` or `PATCH`?
* What if method is lowercase? (spoiler: Go never sends lowercase)
* Can method be empty?

---

## Stretch

* Return `405 Method Not Allowed` for unsupported methods
* Add header:

  ```go
  Allow: GET, POST
  ```
* What happens if `Allow` is missing?

---

---

# 🧩 Exercise 9: Body Truth Validator

## 🎯 Goal

Build a `/validate-body` endpoint that reads request body and validates its exact structure.

---

## Key Tasks

* Accept only `POST`
* Read body using `io.ReadAll`
* Rules:

| Condition                 | Response                               |
| ------------------------- | -------------------------------------- |
| empty body                | `400 body is required`                 |
| body contains only spaces | `400 body cannot be blank`             |
| body is valid             | return `"Valid body received: [body]"` |

---

## Edge cases

* What is the difference between:

  * `""`
  * `"   "`
* What does `len(body)` return for each?
* What happens if body is read twice?

---

## Stretch

* Add header:

  ```go
  Content-Length: <length>
  ```
* What happens if Content-Length is missing?

---

---

# 🧩 Exercise 10: Query Parameter Firewall

## 🎯 Goal

Build `/query-firewall` that validates query parameters strictly.

---

## Key Tasks

Accept request like:

```
/query-firewall?user=John&role=admin
```

Rules:

* `user` is required
* `role` is required
* If missing either → 400 error with specific message

---

## Response rules

| Case         | Response                       |
| ------------ | ------------------------------ |
| missing user | `"user parameter missing"`     |
| missing role | `"role parameter missing"`     |
| both present | `"User: [user], Role: [role]"` |

---

## Edge cases

* What if:

  ```
  ?user=&role=admin
  ```
* What if duplicate keys:

  ```
  ?user=John&user=Mike
  ```

---

## Stretch

* Reject role values except:

  * `admin`
  * `user`
* Anything else → `403 Forbidden`

---

---

# 🧩 Exercise 11: JSON Trap Endpoint

## 🎯 Goal

Teach difference between `FormValue()` and JSON decoding.

---

## Key Tasks

* Accept only `POST`

* Content-Type must be:

  ```
  application/json
  ```

* Body example:

```json
{
  "username": "Ada",
  "language": "Go"
}
```

---

## Rules

| Case                 | Response                      |
| -------------------- | ----------------------------- |
| invalid content-type | 415                           |
| invalid JSON         | 400 "invalid JSON body"       |
| missing username     | 400                           |
| missing language     | 400                           |
| valid                | `"Hello Ada, Go is awesome!"` |

---

## Edge cases

* What if body is:

  ```
  {}
  ```
* What if:

  * extra fields exist
* What if JSON is malformed?

---

## Stretch

* Reject unknown fields strictly
* Add max body size limit:

  ```go
  http.MaxBytesReader(...)
  ```

---

---

# 🧩 Exercise 12: Header Integrity Checker

## 🎯 Goal

Validate multiple headers and enforce strict rules.

---

## Required headers

* `X-API-Key`
* `X-Request-ID`

---

## Rules

| Case               | Response             |
| ------------------ | -------------------- |
| missing API key    | 400                  |
| missing request ID | 400                  |
| both present       | `"Request accepted"` |

---

## Edge cases

* Header casing variations
* Empty header values:

  ```
  X-API-Key:
  ```
* Multiple values in header

---

## Stretch

* `X-Request-ID` must be a valid UUID format
* If invalid → 400 `"invalid request id"`

---

---

# 🧩 Exercise 13: Multi-Stage Router System

## 🎯 Goal

Build a 3-layer routing system:

```
mainMux
  └── /api/
        └── v1Mux
              └── /auth/
              └── /data/
```

---

## Required structure

### v1/auth

* `/login`
* `/logout`

### v1/data

* `/info`

---

## Rules

* `/api/` is mounted using StripPrefix
* `/auth` and `/data` are separate muxes inside v1Mux

---

## Edge cases

* What happens if:

  * StripPrefix is removed?
  * wrong trailing slash is used?
  * handler overlaps?

---

## Stretch

* Add logging middleware that prints:

  ```
  METHOD PATH
  ```

---

---

# 🧩 Exercise 14: Status Code Stress Tester

## 🎯 Goal

Build `/stress-status` endpoint that behaves like a controlled chaos generator.

---

## Rules

* Accept query:

  ```
  ?code=XYZ
  ```

---

## Behavior

| Case         | Response                |
| ------------ | ----------------------- |
| missing code | 400                     |
| not number   | 400                     |
| out of range | 400                     |
| valid        | return that status code |

---

## Edge cases

* `?code=200.5`
* `?code=-1`
* `?code=999`

---

## Stretch

* If code is 500–599 → add header:

  ```
  X-Server-Alert: true
  ```

---

# 🔥 If you want next level

I can now upgrade this into:

* 🧪 full automated test scripts (like real QA suites)
* 🧱 mini-project combining ALL exercises into one API
* 🧭 or a “mock interview backend test” (very close to real job tasks)

Just tell me 👍



# Mini Project: Unified HTTP API Server

## Goal

Combine all previous exercises into a single structured API server with multiple endpoints, middleware rules, and sub-routers.

---

## API Structure

### Root level

* `/method-inspector`
* `/echo`
* `/headers`
* `/form`
* `/status`
* `/render`

---

### Sub-router

Mounted at:

```
/api/v1/
```

Endpoints:

* `/ping`
* `/greet`

---

## Functional Requirements

### 1. Method Inspector

* Accept all methods
* Echo method back in response

---

### 2. Echo Service

* POST only
* returns raw request body

---

### 3. Header Inspector

* Validate required headers
* Return structured response

---

### 4. Form Handler

* Parse URL-encoded forms
* Validate required fields

---

### 5. Status Endpoint

* Dynamic status code responses
* Strict validation rules

---

### 6. Template Renderer

* Inline HTML template rendering
* Query-based input

---

### 7. API Subsystem

* Mounted under `/api/v1/`
* Separate ServeMux
* Uses StripPrefix routing

---

## Middleware Layer (optional but recommended)

* Method validation middleware
* Logging middleware
* Content-Type validation middleware

---

## Edge Cases to Handle

* Missing headers
* Invalid JSON/form mixing
* Wrong HTTP methods
* Empty query params
* Malformed URLs
* Large request bodies

---

## Stretch Goals

* Add `/health` endpoint for system checks
* Add request logging with timestamps
* Add response timing middleware
* Add JSON version of every endpoint (`?format=json`)
* Add rate limiting simulation

---

## Final Outcome

A single Go server that behaves like a **real-world backend API**, demonstrating:

* routing design
* middleware flow
* request validation
* response handling
* sub-router architecture
