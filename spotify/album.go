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

func GetAlbumById(id string) (Album, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://api.spotify.com/v1/albums/%v", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Album{}, fmt.Errorf("GetAlbumById error creating request: %v", err)
	}

	<-GlobalRateLimiter
	resp, err := client.Do(req)
	if err != nil {
		return Album{}, fmt.Errorf("GetAlbumById request error: %v", err)
	}

	// try to UnMarshall
	var album Album
	fullBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(fullBody, &album)
	if err != nil {
		return Album{}, fmt.Errorf("GetAlbumById error unmarshalling album: %v", err)
	}

	return album, nil
}
