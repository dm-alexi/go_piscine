package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type geopoint struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type entry struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location geopoint `json:"location"`
}

func loadData(filename string) []entry {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening csv file: %s", err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var entries []entry
	r.ReadString('\n')
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		t := strings.Split(line, "\t")
		if len(t) != 6 {
			log.Fatalf("Error data format: %s", line)
		}
		id, err1 := strconv.Atoi(t[0])
		lon, err2 := strconv.ParseFloat(t[4], 64)
		lat, err3 := strconv.ParseFloat(strings.Trim(t[5], "\n"), 64)
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Println(err3)
			log.Fatalf("Error data format: %s", line)
		}
		entries = append(entries, entry{ID: id + 1, Name: t[1], Address: t[2], Phone: t[3], Location: geopoint{Longitude: lon, Latitude: lat}})
	}
	if err != nil {
		log.Fatalf("Error reading data file: %s", err)
	}
	return entries
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
	settings := `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":1
		},
		"mappings":{
			"properties":{
				"name":{
					"type":"text"
				},
				"address":{
					"type":"text"
				},
				"phone":{
					"type":"text"
				},
				"location":{
					"type":"geo_point"
				}
			}
		}
	}`
	res, err = es.Indices.Create("places", es.Indices.Create.WithBody(strings.NewReader(settings)))
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "places",
		Client: es,
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}
	for _, ent := range loadData("data.csv") {
		data, err := json.Marshal(ent)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", ent.ID, err)
		}
		err = indexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.Itoa(ent.ID), //strconv.Itoa(i + 1),
				Body:       bytes.NewReader(data),
			},
		)
		if err != nil {
			log.Fatalf("Indexer add error: %s", err)
		}
	}
	if err := indexer.Close(context.Background()); err != nil {
		log.Fatalf("Indexer close error: %s", err)
	}
}
