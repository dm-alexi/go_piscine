package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

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
