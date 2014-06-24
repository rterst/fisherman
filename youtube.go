package main

import (
	"net/http"
//	"fmt"
	"io/ioutil"
	"encoding/json"
)
	
// Structure for parsing YouTube search results from JSON
type youTubeSearchResult struct {
	Items []struct {
		Id struct {
			Videoid string
		}
		Snippet struct {
			Title string
			Description string
			Thumbnails struct {
				Default struct {
					Url string
				}
			}
			ChannelTitle string
		}
	}
}

// Converts youTubeSearchResult list into Assets list
func youTubeConvertToAssets(data youTubeSearchResult) []Asset{
	var ret []Asset

	for _, item := range data.Items {
		asset := Asset {
			Title: item.Snippet.Title,
			Author: item.Snippet.ChannelTitle,
			Description: item.Snippet.Description,
			SourceId: item.Id.Videoid,
			SourceURL: "https://www.youtube.com/watch?v="+item.Id.Videoid,
			Source: "youtube",
			Thumbnail: item.Snippet.Thumbnails.Default.Url,
		}
		ret = append(ret, asset)
	}
	
	return ret
}

//Main function - performs the search on YouTube based on query string
// and returns the list of assets matching it
func youTubePerformSearch(query string) ([]Asset, error) {
	// TODO: Hard-coding API key is a bad idea, it should be in config file
	var searchUrl = "https://www.googleapis.com/youtube/v3/search"+
		"?part=snippet&"+
		"key=AIzaSyDi74Dekt6FUhiK6K9c52Y01avYjvTgIto"
	response, err := http.Get(searchUrl+"&q="+query)
    if err != nil {
        return nil, err
    }
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	var data youTubeSearchResult
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}
	assets := youTubeConvertToAssets(data)
	
	return assets, nil
}

func YouTubeInit() {
	searchFunctions = append(searchFunctions, youTubePerformSearch)
}