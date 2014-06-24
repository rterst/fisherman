package main

import (
	"github.com/go-martini/martini"
	"net/http"
    "fmt"
	"encoding/json"
)

// Generic asset struct
type Asset struct {
	Title       string
	Author      string
	Description string
	SourceId    string
	SourceURL   string
	Source      string
	Thumbnail   string
}

var searchFunctions = []func(string)([]Asset, error) {}

// Performs search for assets from all sources
func performSearchRequest(searchString string) (string, error) {
	var allAssets []Asset
	for _, searchFunction := range searchFunctions {
		assets, err := searchFunction(searchString)
		if err != nil {
			fmt.Println("Error executing one of the searches:", err)
		}
		allAssets = append(allAssets, assets...)
	}
	
	encoded, err := json.MarshalIndent(allAssets, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(encoded), nil
}

func main() {
	// Initialize all the search plugins
	YouTubeInit()
	
	m := martini.Classic()
	m.Get("/search", func(req *http.Request, writer http.ResponseWriter) (int, string) {
		result, err := performSearchRequest(req.FormValue("q"))
		
		if err != nil {
			fmt.Println("error:", err)
			return http.StatusInternalServerError, err.Error()
		}

		writer.Header().Set("Content-Type", "application/json")
		return http.StatusOK, result
	})
	m.Run()
}