package spotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Track struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Number  uint8    `json:"track_number:`
	Artists []Artist `json:"artists"`
}

type TrackList struct {
	TrackItems []Track `json:"items"`
}

type Album struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Artists []Artist  `json:"artists"`
	Tracks  TrackList `json:"tracks"`
}

type Catalog struct {
	Albums []Album `json:"items"`
}

func GetAlbumById(id string) Album {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%v", id)
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	// try to UnMarshall
	var album Album
	fullBody, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(fullBody, &album)

	return album
}
