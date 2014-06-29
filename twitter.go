package main

import (
	"net/http"
	"net/url"
	"encoding/base64"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"os"
)

// Parent structure
type TwitterAssetGatherer struct {
	accessToken string
}

// Structure used for parsing authentication response from server
type twitterAuthResponse struct {
	Token_type   string
	Access_token string
}

// Structure used for parsing search response from server
type twitterSearchResult struct {
	Statuses []struct {
		Id_str string
		Text string
		User struct {
			Profile_image_url string
			Name string
			Id_str string
		}
	}
}

// Custom error object for this plugin
type TwitterAssetGathererError struct {
    msg string
}
func (e TwitterAssetGathererError) Error() string { return "[Twitter plugin] "+e.msg }

// Converts twitterSearchResult list into Assets list
func (this *TwitterAssetGatherer) convertToAssets(data twitterSearchResult) []Asset{
	var ret []Asset

	for _, item := range data.Statuses {
		asset := Asset {
			Title: item.Text,
			Author: item.User.Name,
			Description: item.Text,
			SourceId: item.Id_str,
			SourceURL: "http://twitter.com/"+item.User.Id_str+"/status/"+item.Id_str,
			Source: "twitter",
			Thumbnail: item.User.Profile_image_url,
		}
		ret = append(ret, asset)
	}
	
	return ret
}

// Main function - performs the search on Twitter based on query string
// and returns the list of assets matching it
func (this *TwitterAssetGatherer) Search(query string) ([]Asset, error) {
	client := &http.Client{}
	request, err := http.NewRequest(
		"GET",
		"https://api.twitter.com/1.1/search/tweets.json"+"?q="+query+"&count=5",
		nil)
	request.Header.Add("Authorization", "Bearer "+this.accessToken)
	if err != nil {
		return nil, TwitterAssetGathererError{err.Error()}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, TwitterAssetGathererError{err.Error()}
	}
	
	responseJSON, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, TwitterAssetGathererError{err.Error()}
	}
	if response.Status != "200 OK" {
		return nil, TwitterAssetGathererError{"Search failed. Invalid response from server: "+response.Status+
			"\n"+string(responseJSON)}
	}

	var data twitterSearchResult
	err = json.Unmarshal(responseJSON, &data)
	if err != nil {
		return nil, err
	}

	assets := this.convertToAssets(data)

	return assets, nil
}

// Initialization function
func (this *TwitterAssetGatherer) Init() (error) {
	// Fetch apiKey and apiSecret from config file
	type configuration struct {
    	ApiKey    string
    	ApiSecret string
	}
	file, err := os.Open("twitter.cfg")
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}
	decoder := json.NewDecoder(file)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}

	// Obtain authentication token from twitter server
	bearerTokenBase64 := base64.StdEncoding.EncodeToString([]byte(conf.ApiKey+":"+conf.ApiSecret))
	client := &http.Client{}
	parameters := url.Values{}
	parameters.Set("grant_type", "client_credentials")
	request, err := http.NewRequest(
		"POST",
		"https://api.twitter.com/oauth2/token",
		bytes.NewBufferString(parameters.Encode()))
	request.Header.Add("Authorization", "Basic "+bearerTokenBase64)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}

	response, err := client.Do(request)
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}
	if response.Status != "200 OK" {
		return TwitterAssetGathererError{"Auth failed. Invalid response from server"}
	}
	responseJSON, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}
	var data twitterAuthResponse
	err = json.Unmarshal(responseJSON, &data)
	if err != nil {
		return TwitterAssetGathererError{err.Error()}
	}
	if data.Token_type != "bearer" {
		return TwitterAssetGathererError{"Auth failed. Invalid response from server"}
	}

	// Save authentication token for further use
	this.accessToken = data.Access_token

	return nil
}

func init() {
	// Register plugin with the server
	assetGatherers = append(assetGatherers, &TwitterAssetGatherer{})
}