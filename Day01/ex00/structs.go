package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"strings"
)

// Recipes is an aggregate type structure, comprising both xml and json representations
type Recipes struct {
	XMLName xml.Name `xml:"recipes"`
	Cake    []struct {
		Name      string `xml:"name" json:"name"`
		Stovetime string `xml:"stovetime" json:"time"`
		Item      []struct {
			Itemname  string `xml:"itemname" json:"ingredient_name"`
			Itemcount string `xml:"itemcount" json:"ingredient_count"`
			Itemunit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
		} `xml:"ingredients>item" json:"ingredients"`
	} `xml:"cake" json:"cake"`
}

// DBReader is reading xml and json files to Recipe structure
type DBReader interface {
	ReadDB(filename string) Recipes
}

// XMLReader implements reading from xml file to Recipe
type XMLReader struct {
}

// JSONReader implements reading from json file to Recipe
type JSONReader struct {
}

// ReadDB from xml file
func (o *XMLReader) ReadDB(filename string) (Recipes, error) {
	var record Recipes
	in, err := ioutil.ReadFile(filename)
	if err == nil {
		err = xml.Unmarshal(in, &record)
	}
	return record, err
}

// ReadDB from json file
func (o *JSONReader) ReadDB(filename string) (Recipes, error) {
	var record Recipes
	in, err := ioutil.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(in, &record)
	}
	return record, err
}

func readUniversal(filename string) (Recipes, error) {
	if len(filename) > 4 && strings.HasSuffix(filename, ".xml") {
		return new(XMLReader).ReadDB(filename)
	} else if len(filename) > 5 && strings.HasSuffix(filename, ".json") {
		return new(JSONReader).ReadDB(filename)
	}
	var record Recipes
	return record, errors.New("invalid file " + filename)
}
