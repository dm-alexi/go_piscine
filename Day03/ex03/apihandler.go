package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	page := 1
	if paging, ok := r.URL.Query()["page"]; ok {
		page, err = strconv.Atoi(paging[0])
		if err != nil || page < 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(JSONerror{Err: fmt.Sprintf("Invalid 'page' value: %s\n", paging[0])})
			return
		}
	}
	var req requestAll
	req.es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	places, hits, err := req.GetPlaces(numEntries, numEntries*(page-1))
	if err != nil {
		log.Fatalf("Error processing the request: %s", err)
	}
	lastPage := (hits + numEntries - 1) / numEntries
	if page > lastPage {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JSONerror{Err: fmt.Sprintf("Invalid 'page' value: %d\n", page)})
		return
	}
	w.WriteHeader(http.StatusCreated)
	js := JSONplaces{
		Name:   "places",
		Total:  hits,
		Places: places,
		Prev:   page - 1,
		Next:   page + 1,
		Last:   lastPage,
	}
	if page == 1 {
		js.Prev = lastPage
	}
	if page == lastPage {
		js.Next = 1
	}
	json.NewEncoder(w).Encode(js)
}
