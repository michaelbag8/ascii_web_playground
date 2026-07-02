package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"log"
	"net/http"
	"text/template"
)

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

func fetchData(url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return errors.New("error fetching API")

	}
	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return errors.New("error fetching API")
	// }
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return errors.New("error fetching API")
	}

	// for _, artist := range artists{
	// 	fmt.Printf("Name: %s|Created: %d|First Album: %s|Location: %s\n",artist.Name, artist.CreationDate,artist.FirstAlbum, artist.Locations)
	// }

	return nil
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "wrong path", http.StatusNotFound)
		return
	}
	templ, _ := template.ParseFiles("index.html")

	templ.Execute(w, artists)

}

func artistDetail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}
	idC, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "failed", http.StatusBadRequest)
		return
	}
	temps, err := template.ParseFiles("artist.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	//var found *Artist
	found := false
	for _, artist := range artists {
		if artist.ID == idC {
			found = true
			err = temps.Execute(w, artist)
			if err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
				return
			}

		}
	}

	if !found {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

}

func searchHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path !="/search"{
		http.Error(w, "wrong route", http.StatusBadRequest)
		return
	}
	q := r.URL.Query().Get("q")
	if r.Method !=http.MethodGet{
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if strings.TrimSpace(q) ==""{
		http.Error(w, "search query is required", http.StatusBadRequest)
		return
	} 

	var result []Artist
		templ, _ := template.ParseFiles("index.html")
	for _, artist :=range artists{
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(q)){
			result = append(result, artist)
		}
		
	}
	err := templ.Execute(w, result)
			if err != nil {
				http.Error(w, "template error", http.StatusInternalServerError)
				return
			}
}

func filterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	creationDateStr := r.URL.Query().Get("creationDate")
	memberCountStr := r.URL.Query().Get("memberCount")

	if creationDateStr == "" && memberCountStr == "" {
		http.Error(w, "at least one filter is required", http.StatusBadRequest)
		return
	}

	var (
		creationDate int
		memberCount  int
		err          error
	)

	if creationDateStr != "" {
		creationDate, err = strconv.Atoi(creationDateStr)
		if err != nil {
			http.Error(w, "invalid creationDate", http.StatusBadRequest)
			return
		}
	}

	if memberCountStr != "" {
		memberCount, err = strconv.Atoi(memberCountStr)
		if err != nil {
			http.Error(w, "invalid memberCount", http.StatusBadRequest)
			return
		}
	}

	var result []Artist

	for _, artist := range artists {

		match := true

		if creationDateStr != "" && artist.CreationDate != creationDate {
			match = false
		}

		if memberCountStr != "" && len(artist.Members) != memberCount {
			match = false
		}

		if match {
			result = append(result, artist)
		}
	}

	templ, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, result)
	if err != nil {
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func main() {
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/artist", artistDetail)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/filter", filterHandle)
	log.Fatal(http.ListenAndServe(":8080", nil))
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