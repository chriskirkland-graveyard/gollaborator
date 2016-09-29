package worker

import (
	"fmt"
	"runtime"
	"time"

	spotify "../spotify"
	utils "../utils"
)

var visitedArtists = utils.SafeMap{Map: make(map[string]int)}
var NumActiveProcessors = utils.SafeWaitGroup{}

type Processor struct {
	StartArtistId   string
	CurrentArtistId string
	TargetArtistId  string
	Path            []spotify.Artist
	MaxPathLength   int
	Results         chan<- []spotify.Artist
	Printer         chan<- string
}

func (p Processor) pathTooLong() bool {
	return len(p.Path) >= p.MaxPathLength
}

func (p Processor) Print(msg string) {
	p.Printer <- msg
}

type ArtistProcessor struct {
	Processor
}

type AlbumProcessor struct {
	AlbumId string
	Processor
}

func (ap ArtistProcessor) Do() {
	NumActiveProcessors.Add(1)
	defer NumActiveProcessors.Done()

	// fmt.Printf("Entering ArtistProcessor.Do w/ %v @ %v degrees (%v)\n", ap.CurrentArtistId, len(ap.Path), ap.Path)
	// defer fmt.Printf("Exiting ArtistProcessor.Do w/ %v @ %v degrees (%v)\n", ap.CurrentArtistId, len(ap.Path), ap.Path)

	ap.Print(fmt.Sprintf("Entering ArtistProcessor.Do w/ %v @ %v degrees (%v)\n", ap.CurrentArtistId, len(ap.Path), ap.Path))
	defer ap.Print(fmt.Sprintf("Exiting ArtistProcessor.Do w/ %v @ %v degrees (%v)\n", ap.CurrentArtistId, len(ap.Path), ap.Path))

	currentPathLength := len(ap.Path)

	visitedArtists.Lock()
	// fmt.Printf("Visiting artist LOCK %v\n", ap.CurrentArtistId)
	ap.Print(fmt.Sprintf("Visiting artist LOCK %v\n", ap.CurrentArtistId))

	// check if we've been here
	shortestPath, ok := visitedArtists.Map[ap.CurrentArtistId]
	if !ok {
		// fmt.Printf("Visiting artist for the first time: %v\n", ap.CurrentArtistId)
		ap.Print(fmt.Sprintf("Visiting artist for the first time: %v\n", ap.CurrentArtistId))
		visitedArtists.Map[ap.CurrentArtistId] = currentPathLength
	} else if currentPathLength < shortestPath {
		// fmt.Printf("Found shorter path for artist: %v\n", ap.CurrentArtistId)
		ap.Print(fmt.Sprintf("Found shorter path for artist: %v\n", ap.CurrentArtistId))
		visitedArtists.Map[ap.CurrentArtistId] = currentPathLength
	} else {
		// fmt.Printf("On a longer path for artist \"%v\". Exiting...\n", ap.CurrentArtistId)
		ap.Print(fmt.Sprintf("On a longer path for artist \"%v\". Exiting...\n", ap.CurrentArtistId))
		visitedArtists.Unlock()
		return
	}

	// fmt.Printf("Visiting artist UNLOCK %v\n", ap.CurrentArtistId)
	ap.Print(fmt.Sprintf("Visiting artist UNLOCK %v\n", ap.CurrentArtistId))
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
}

func (ap AlbumProcessor) Do() {
	NumActiveProcessors.Add(1)
	defer NumActiveProcessors.Done()

	// fmt.Printf("Entering AlbumProcessor.Do w/ %v @ %v degrees (%v)\n", ap.AlbumId, len(ap.Path), ap.Path)
	// defer fmt.Printf("Exiting AlbumProcessor.Do w/ %v @ %v degrees (%v)\n", ap.AlbumId, len(ap.Path), ap.Path)
	ap.Print(fmt.Sprintf("Entering AlbumProcessor.Do w/ %v @ %v degrees (%v)\n", ap.AlbumId, len(ap.Path), ap.Path))
	defer ap.Print(fmt.Sprintf("Exiting AlbumProcessor.Do w/ %v @ %v degrees (%v)\n", ap.AlbumId, len(ap.Path), ap.Path))

	// lookup album by id
	album := spotify.GetAlbumById(ap.AlbumId)

	for _, track := range album.Tracks.TrackItems {
		// for artist on track
		for _, artist := range track.Artists {
			if artist.Id == ap.StartArtistId {
				continue
			}

			newPath := append(ap.Path, artist)
			if artist.Id == ap.TargetArtistId {
				ap.Results <- newPath

			} else if !ap.pathTooLong() {
				artistProcessor := ArtistProcessor{
					Processor: Processor{
						StartArtistId:   ap.StartArtistId,
						CurrentArtistId: artist.Id,
						TargetArtistId:  ap.TargetArtistId,
						Path:            newPath,
						MaxPathLength:   ap.MaxPathLength,
						Results:         ap.Results,
						Printer:         ap.Printer,
					},
				}
				go artistProcessor.Do()
			}
		}
	}
}

func ProcessResults(maxPathLength int, results <-chan []spotify.Artist, printQueue <-chan string) ([]spotify.Artist, error) {
	bestPath := make([]spotify.Artist, maxPathLength+1)
	ticker := time.NewTicker(time.Millisecond * 500)

	breakOut := false

	for {
		select {
		case path, ok := <-results:
			if !ok {
				// results channel closed
				fmt.Println("Results channel was closed!!!")
				breakOut = true
			} else if len(path) < len(bestPath) {
				fmt.Printf("Found a better path!!! w/ %v, old path is: %v\n", path, bestPath)
				bestPath = path
			}
		case <-ticker.C:
			fmt.Printf("NUMBER OF GOROUTINES: %v\n", runtime.NumGoroutine())
			fmt.Printf("LENGTH OF PRINTQUEUE: %v\n", len(printQueue))
			fmt.Println()
		}

		if breakOut {
			break
		}

	}

	if len(bestPath) < maxPathLength+1 {
		return bestPath, nil
	} else {
		return nil, fmt.Errorf("Could not find a valid best path!")
	}
}
