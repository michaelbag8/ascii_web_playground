Here are four focused Go exercises that use only the standard `net/http` package (plus other standard library packages such as `html/template` and `encoding/json`).

---

## Task 1: Server → Browser (Render HTML Template)

### Goal

Create a web server that displays a list of books in a browser.

### Requirements

1. Create a `Book` struct:

```go
type Book struct {
    Title  string
    Author string
}
```

2. Create a slice with at least 3 books.
3. Create an HTTP handler at `/books`.
4. Pass the slice to an HTML template.
5. Render the books as an unordered list (`<ul>`).

### Expected URL

```text
http://localhost:8080/books
```

### Example Output

```html
<h1>Books</h1>
<ul>
    <li>Go Programming - John Doe</li>
    <li>Learning HTTP - Jane Smith</li>
    <li>Mastering APIs - Alex Brown</li>
</ul>
```

### Concepts Practiced

* `http.HandleFunc`
* `html/template`
* Passing data from server to browser

---

## Task 2: Browser → Server (Read Query Parameters)

### Goal

Build a greeting endpoint that reads data from the URL.

### Requirements

1. Create a handler at `/greet`.
2. Read the query parameter `name`.
3. If provided, return:

```text
Hello, Alice!
```

4. If missing, return:

```text
Hello, Guest!
```

### Example URLs

```text
http://localhost:8080/greet?name=Alice
```

Response:

```text
Hello, Alice!
```

---

```text
http://localhost:8080/greet
```

Response:

```text
Hello, Guest!
```

### Concepts Practiced

* `r.URL.Query()`
* `Get()`
* Reading browser input

---

## Task 3: Struct → JSON

### Goal

Convert a Go struct into JSON and send it to the browser.

### Requirements

Create:

```go
type User struct {
    ID    int
    Name  string
    Email string
}
```

1. Create a sample user.
2. Create a handler at `/user`.
3. Set:

```go
w.Header().Set("Content-Type", "application/json")
```

4. Encode the struct as JSON and write it to the response.

### Expected URL

```text
http://localhost:8080/user
```

### Expected Response

```json
{
  "ID": 1,
  "Name": "Alice",
  "Email": "alice@example.com"
}
```

### Concepts Practiced

* `encoding/json`
* `json.NewEncoder`
* Sending JSON responses

---

## Task 4: JSON → Struct

### Goal

Accept JSON from the browser and convert it into a Go struct.

### Requirements

Use:

```go
type User struct {
    Name string
    Age  int
}
```

1. Create a POST endpoint `/register`.
2. Read JSON from the request body.
3. Decode it into a `User` struct.
4. Return a confirmation message using the decoded values.

### Example Request

```http
POST /register
Content-Type: application/json

{
    "Name": "David",
    "Age": 25
}
```

### Example Response

```text
User David registered successfully. Age: 25
```

### Concepts Practiced

* `json.NewDecoder`
* `r.Body`
* Receiving JSON from clients

---

### Progression

These four tasks follow the exact data flow you'll use in most web applications:

1. **Server → Browser** → Render HTML templates.
2. **Browser → Server** → Read query parameters.
3. **Struct → JSON** → Build APIs.
4. **JSON → Struct** → Accept API requests.

Together, they cover the fundamental request/response patterns used in Go web development with the standard library.
