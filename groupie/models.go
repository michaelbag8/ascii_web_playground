package main

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
	ID    int `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relations struct{
	Index []Relation `json:"index"`
}