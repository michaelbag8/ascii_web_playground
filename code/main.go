package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)
type Page struct{
	Username string
	Language string
}
// switch, if, and map
func secureHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodPost:
		fmt.Fprintf(w, "%s request method accepted\n", r.Method)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s request is not supported\n", r.Method)
		return
	}

}

// map variations
var supported = map[string]bool{
	http.MethodGet:  true,
	http.MethodPost: true,
}

func methodHandler(w http.ResponseWriter, r *http.Request) {
	if supported[r.Method] {
		fmt.Fprintf(w, "%s request method accepted\n", r.Method)
		return
	}
	w.Header().Set("Allow", "POST, GET")
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "%s request is not supported\n", r.Method)

}

func validatehandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to parse body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		body := string(data)
		if len(body) == 0 {
			http.Error(w, "body is required", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(body) == "" {
			http.Error(w, "body cannot be blank",http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Valid body received: %s\n", body)

	} else{
		http.Error(w, "method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func validateBodyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to parse body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		body := string(data)

		if len(body) == 0 {
			http.Error(w, "body is required", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(body) == "" {
			http.Error(w, "body cannot be blank", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "Valid body received: %s\n", body)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Please send a POST request")
		return
	}

}

func firewallHandler(w http.ResponseWriter, r *http.Request){
	user :=r.URL.Query().Get("user")
	role :=r.URL.Query().Get("role")

	if user == ""{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "user parameter missing")
		return
	}
	if role == ""{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "role parameter missing")
		return
		
	}

	if role != "admin" && role !="user"{
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w,"Only admin and user are accepted")
	}

	fmt.Fprintf(w, "User: %s, Role: %s", user, role)
}

func formHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")

	if contentType!="application/json"{
		http.Error(w, "Only JSON is accepted", http.StatusUnsupportedMediaType)
		return
	}
	var req Page

	err := json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Language == ""{
		http.Error(w, "fields cannot be empty", http.StatusBadRequest)
		return
	}


	fmt.Fprintf(w,"Hello %s, %s is awesome!", req.Username, req.Language)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "login endpoint reached")
}
func logoutHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "logout endpoint reached")
}
func infoHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "data info endpoint reached")
}


func main() {
	http.HandleFunc("/method-inspector", methodHandler)
	http.HandleFunc("/validate-body", validatehandler)
	http.HandleFunc("/query-firewall", firewallHandler)
	http.HandleFunc("/form", formHandler)
	
	/*
	mainMux := http.NewServeMux()
	apiMux:= http.NewServeMux()
	
	apiMux.HandleFunc("/v1/auth/login", loginHandler)
	apiMux.HandleFunc("/v1/auth/logout", logoutHandler)
	apiMux.HandleFunc("/v1/data/info", infoHandler)

	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))

	http.ListenAndServe(":8080", mainMux)
	*/

	fmt.Println("Server is running at port 8080")
	http.ListenAndServe(":8080", nil)


}

