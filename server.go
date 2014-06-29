package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
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

type AssetGatherer interface {
	Search(string) ([]Asset, error)
	Init() error
}

var assetGatherers []AssetGatherer

// Performs search for assets from all sources
func performSearchRequest(searchString string) (string, error) {
	var allAssets []Asset
	for _, assetGatherer := range assetGatherers {
		assets, err := assetGatherer.Search(searchString)
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
	// Initialize all the asset gathering plugins
	for _, assetGatherer := range assetGatherers {
		err := assetGatherer.Init()
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	// Register '/search' URI that displays result from all asset gatherers
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
