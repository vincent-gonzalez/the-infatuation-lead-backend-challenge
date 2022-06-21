package main

import (
	"fmt"
	"os"
	"sort"
)

func main() {
	fmt.Println("Starting application...");

	likeEvents, err := ReceiveEvents("tcp", "localhost", 9090)
	if err != nil {
		// exit here. we can't process further without a list of valid events
		fmt.Printf("Failed while receive events: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Sorting events...");
	sort.Slice(likeEvents, func(i, j int) bool {
		return likeEvents[i].SequenceNum < likeEvents[j].SequenceNum
	});

	fmt.Println("Finding matches...");
	matchSequenceNumbers, err := FindMatchEvents(likeEvents)
	if err != nil {
		fmt.Printf("Failed while finding match events: %v\n", err.Error())
	}

	// fmt.Println("Sorting sequence numbers...");
	// sort.Slice(matchSequenceNumbers, func (i, j int) bool {
	// 	return matchSequenceNumbers[i] < matchSequenceNumbers[j]
	// });

	if len(matchSequenceNumbers) < 1 {
		fmt.Println("No matches found.")
	} else {
		// fmt.Println("Sorting sequence numbers...");
		// sort.Slice(matchSequenceNumbers, func (i, j int) bool {
		// 	return matchSequenceNumbers[i] < matchSequenceNumbers[j]
		// });

		err := SendMatchEvents(matchSequenceNumbers, "tcp", "localhost", 9099)
		if err != nil {
			fmt.Printf("Failed while sending match events: %v", err.Error())
			os.Exit(1)
		}
	}

	fmt.Println("Exiting application.");
}
