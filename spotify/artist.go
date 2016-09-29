package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ArtistList struct {
	ArtistItem []Artist `json:"items"`
}

type ArtistResponse struct {
	Artists ArtistList `json:"artists"`
}

func GetArtistsByName(name string) ([]Artist, error) {
	client := &http.Client{}

	// url parameters
	params := url.Values{}
	params.Add("q", name)
	params.Add("type", "artist")

	url := fmt.Sprintf("https://api.spotify.com/v1/search?=%v", params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("GetArtistsByName error creating request: %v", err)
	}

	<-GlobalRateLimiter
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetArtistsByName request error: %v", err)
	}

	// try to UnMarshall
	var ar ArtistResponse
	defer resp.Body.Close()
	fullBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetArtistsByName error reading response: %v", err)
	}

	_ = json.Unmarshal(fullBody, &ar)
	if err != nil {
		return nil, fmt.Errorf("GetArtistsByName error unmarshalling artists: %v", err)
	}

	return ar.Artists.ArtistItem, nil
}

func GetArtistById(id string) (Artist, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Artist{}, fmt.Errorf("GetArtistById error creating request: %v", err)
	}

	<-GlobalRateLimiter
	resp, err := client.Do(req)
	if err != nil {
		return Artist{}, fmt.Errorf("GetArtistById request error: %v", err)
	}

	// try to UnMarshall
	var artist Artist
	defer resp.Body.Close()
	fullBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Artist{}, fmt.Errorf("GetArtistById error reading response: %v", err)
	}

	err = json.Unmarshal(fullBody, &artist)
	if err != nil {
		return Artist{}, fmt.Errorf("GetArtistById error unmarshalling artist: %v", err)
	}

	return artist, nil
}

func GetArtistCatalog(id string) (Catalog, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v/albums?market=US&limit=50", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Catalog{}, fmt.Errorf("GetArtistCatalog error creating request: %v", err)
	}

	<-GlobalRateLimiter
	resp, err := client.Do(req)
	if err != nil {
		return Catalog{}, fmt.Errorf("GetArtistCatalog request error: %v", err)
	}

	// try to UnMarshall
	var catalog Catalog
	// defer resp.Body.Close()
	fullBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Catalog{}, fmt.Errorf("catalog response body is nil!!!")
	}
	resp.Body.Close()

	err = json.Unmarshal(fullBody, &catalog)
	if err != nil {
		return Catalog{}, fmt.Errorf("GetArtistCatalog error unmarshalling catalog: %v", err)
	}

	return catalog, nil
}
