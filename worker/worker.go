package worker

import (
	"fmt"
	"sync"

	spotify "../spotify"
)

type safeMap struct {
	Map map[string]int
	sync.Mutex
}

var visitedArtists = safeMap{Map: make(map[string]int)}

type Processor struct {
	StartArtistId   string
	CurrentArtistId string
	TargetArtistId  string
	Path            []spotify.Artist
	MaxPathLength   int
	Results         chan<- []spotify.Artist
}

func (p Processor) pathTooLong() bool {
	return len(p.Path) >= p.MaxPathLength
}

type ArtistProcessor struct {
	Processor
}

type AlbumProcessor struct {
	AlbumId string
	Processor
}

func (ap ArtistProcessor) Do() {
	fmt.Printf("Entering ArtistProcessor.Do w/ %v @ %v degrees\n", ap.CurrentArtistId, len(ap.Path))

	currentPathLength := len(ap.Path)

	visitedArtists.Lock()

	// check if we've been here
	shortestPath, ok := visitedArtists.Map[ap.CurrentArtistId]
	if !ok || currentPathLength < shortestPath {
		visitedArtists.Map[ap.CurrentArtistId] = len(ap.Path)
	} else {
		return
	}

	visitedArtists.Unlock()

	// get artist catalog
	catalog := spotify.GetArtistCatalog(ap.CurrentArtistId)

	// get album ids
	for _, album := range catalog.Albums {
		// create album processor
		albumProcessor := AlbumProcessor{
			// CurrentArtistId: ap.CurrentArtistId,
			// TargetArtistId:  ap.TargetArtistId,
			Processor: ap.Processor,
			AlbumId:   album.Id,
		}

		// have at it
		go albumProcessor.Do()
	}

	fmt.Printf("Exiting ArtistProcessor.Do w/ %v @ %v degrees\n", ap.CurrentArtistId, len(ap.Path))
}

func (ap AlbumProcessor) Do() {
	fmt.Printf("Entering AlbumProcessor.Do w/ %v @ %v degrees\n", ap.AlbumId, len(ap.Path))

	// lookup album by id
	album := spotify.GetAlbumById(ap.AlbumId)

	for _, track := range album.Tracks.TrackItems {
		// for artist on track
		for _, artist := range track.Artists {
			if artist.Id == ap.StartArtistId {
				continue
			}

			// update Processor
			ap.Path = append(ap.Path, artist)
			ap.CurrentArtistId = artist.Id

			if artist.Id == ap.TargetArtistId {
				ap.Results <- ap.Path
				return
			} else if !ap.pathTooLong() {
				artistProcessor := ArtistProcessor{
					Processor: ap.Processor,
				}
				go artistProcessor.Do()
			}
		}
	}
	fmt.Printf("Exiting AlbumProcessor.Do w/ %v @ %v degrees\n", ap.AlbumId, len(ap.Path))
}

func ProcessResults(maxPathLength int, results <-chan []spotify.Artist) ([]spotify.Artist, error) {
	bestPath := make([]spotify.Artist, maxPathLength+1)

	for path := range results {
		if len(path) < len(bestPath) {
			bestPath = path
		}
	}

	if len(bestPath) < maxPathLength+1 {
		return bestPath, nil
	} else {
		return nil, fmt.Errorf("Could not find a valid best path!")
	}
}
