package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func jwtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("alexisawesome"))
	if err != nil {
		log.Fatalf("Error creating token: %s", err)
	}
	js := JSONtoken{
		Token: tokenString,
	}
	json.NewEncoder(w).Encode(js)
}
