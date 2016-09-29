package main

import (
	"fmt"
	"log"
	"time"

	http "net/http"
	_ "net/http/pprof"

	spotify "./spotify"
	utils "./utils"
	worker "./worker"
)

func main() {
	startArtistName := "Banks & Steelz"
	targetArtistName := "Snoop Dogg"
	maxPathLength := 4

	// var startArtistName = flag.String("starting artist", "Banks & Steelz", "starting artist")
	// var targetArtistName = flag.String("target artist", "Snoop Dogg", "target artist")

	// assume the first result is good... heh.
	startArtistPossibilities, err := spotify.GetArtistsByName(startArtistName)
	if err != nil {
		panic(err)
	}
	startArtist := startArtistPossibilities[0]

	targetArtistPossibilities, err := spotify.GetArtistsByName(targetArtistName)
	if err != nil {
		panic(err)
	}
	targetArtist := targetArtistPossibilities[0]
	// targetArtist := spotify.Artist{Id: "6NYHiotMjdo6XqeVcCphtf", Name: "Space"}

	fmt.Printf("%+v\n", startArtist)
	fmt.Printf("%+v\n", targetArtist)
	fmt.Println()

	results := make(chan []spotify.Artist, 10)
	printQueue := make(chan string, 100000)

	// pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:1234", nil))
	}()

	// start printer
	go utils.Printer(printQueue)

	// do stuff
	go worker.ArtistProcessor{
		worker.Processor{
			StartArtistId:   startArtist.Id,
			CurrentArtistId: startArtist.Id,
			TargetArtistId:  targetArtist.Id,
			Path:            []spotify.Artist{startArtist},
			MaxPathLength:   maxPathLength,
			Results:         results,
			Printer:         printQueue,
		}}.Do()

	<-time.After(1000 * time.Millisecond)

	// close channels
	go utils.WaitAndClose(worker.NumActiveProcessors.WaitGroup, results, printQueue)

	// start processing results
	path, err := worker.ProcessResults(maxPathLength, results, printQueue)
	if err != nil {
		fmt.Println("YOU FAILED!!!!")
	} else {
		fmt.Printf("GREAT SUCCESS: %v\n", path)
	}
}
