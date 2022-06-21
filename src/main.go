package main

import (
	// "bufio"
	"fmt"
	// "net"
	// "strings"
	// "strconv"
	"sort"
)

func main() {
	fmt.Println("Starting application...");

	// fmt.Println("Connecting to EVENT SOURCE...");
	// connection, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", 9090));
	// if err != nil {
	// 	fmt.Println("Error connecting to EVENT SOURCE:", err.Error());
	// }
	// defer connection.Close();

	// fmt.Println("Waiting for EVENT SOURCE...");

	// reader := bufio.NewReader(connection);
	// isAllEventsReceived := false;
	// var likeEvents []LikeEvent;
	// likeMatches := make(map[string][]LikeEvent);
	// var matchSequenceNumbers []uint64;

	// for !isAllEventsReceived {
	// 	message, err := reader.ReadString('\n');
	// 	if err != nil {
	// 		fmt.Println("Error reading event:", err.Error());
	// 	}
	// 	message = message[:len(message)-1]

	// 	if message == "EVENT BEGIN" {
	// 		fmt.Println("START receiving events...");
	// 		continue;
	// 	} else if message == "EVENT END" {
	// 		fmt.Println("END receiving events...");
	// 		isAllEventsReceived = true;
	// 	} else {
	// 		messageParts := strings.Split(message, "|");
	// 		// sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64);
	// 		// 	if err != nil {
	// 		// 		fmt.Println("Sequence Num not number: ", err.Error());
	// 		// 	}
	// 		// 	newLikeEvent := LikeEvent{
	// 		// 		SequenceNum: sequenceNum,
	// 		// 		LikeType: messageParts[1],
	// 		// 		FromUserId: messageParts[2],
	// 		// 		ToUserId: messageParts[3],
	// 		// 	};
	// 		// 	likeEvents = append(likeEvents, newLikeEvent);
	// 		if messageParts[1] == "LIKE_LIKED" {
	// 			sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64);
	// 			if err != nil {
	// 				fmt.Println("Sequence Num not number: ", err.Error());
	// 			}

	// 			newLikeEvent := LikeEvent{
	// 				SequenceNum: sequenceNum,
	// 				LikeType: messageParts[1],
	// 				FromUserId: messageParts[2],
	// 				ToUserId: messageParts[3],
	// 			};
	// 			likeEvents = append(likeEvents, newLikeEvent);

	// 			if events, found := likeMatches[newLikeEvent.ToUserId]; found {
	// 				var matchingEvent LikeEvent;
	// 				var matchSequenceNum uint64;
	// 				foundMatchingEvent := false;

	// 				for _, m := range events {
	// 					if m.ToUserId == newLikeEvent.FromUserId {
	// 						matchingEvent = m;
	// 						foundMatchingEvent = true;
	// 						break;
	// 					}
	// 				}

	// 				if foundMatchingEvent {
	// 					if matchingEvent.SequenceNum > newLikeEvent.SequenceNum {
	// 						matchSequenceNum = matchingEvent.SequenceNum;
	// 					} else {
	// 						matchSequenceNum = newLikeEvent.SequenceNum;
	// 					}

	// 					matchSequenceNumbers = append(matchSequenceNumbers, matchSequenceNum);
	// 				}
	// 			} else {
	// 				likeMatches[newLikeEvent.FromUserId] = append(likeMatches[newLikeEvent.FromUserId], newLikeEvent);
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Println("START receiving events...")
	likeEvents, _ := ReceiveEvents("tcp", "localhost", 9090)

	fmt.Println("Sorting events...");
	sort.Slice(likeEvents, func(i, j int) bool {
		return likeEvents[i].SequenceNum < likeEvents[j].SequenceNum;
	});
	fmt.Println("Finding matches...");
	matchSequenceNumbers, _ := FindMatchEvents(likeEvents)

	// fmt.Println("All messages received.");
	// fmt.Println("Connection to EVENT SOURCE closed.");
	// fmt.Println("Server shutdown.");
	fmt.Println("Printing received messages...");

	sort.Slice(matchSequenceNumbers, func (i, j int) bool {
		return matchSequenceNumbers[i] < matchSequenceNumbers[j];
	});
	// for _, event := range likeEvents {
	// 	fmt.Printf("Sequence Num: %d Like Type: %s From User: %s To User: %s\n", event.SequenceNum, event.LikeType, event.FromUserId, event.ToUserId);
	// }

	if len(matchSequenceNumbers) < 1 {
		fmt.Println("No matches found.");
		fmt.Println("Exiting application.");
	} else {
		// for _, s := range matchSequenceNumbers {
		// 	fmt.Printf("Match num: %d\n", s);
		// }
		// for _, s := range matchInfo {
		// 	// fmt.Printf("Match: %s\n", s);
		// }
		fmt.Println("Connecting to EVENT LISTENER...")
		err := SendMatchEvents(matchSequenceNumbers, "tcp", "localhost", 9099)
		if err != nil {
			fmt.Println(err.Error());
		}
	}

}
