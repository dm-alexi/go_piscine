package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
)

const numEntries = 10
const numRecs = 3

// Store is a database access wrapper
type Store interface {
	// turns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type requestAll struct {
	es *elasticsearch.Client
}

type requestRecs struct {
	es  *elasticsearch.Client
	lat float64
	lon float64
}

func (o *requestRecs) GetPlaces(limit int, offset int) ([]Place, int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": map[string]interface{}{
			"_geo_distance": map[string]interface{}{
				"location": map[string]interface{}{
					"lat": o.lat,
					"lon": o.lon,
				},
				"order":           "asc",
				"unit":            "km",
				"mode":            "min",
				"distance_type":   "arc",
				"ignore_unmapped": true,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := o.es.Search(
		o.es.Search.WithContext(context.Background()),
		o.es.Search.WithIndex("places"),
		o.es.Search.WithBody(&buf),
		o.es.Search.WithTrackTotalHits(true),
		o.es.Search.WithPretty(),
		o.es.Search.WithSize(limit),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	var r map[string]interface{}
	places := make([]Place, limit)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	for i, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		js, err := json.Marshal(hit.(map[string]interface{})["_source"])
		if err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		if err = json.Unmarshal(js, &places[i]); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
	}
	return places, int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), err
}

func (o *requestAll) GetPlaces(limit int, offset int) ([]Place, int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := o.es.Search(
		o.es.Search.WithContext(context.Background()),
		o.es.Search.WithIndex("places"),
		o.es.Search.WithBody(&buf),
		o.es.Search.WithTrackTotalHits(true),
		o.es.Search.WithPretty(),
		o.es.Search.WithFrom(offset),
		o.es.Search.WithSize(limit),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	var r map[string]interface{}
	places := make([]Place, limit)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	for i, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		js, err := json.Marshal(hit.(map[string]interface{})["_source"])
		if err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		if err = json.Unmarshal(js, &places[i]); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
	}
	return places, int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)), err
}

func main() {
	PORT := ":8888"
	http.HandleFunc("/api/places/", apiHandler)
	http.HandleFunc("/api/recommend", recHandler)
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
