package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

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

func main() {
	http.HandleFunc("/method-inspector", methodHandler)
	http.HandleFunc("/validate-body", validatehandler)
	http.HandleFunc("/query-firewall", firewallHandler)

	fmt.Println("Server is running at port 8080")
	http.ListenAndServe(":8080", nil)

}

