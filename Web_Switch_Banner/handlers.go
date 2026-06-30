package main

import (
	"net/http"
	"text/template"
)

type PageData struct {
	Result string
	Banner string
	Text   string
}

var templ = template.Must(template.ParseFiles("templates/index.html"))

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	err := templ.Execute(w, PageData{})
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "failed to parse form", http.StatusInternalServerError)
		return
	}
	text := r.FormValue("text")
	banner := r.FormValue("banner")

	if text == "" {
		http.Error(w, "text field is empty", http.StatusBadRequest)
		return
	}

	if banner == "" {
		http.Error(w, "banner field is empty", http.StatusNotFound)
		return
	}

	bannerFile := "banners/" + banner + ".txt"
	data, err := LoadBanner(bannerFile)
	if err != nil {
		http.Error(w, "loading banner file failed", http.StatusInternalServerError)
		return
	}
	result := GenerateArt(text, data)

	final := PageData{
		Result: result,
		Banner: banner,
		Text:   text,
	}

	err = templ.Execute(w, final)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}

func handleSwitch(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	banner := r.URL.Query().Get("banner")

	if text == "" {
		http.Error(w, "text field is empty", http.StatusBadRequest)
		return
	}

	if banner == "" {
		http.Error(w, "banner field is empty", http.StatusNotFound)
		return
	}

	bannerFile := "banners/" + banner + ".txt"
	data, err := LoadBanner(bannerFile)
	if err != nil {
		http.Error(w, "loading banner file failed", http.StatusInternalServerError)
		return
	}
	result := GenerateArt(text, data)

	final := PageData{
		Result: result,
		Banner: banner,
		Text:   text,
	}

	err = templ.Execute(w, final)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
}
