package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type order struct {
	Money int    `json:"money"`
	Candy string `json:"candyType"`
	Count int    `json:"candyCount"`
}

type answer struct {
	Thanks string `json:"thanks"`
	Change int    `json:"change"`
	Err    string `json:"error"`
}

func getClient() *http.Client {
	data, err := ioutil.ReadFile("../../minica/minica.pem")
	if err != nil {
		log.Fatalf("Can't read root certificate: %v\n", err)
	}
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)
	config := &tls.Config{
		RootCAs:      cp,
		Certificates: make([]tls.Certificate, 1),
	}
	config.Certificates[0], err = tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Can't load client certificate: %v\n", err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
	return client
}

func main() {
	var ord order
	flag.StringVar(&ord.Candy, "k", "", "Candy type (CE, AA, NT, DE, YR)")
	flag.IntVar(&ord.Money, "m", 0, "Amount of money")
	flag.IntVar(&ord.Count, "c", 0, "Number of candies")
	flag.Parse()
	requestBody, err := json.Marshal(ord)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	client := getClient()
	response, err := client.Post("https://localhost:3333/buy_candy", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	var ans answer
	err = json.Unmarshal(body, &ans)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	if ans.Err != "" {
		fmt.Printf("Error: %s\n", ans.Err)
	} else {
		fmt.Printf("%s Your change is %d\n", ans.Thanks, ans.Change)
	}
}
