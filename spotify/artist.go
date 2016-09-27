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

func GetArtistsByName(name string) []Artist {
	client := &http.Client{}

	// url parameters
	params := url.Values{}
	params.Add("q", name)
	params.Add("type", "artist")

	url := fmt.Sprintf("https://api.spotify.com/v1/search?=%v", params.Encode())
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	// try to UnMarshall
	var ar ArtistResponse
	defer resp.Body.Close()
	fullBody, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(fullBody, &ar)

	return ar.Artists.ArtistItem
}

func GetArtistById(id string) Artist {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v", id)
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	// try to UnMarshall
	var artist Artist
	defer resp.Body.Close()
	fullBody, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(fullBody, &artist)

	return artist
}

func GetArtistCatalog(id string) Catalog {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%v/albums?market=US", id)
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	// try to UnMarshall
	var catalog Catalog
	// defer resp.Body.Close()
	fullBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("catalog response body is nil!!!")
	}
	resp.Body.Close()
	_ = json.Unmarshal(fullBody, &catalog)

	return catalog

}
