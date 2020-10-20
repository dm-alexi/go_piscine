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
