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

// Performs search for assets from all sources
func performSearchRequest(search_string string) (string, error) {
//TODO: rewrite using plugin architecture
	assets_youtube, err := YouTube_PerformSearch(search_string)
	if err != nil {
		return "", err
	}
	
	encoded, err := json.MarshalIndent(assets_youtube, "", "  ")
	if err != nil {
		return "", err
	}
	
	return string(encoded), nil
}

func main() {
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