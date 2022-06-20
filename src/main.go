package main

import (
	"bufio"
	"fmt"
	//"os"
	//"io"
	"net"
	"strings"
	"strconv"
	// "time"
	"sort"
)

func main() {
	fmt.Println("Hello World!");
	// fmt.Println("Starting server...");
	// server, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090));
	// if err != nil {
	// 	fmt.Println("Error starting server:", err.Error());
	// 	os.Exit(1);
	// }
	// fmt.Println("Server started.");
	// defer server.Close();
	// fmt.Println("Waiting for client...");
	// isAllEventsReceived := false;
	// for !isAllEventsReceived {
	// 	connection, err := server.Accept();
	// 	if err != nil {
	// 		fmt.Println("Error in client connection:", err.Error());
	// 		os.Exit(1);
	// 	}
	// 	fmt.Println("client connected");
	// 	// do go routine here
	// 	buffer := make([]byte, 1024);
	// 	messageLength, err := connection.Read(buffer);
	// 	if err != nil {
	// 		fmt.Println("Error reading:", err.Error());
	// 	}
	// 	message := string(buffer[:messageLength]);
	// 	fmt.Println("Message received: ", message);
	// 	if message == "EVENT END" {
	// 		isAllEventsReceived = true;
	// 		connection.Close();
	// 	}
	// }
	connection, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", 9090));
	if err != nil {
		fmt.Println("Error connecting to EVENT SOURCE:", err.Error());
	}
	defer connection.Close()
	reader := bufio.NewReader(connection);
	fmt.Println("Waiting for server...");
	isAllEventsReceived := false;
	var likeEvents []LikeEvent;
	likeMatches := make(map[string][]LikeEvent);
	var matchSequenceNumbers []uint64;
	for !isAllEventsReceived {
		//buffer := make([]byte, 1024);
		//messageLength, err := connection.Read(buffer);
		message, err := reader.ReadString('\n');
		if err != nil {
			fmt.Println("Error reading:", err.Error());
		}
		//message := string(buffer[:messageLength]);
		message = message[:len(message)-1]
		//fmt.Println(message);
		//fmt.Println(message == "EVENT BEGIN");
		//time.Sleep(2 * time.Second);
		if message == "EVENT BEGIN" {
			continue;
		} else if message == "EVENT END" {
			isAllEventsReceived = true;
		} else {
			//message = message[:len(message)-1]
			messageParts := strings.Split(message, "|");
			if messageParts[1] == "LIKE_LIKED" {
				sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64);
				if err != nil {
					fmt.Println("Sequence Num not number: ", err.Error());
				}
				newLikeEvent := LikeEvent{
					SequenceNum: sequenceNum,
					LikeType: messageParts[1],
					FromUserId: messageParts[2],
					ToUserId: messageParts[3],
				};
				likeEvents = append(likeEvents, newLikeEvent);
				if events, found := likeMatches[newLikeEvent.ToUserId]; found {
					var matchingEvent LikeEvent;
					var matchSequenceNum uint64;
					for _, m := range events {
						if m.ToUserId == newLikeEvent.FromUserId {
							matchingEvent = m;
							break;
						}
					}
					if matchingEvent.SequenceNum > newLikeEvent.SequenceNum {
						matchSequenceNum = matchingEvent.SequenceNum;
					} else {
						matchSequenceNum = newLikeEvent.SequenceNum;
					}

					matchSequenceNumbers = append(matchSequenceNumbers, matchSequenceNum);
				} else {
					likeMatches[newLikeEvent.FromUserId] = append(likeMatches[newLikeEvent.FromUserId], newLikeEvent);
				}
			}
		}
		// message = message[:len(message)-1]
		// messageParts := strings.Split(message, "|");
		// sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64);
		// if err != nil {
		// 	fmt.Println("Sequence Num not number: ", err.Error());
		// }
		// likeEvents = append(likeEvents, LikeEvent{
		// 	SequenceNum: sequenceNum,
		// 	LikeType: messageParts[1],
		// 	FromUserId: messageParts[2],
		// 	ToUserId: messageParts[3],
		// });
		//fmt.Println("Message received: ", message);
		// if message == "EVENT END" {
		// 	isAllEventsReceived = true;
		// 	//connection.Close();
		// 	//break;
		// }
	}
	//connection.Close();
	fmt.Println("All messages received.");
	fmt.Println("Connection closed.");
	fmt.Println("Server shutdown.");
	fmt.Println("Printing received messages...");
	sort.Slice(likeEvents, func(i, j int) bool {
		return likeEvents[i].SequenceNum < likeEvents[j].SequenceNum;
	});
	sort.Slice(matchSequenceNumbers, func (i, j int) bool {
		return matchSequenceNumbers[i] < matchSequenceNumbers[j];
	});
	for _, event := range likeEvents {
		fmt.Printf("Sequence Num: %d Like Type: %s From User: %s To User: %s\n", event.SequenceNum, event.LikeType, event.FromUserId, event.ToUserId);
	}
	for _, s := range matchSequenceNumbers {
		fmt.Printf("Match num: %d\n", s);
	}
}
