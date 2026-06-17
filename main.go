package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type DataPage struct {
	Title string
	Items []string
}
var tmpl =template.Must(template.ParseFiles("templates/index.html"))
	
func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		fmt.Fprintln(w, "Hello, Stranger")
	} else {
		fmt.Fprintf(w, "Hello %s\n", name)
	}

}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		read, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		data := string(read)
		fmt.Fprintln(w, data)

	} else {
		fmt.Fprintln(w, "Send me a POST request with a body to echo it back")
	}
}

func bannerHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("banners/standard.txt")
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/list"{
		http.Error(w, "Wrong Path", http.StatusBadRequest)
		return
	}
	
	result := DataPage{
		Title: "Banner Title",
		Items: []string{"Apple","Mango","Guava", "Coconut"},
	}

	err := tmpl.Execute(w, result)
	if err != nil {
		log.Println(w, "render error", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is the about page")
}

func main() {
	http.HandleFunc("GET /{$}", handler)
	http.HandleFunc("GET /about", aboutHandler)
	http.HandleFunc("GET /echo", echoHandler)
	http.HandleFunc("GET /banner", bannerHandler)
	http.HandleFunc("GET /list", listHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
