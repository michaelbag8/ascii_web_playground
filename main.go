package main

import (
	"code/ascii"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
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
var res = template.Must(template.ParseFiles("templates/result.html"))
var tmplReverse = template.Must(template.ParseFiles("tem/reverse.html"))

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
		if len(read) == 0 {
			http.Error(w, "body is empty", http.StatusBadRequest)
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

type ArtData struct {
	Result string
	Error  string
}

func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl.Execute(w, ArtData{})
	} else {
		text := r.FormValue("input")
		banner := r.FormValue("banner")
		if text == "" || banner == "" {
			http.Error(w, "fields are empty", http.StatusBadRequest)
			return
		}

		data, err := ascii.Render(text, banner)

		if err != nil {
			if errors.Is(err, ascii.ErrBannerNotFound) {
				w.WriteHeader(http.StatusNotFound)
				res.Execute(w, ArtData{Error: "banner not found"})
				return
			} else if errors.Is(err, ascii.ErrInvalidChar) {
				w.WriteHeader(http.StatusBadRequest)
				res.Execute(w, ArtData{Error: "invalid character"})
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				res.Execute(w, ArtData{Error: "internal server error"})
				return
			}

		}
		res.Execute(w, ArtData{Result: data})
	}
}

func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler)
	mux.HandleFunc("GET /about", aboutHandler)
	mux.HandleFunc("/echo", echoHandler)
	mux.HandleFunc("GET /banner", bannerHandler)
	mux.HandleFunc("GET /list", listHandler)
	mux.HandleFunc("/reverse", reverseHandler)
	mux.HandleFunc("/books", bookHandler)
	mux.HandleFunc("/greet", greetHandler)
	mux.HandleFunc("/users", userHandler)
	mux.HandleFunc("POST /register", registerHandler)
	mux.HandleFunc("GET /ascii-art", asciiHandler)
	mux.HandleFunc("POST /ascii-art", asciiHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", withLogging(mux))
}
