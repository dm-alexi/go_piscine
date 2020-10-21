package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/elastic/go-elasticsearch/v7"
)

func middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		line := r.Header.Get("Authorization")
		if len(line) < 7 || !strings.HasPrefix(line, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusBadRequest)
			return
		}
		token, err := jwt.Parse(line[7:], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("alexisawesome"), nil
		})
		fmt.Println(line[7:])
		fmt.Println(*token)
		if err != nil || !token.Valid {
			http.Error(w, fmt.Sprintf("Authorization error: %s\n", err), http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}

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
