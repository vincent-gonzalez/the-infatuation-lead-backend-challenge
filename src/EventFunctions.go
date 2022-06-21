package main

import (
	"bufio"
	"fmt"
	// "math/rand"
	"net"
	"strconv"
	"strings"
	// "time"
)

func ReceiveEvents(protocol string, eventSourceURL string, port uint) ([]LikeEvent, error) {
	// Put param checking here
	eventSourceConnection, err := net.Dial(protocol, fmt.Sprintf("%s:%d", eventSourceURL, port));
	if err != nil {
		fmt.Println("Error connecting to EVENT SOURCE:", err.Error())
		return nil, err
	}
	defer eventSourceConnection.Close()

	events, err := GetLikeEvents(eventSourceConnection);
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return events, nil
}

func GetLikeEvents(eventSourceConnection net.Conn) ([]LikeEvent, error) {
	reader := bufio.NewReader(eventSourceConnection)
	isAllEventsReceived := false
	var likeEvents []LikeEvent

	for !isAllEventsReceived {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading event:", err.Error())
		}
		message = message[:len(message)-1]

		if message == "EVENT BEGIN" {
			fmt.Println("START receiving events...")
			continue
		} else if message == "EVENT END" {
			fmt.Println("END receiving events...")
			isAllEventsReceived = true
		} else {
			messageParts := strings.Split(message, "|")

			// if messageParts[1] == "LIKE_LIKED" {
			// 	sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64)
			// 	if err != nil {
			// 		fmt.Println("Sequence Num not number: ", err.Error())
			// 	}

			// 	newLikeEvent := LikeEvent{
			// 		SequenceNum: sequenceNum,
			// 		LikeType: messageParts[1],
			// 		FromUserId: messageParts[2],
			// 		ToUserId: messageParts[3],
			// 	}
			// 	likeEvents = append(likeEvents, newLikeEvent)
			// }
			sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64)
			if err != nil {
				fmt.Println("Sequence Num not number: ", err.Error())
			}

			newLikeEvent := LikeEvent{
				SequenceNum: sequenceNum,
				LikeType: messageParts[1],
				FromUserId: messageParts[2],
				ToUserId: messageParts[3],
			}
			likeEvents = append(likeEvents, newLikeEvent)
		}
	}

	return likeEvents, nil
}

func FindMatchEvents(likeEvents []LikeEvent) ([]uint64, error) {
	matchMap := make(map[string][]LikeEvent)
	var matchSequenceNumbers []uint64

	for _, event := range likeEvents {
		if event.LikeType == "LIKE_LIKED" {
			if eventList, found := matchMap[event.ToUserId]; found {
				var matchingEvent LikeEvent
				var matchSequenceNum uint64
				foundMatchingEvent := false

				for _, m := range eventList {
					if m.ToUserId == event.FromUserId {
						matchingEvent = m
						foundMatchingEvent = true
						break
					}
				}

				if foundMatchingEvent {
					if matchingEvent.SequenceNum > event.SequenceNum {
						matchSequenceNum = matchingEvent.SequenceNum
					} else {
						matchSequenceNum = event.SequenceNum
					}

					matchSequenceNumbers = append(matchSequenceNumbers, matchSequenceNum)
				}
			}
			// else {
			// 	matchMap[event.FromUserId] = append(matchMap[event.FromUserId], event)
			// }
			matchMap[event.FromUserId] = append(matchMap[event.FromUserId], event)
		}
	}
	matchSequenceNumbers = filterDups(matchSequenceNumbers)
	return matchSequenceNumbers, nil
}

func filterDups(input []uint64) []uint64 {
	inputMap := make(map[uint64]bool)
	var results []uint64
	for _, i := range input {
		if _, found := inputMap[i]; !found {
			inputMap[i] = true
			results = append(results, i)
		}
	}
	return results
}

func SendMatchEvents(matchSequenceNumbers []uint64, protocol string, eventListenerURL string, port uint) error {
	// Put param checking here
	eventListenerConnection, err := net.Dial(protocol, fmt.Sprintf("%s:%d", eventListenerURL, port));
	if err != nil {
		fmt.Println("Error connecting to EVENT LISTENER:", err.Error())
		return err
	}
	defer eventListenerConnection.Close()

	scanner := bufio.NewScanner(eventListenerConnection)

	for scanner.Scan() && scanner.Text() != "MATCH BEGIN" {
		fmt.Println("Waiting for EVENT LISTINER to be ready...")
	}

	for _, sequenceNumber := range matchSequenceNumbers {
		fmt.Printf("Sending: %d\n", sequenceNumber);
		eventListenerConnection.Write([]byte(fmt.Sprintf("%d\n", sequenceNumber)))
	}

	matchEndMessage := "Unknown"
	for scanner.Scan() {
		matchEndMessage = scanner.Text()
	}
	fmt.Println(matchEndMessage)
	return nil
}
