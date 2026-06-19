package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Exercise 1
func pongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

// Exercise 2
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name != "" {
		fmt.Fprintf(w, "Hello, %s!", name)
	} else {
		fmt.Fprintf(w, "Hello, Guest!")
	}
}

// Exercise 3
func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading the body", http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, len(data))

	} else {
		fmt.Fprintln(w, "Send a POST request with text to count words")
	}

}

// Exercise 4
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	op := r.URL.Query().Get("op")
	a := r.URL.Query().Get("a")
	b := r.URL.Query().Get("b")

	conA, err := strconv.Atoi(a)
	if err != nil {
		http.Error(w, "error parsing parameter a", http.StatusBadRequest)
		return
	}
	conB, err := strconv.Atoi(b)
	if err != nil {
		http.Error(w, "error parsing parameter b", http.StatusBadRequest)
		return
	}

	switch op {
	case "add":
		fmt.Fprintln(w, "Result:", conA+conB)
	case "subtract":
		fmt.Fprintln(w, "Result:", conA-conB)
	case "multiply":
		fmt.Fprintln(w, "Result:", conA*conB)
	default:
		fmt.Fprintln(w, "operation is unknown", http.StatusBadRequest)
		return
	}

}

func agentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/agent" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, r.Header)
}
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	APIKey := r.Header.Get("X-API-Key")
	userKey := "secret123"

	if userKey != APIKey {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}

	fmt.Fprintln(w, "Welcome to your dasdboard")

}

func legacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/legacy" {
		http.Redirect(w, r, "/v2", http.StatusMovedPermanently)
	}

}
func v2Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to version 2")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", pongHandler)
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/count", countHandler)
	mux.HandleFunc("/calculate", calculateHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/agent", agentHandler)
	mux.HandleFunc("/legacy", legacyHandler)
	mux.HandleFunc("/v2", v2Handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
