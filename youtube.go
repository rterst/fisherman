package main

import (
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
)

// Parent structure
type YouTubeAssetGatherer struct {
	apiKey string
}

// Structure used for parsing search results from server
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

// Custom error object for this plugin
type YouTubeAssetGathererError struct {
    msg string
}
func (e YouTubeAssetGathererError) Error() string { return "[YouTube plugin] "+e.msg }

// Converts youTubeSearchResult list into Assets list
func (this *YouTubeAssetGatherer) convertToAssets(data youTubeSearchResult) []Asset{
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

// Main function - performs the search on YouTube based on query string
// and returns the list of assets matching it
func (this *YouTubeAssetGatherer) Search(query string) ([]Asset, error) {
	var searchUrl = "https://www.googleapis.com/youtube/v3/search"+
		"?part=snippet&"+
		"key="+this.apiKey
	response, err := http.Get(searchUrl+"&q="+query)
    if err != nil {
        return nil, YouTubeAssetGathererError{err.Error()}
    }
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, YouTubeAssetGathererError{err.Error()}
	}
	
	var data youTubeSearchResult
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, YouTubeAssetGathererError{err.Error()}
	}
	assets := this.convertToAssets(data)
	
	return assets, nil
}

// Initialization function
func (this *YouTubeAssetGatherer) Init() (error){
	// Fetch apiKey from config file
	type configuration struct {
    	Key    string
	}
	file, err := os.Open("youtube.cfg")
	if err != nil {
		return YouTubeAssetGathererError{err.Error()}
	}
	decoder := json.NewDecoder(file)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return YouTubeAssetGathererError{err.Error()}
	}

	// Save apiKey for further use
	this.apiKey = conf.Key

	return nil
}

func init() {
	// Register plugin with the server
	assetGatherers = append(assetGatherers, &YouTubeAssetGatherer{})
}