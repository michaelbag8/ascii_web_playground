package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

// Knowledge Gap - Go template, nested struct and json
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	Image        string   `json:"image"`
	Locations    string   `json:"locations"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Locations struct {
	Index []Location `json:"index"`
}

type Dates struct {
	Index []Date `json:"index"`
}
type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

var artists []Artist
var locations Locations
var dates Dates
var relations Relations

func fetchData(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}
	return nil
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	templ, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, artists)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	templ, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	found := false
	for _, artist := range artists {
		if artist.ID == id {
			found = true
			err = templ.Execute(w, artist)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}
	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query().Get("q")

	var result []Artist

	templ, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(q)) {

			result = append(result, artist)

		}
	}

	err = templ.Execute(w, result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
/*
func artistDetail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	idC, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var foundArtist *Artist

	for _, artist := range artists {
		if artist.ID == idC {
			foundArtist = &artist
			break
		}
	}

	if foundArtist == nil {
		http.Error(w, "artist not found", http.StatusNotFound)
		return
	}

	temps, err := template.ParseFiles("artist.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	err = temps.Execute(w, foundArtist)
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
*/

func main() {
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		log.Fatal("error fetching artists", err)
	}
	err = fetchData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	if err != nil {
		log.Fatal("error fetching loactions", err)
	}
	err = fetchData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	if err != nil {
		log.Fatal("error fetching dates", err)
	}
	err = fetchData("https://groupietrackers.herokuapp.com/api/relation", &relations)
	if err != nil {
		log.Fatal("error fetching relation", err)
	}
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist", ArtistHandler)
	http.HandleFunc("/search", SearchHandler)
	fmt.Println("server is running at port http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
