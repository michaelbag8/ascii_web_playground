package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var temp = template.Must(template.ParseFiles("temp/index.html"))
var tempest = template.Must(template.ParseFiles("temp/layout.html"))

type Book struct {
	Title  string
	Author string
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/books" {
		log.Println("wrong path")
	}

	books := []Book{
		{
			Title:  "Golang For Idiots",
			Author: "Michael Bag",
		},
		{
			Title:  "JS Playground For Dummies",
			Author: "James Baba",
		},
		{
			Title:  "Great To Good",
			Author: "Stephen Covey",
		},
	}
	if err := temp.Execute(w, books); err != nil {
		log.Println("error reading books")
	}
}

type Greet struct {
	Name string
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	parameters := r.URL.Query().Get("name")

	if parameters != "" {
		fmt.Fprintf(w, "Hello, %s\n", parameters)
		return
	} else {
		fmt.Fprintln(w, "Hello , Guest!")
	}

}

type User struct {
	ID    int
	Name  string
	Email string
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := User{
		ID:    1,
		Name:  "John Doe",
		Email: "michaelbag8@gmail.com",
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Bad request", http.StatusInternalServerError)
		return
	}

}

type Users struct {
	Name string
	Age  int
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req Users
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusInternalServerError)
		return
	}

	//json.NewEncoder(w).Encode(req)
	fmt.Fprintf(w,"%s of %d have successfully registered\n",req.Name, req.Age)

}
