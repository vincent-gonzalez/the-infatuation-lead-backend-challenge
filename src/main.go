package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	fmt.Println("Starting application...");

	fmt.Println("Connecting to EVENT SOURCE...")
	eventSourceConnection, err := CreateConnection("tcp", "localhost", 9090)
	if err != nil {
		fmt.Printf("Error while connecting to EVENT SOURCE: %v\n", err.Error())
		fmt.Println("Exiting application.")
		os.Exit(1)
	}
	defer eventSourceConnection.Close()

	likeEvents, err := ReceiveEvents(eventSourceConnection)
	if err != nil {
		// exit here. we can't process further without a list of valid events
		// TODO - we could make a new request to the server to attempt to get
		// the events again.
		fmt.Printf("Failed while receiving events: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Sorting events...");
	// TODO - built in sorting is okay for simple data sets, but we may want
	// to upgrade to a better algorithm if the data becomes larger or more
	// complex.
	sort.Slice(likeEvents, func(i, j int) bool {
		return likeEvents[i].SequenceNum < likeEvents[j].SequenceNum
	});

	fmt.Println("Finding matches...");
	matchSequenceNumbers, err := FindMatchEvents(likeEvents)
	if err != nil {
		fmt.Printf("Failed while finding match events: %v\n", err.Error())
		fmt.Println("Exiting application.")
		os.Exit(1)
	}

	if len(matchSequenceNumbers) < 1 {
		fmt.Println("No matches found.")
	} else {
		fmt.Println("Connecting to EVENT LISTENER...")
		eventListenerConnection, err := CreateConnection("tcp", "localhost", 9099)
		if err != nil {
			fmt.Printf("Error while connecting to EVENT LISTENER: %v\n", err.Error())
			fmt.Println("Exiting application.")
			os.Exit(1)
		}
		defer eventListenerConnection.Close()

		matchStatus, err := SendMatchEvents(eventListenerConnection, matchSequenceNumbers)
		if err != nil {
			fmt.Printf("Failed while sending match events: %v", err.Error())
			fmt.Println("Exiting application.")
			os.Exit(1)
		}

		// Display if all matches received in order or not.
		fmt.Println(matchStatus)
	}

	fmt.Println("Exiting application.");
}
