package main

type geopoint struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

// Place is the entry structure
type Place struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location geopoint `json:"location"`
}

// JSONplaces is the return value
type JSONplaces struct {
	Name   string  `json:"name"`
	Total  int     `json:"total"`
	Places []Place `json:"places"`
	Prev   int     `json:"prev_page"`
	Next   int     `json:"next_page"`
	Last   int     `json:"last_page"`
}

// JSONerror is the error value
type JSONerror struct {
	Err string `json:"error"`
}
