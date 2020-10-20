package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

const numEntries = 10

// Store is a database access wrapper
type Store interface {
	// turns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type requestAll struct {
	es *elasticsearch.Client
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

func getEntryList(places []Place) string {
	var sb strings.Builder
	sb.WriteString("<ul>\n")
	for _, place := range places {
		sb.WriteString(fmt.Sprintf("<li><div>%s</div><div>%s</div><div>%s</div></li>\n", place.Name, place.Address, place.Phone))
	}
	sb.WriteString("</ul>\n")
	return sb.String()
}

func getLinks(page int, last int) string {
	var sb strings.Builder
	if page > 1 {
		sb.WriteString(fmt.Sprintf("<a href=\"/?page=%d\">Previous</a>\n", page-1))
	}
	if page < last {
		sb.WriteString(fmt.Sprintf("<a href=\"/?page=%d\">Next</a>\n", page+1))
		sb.WriteString(fmt.Sprintf("<a href=\"/?page=%d\">Last</a>\n", last))
	}
	return sb.String()
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	var body string
	var err error
	page := 1
	if paging, ok := r.URL.Query()["page"]; ok {
		page, err = strconv.Atoi(paging[0])
		if err != nil || page < 1 {
			http.Error(w, fmt.Sprintf("Invalid 'page' value: %s\n", paging[0]), http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf("Invalid 'page' value: %d\n", page), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, `<!doctype html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>Places</title>
		<meta name="description" content="">
		<meta name="viewport" content="width=device-width, initial-scale=1">
	</head>`)
	body = getEntryList(places) + getLinks(page, lastPage)
	fmt.Fprintf(w, "<body><h5>Total: %d</h5>%s</body></html>", hits, body)
}

func main() {
	PORT := ":8888"
	http.HandleFunc("/", httpHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
