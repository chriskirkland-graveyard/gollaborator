package main

import (
	"fmt"

	spotify "./spotify"
	worker "./worker"
)

func main() {
	startArtistName := "Banks & Steelz"
	targetArtistName := "Snoop Dogg"
	maxPathLength := 6

	// assume the first result is good... heh.
	startArtist := spotify.GetArtistsByName(startArtistName)[0]
	targetArtist := spotify.GetArtistsByName(targetArtistName)[0]
	fmt.Printf("%+v\n", startArtist)
	fmt.Printf("%+v\n", targetArtist)
	fmt.Println()

	results := make(chan []spotify.Artist, 10)

	// do stuff
	go worker.ArtistProcessor{
		worker.Processor{
			StartArtistId:   startArtist.Id,
			CurrentArtistId: startArtist.Id,
			TargetArtistId:  targetArtist.Id,
			Path:            []spotify.Artist{startArtist},
			MaxPathLength:   maxPathLength,
			Results:         results,
		}}.Do()

	// start processing results
	path, err := worker.ProcessResults(maxPathLength, results)
	if err != nil {
		fmt.Println("YOU FAILED!!!!")
	} else {
		fmt.Printf("GREAT SUCCESS: %v\n", path)
	}
}
