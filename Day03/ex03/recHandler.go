package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
)

func recHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	var req requestRecs
	req.es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	if lon, ok := r.URL.Query()["lon"]; ok {
		req.lon, err = strconv.ParseFloat(lon[0], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSONerror{Err: fmt.Sprintf("Invalid 'lon' value: %s\n", lon[0])})
			return
		}
	}
	if lat, ok := r.URL.Query()["lat"]; ok {
		req.lat, err = strconv.ParseFloat(lat[0], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSONerror{Err: fmt.Sprintf("Invalid 'lat' value: %s\n", lat[0])})
			return
		}
	}
	places, _, err := req.GetPlaces(numRecs, 0)
	if err != nil {
		log.Fatalf("Error processing the request: %s", err)
	}
	w.WriteHeader(http.StatusCreated)
	js := JSONrec{
		Name:   "Recommendations",
		Places: places,
	}
	json.NewEncoder(w).Encode(js)
}
