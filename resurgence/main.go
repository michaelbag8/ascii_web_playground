package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func methodHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You made a %s request\n", r.Method)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "fail to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(data) == 0 {
		http.Error(w, "body cannot be empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, string(data))

}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	head := r.Header.Get("X-Custom-Token")
	if head == "" {
		http.Error(w, "X-Custom-Token header is missing", http.StatusBadRequest)
		return
	}
	contentType := r.Header.Get("Content-Type")

	if contentType == ""{
		contentType = "Content-Type not provided"
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Token received: %s", head)
	fmt.Fprintf(w, " Content-Type: %s",contentType)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	// if contentType != "application/x-www-form-urlencoded"{
	// 	http.Error(w, "Content-Type is not application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
	// 	return
	// }
	if !strings.HasPrefix(contentType, "application/x-www-form-urlencoded"){
		http.Error(w, "Content-Type is not application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}

	err := r.ParseForm()
	if err !=nil{
		http.Error(w, "form parsing failed", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	language := r.FormValue("language")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	if language == "" {
		http.Error(w, "language is required", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello %s, you are coding in %s", username, language)

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code parameter is required", http.StatusNotFound)
		return
	}
	convCode, err := strconv.Atoi(code)
	if err != nil {
		http.Error(w, "code must be a valid integer", http.StatusNotFound)
		return
	}
	if convCode >= 100 || convCode <= 599 {
		http.Error(w, "code must be a valid HTTP status code (100–599)", http.StatusNotFound)
		return
	}
	w.WriteHeader(convCode)
	fmt.Fprintf(w, "Responding with status %d", convCode)
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	const tmplStr = `
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
  <h1>{{.Title}}</h1>
  <p>{{.Body}}</p>
</body>
</html>
`
	type PageData struct {
		Title string
		Body  string
	}
	var temp = template.Must(template.New("page").Parse(tmplStr))
	r.ParseForm()
	title := r.FormValue("title")
	body := r.FormValue("body")

	if title == "" || body == "" {
		http.Error(w, "title and body are required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := temp.Execute(w, PageData{Title: title, Body: body}); err != nil {
		http.Error(w, "template execution failed", http.StatusInternalServerError)
		return
	}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/method-inspector", methodHandler)
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("/headers", headerHandler)
	mux.HandleFunc("/form", formHandler)
	mux.HandleFunc("/render", renderHandler)

	fmt.Println("server is running at port 8080")

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}
