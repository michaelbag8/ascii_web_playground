package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type DataPage struct {
	Title string
	Items []string
}

type Data struct {
	Result string
	Error  string
}

var tmpl = template.Must(template.ParseFiles("templates/index.html"))
var tmplReverse = template.Must(template.ParseFiles("templates/reverse.html"))

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

// hello   你好 >> 好你  olleh
func reverseWord(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// hello   你 好 >> olleh 好你
func reversEachWord(str string) string {
	words := strings.Fields(str)
	for i, word := range words {
		words[i] = reverseWord(word)
	}
	return strings.Join(words, " ")
}
func reverseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmplReverse.Execute(w, Data{})
	} else {
		r.ParseForm()
		word := r.FormValue("word")
		if word == "" {
			data := Data{Error: "Please enter a word."}
			if err := tmplReverse.Execute(w, data); err != nil {
				log.Println(err)
			}
			return
		}
		rev := reverseWord(word)

		userInfo := Data{
			Result: rev,
			Error:  "",
		}
		if err := tmplReverse.Execute(w, userInfo); err != nil {
			log.Println("error parsing file")
		}
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
	if r.URL.Path != "/list" {
		http.Error(w, "Wrong Path", http.StatusBadRequest)
		return
	}

	result := DataPage{
		Title: "Banner Title",
		Items: []string{"Apple", "Mango", "Guava", "Coconut"},
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
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("GET /banner", bannerHandler)
	http.HandleFunc("GET /list", listHandler)
	http.HandleFunc("/reverse", reverseHandler)
	http.HandleFunc("/books", bookHandler)
	http.HandleFunc("/greet", greetHandler)
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("POST /register", registerHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}
