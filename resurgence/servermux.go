package main

import (
	"fmt"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
func greethandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name != "" {
		fmt.Fprintf(w, "Greetings, %s!", name)
	} else {
		fmt.Fprintln(w, "Greetings, Stranger!")
	}
}

func main() {
	mainMux := http.NewServeMux()

	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/v1/ping", pingHandler)
	apiMux.HandleFunc("/v1/greet", greethandler)

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	mainMux.HandleFunc("/ping", pingHandler)

	http.ListenAndServe(":8080", mainMux)

}
