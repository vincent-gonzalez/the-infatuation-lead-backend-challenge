package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ReceiveEvents(eventSourceConnection net.Conn) ([]LikeEvent, error) {
	var likeEvents []LikeEvent
	// TODO - these magic constants should be defined in the environment configuration
	const dataDelimiter = "|"
	const numberOfEventFields = 4

	scanner := bufio.NewScanner(eventSourceConnection)
	// wait for EVENT SOURCE to send start message
	for scanner.Scan() && scanner.Text() != "EVENT BEGIN" {
		fmt.Println("Waiting for EVENT SOURCE to send events...")
	}
	fmt.Println("EVENT BEGIN")
	fmt.Println("Receiving events...")

	for scanner.Scan() {
		message := scanner.Text()

		// EVENT SOURCE has finished sending messages
		if message == "EVENT END" {
			fmt.Println(message)
			break
		} else {
			newLikeEvent, err := ParseEvent(message, dataDelimiter, numberOfEventFields)
			if err != nil {
				fmt.Println("Failed to parse event message.")
				// in a simple case, dropping the record is fine, but in a real
				// world scenario, this may not be acceptable and should be handled.
				continue
			}

			likeEvents = append(likeEvents, *newLikeEvent)
		}
	}

	return likeEvents, nil
}

func ParseEvent(eventMessage, dataDelimiter string, numberOfFields int) (*LikeEvent, error) {
	messageParts := strings.Split(eventMessage, dataDelimiter)
	if len(messageParts) < numberOfFields {
		return nil, fmt.Errorf("cannot parse fields from message")
	}

	sequenceNum, err := strconv.ParseUint(messageParts[0], 10, 64)
	if err != nil {
		fmt.Println("Input sequence number is not a number: ", err.Error())
		// if sequence number is malformed, we can't parse it correctly.
		// we can't eventually find matches without a sequence number.
		return nil, err
	}

	return &LikeEvent{
		SequenceNum: sequenceNum,
		LikeType:    messageParts[1],
		FromUserId:  messageParts[2],
		ToUserId:    messageParts[3],
	}, nil
}

func FindMatchEvents(likeEvents []LikeEvent) ([]uint64, error) {
	if len(likeEvents) < 1 {
		return nil, fmt.Errorf("likeEvents: slice cannot be empty. len(likeEvents): %d", len(likeEvents))
	}

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

			matchMap[event.FromUserId] = append(matchMap[event.FromUserId], event)
		}
	}

	return matchSequenceNumbers, nil
}

func SendMatchEvents(eventListenerConnection net.Conn, matchSequenceNumbers []uint64) (string, error) {
	if len(matchSequenceNumbers) < 1 {
		return "", fmt.Errorf("matchSequenceNumbers: slice cannot be empty. len(matchSequenceNumbers): %d", len(matchSequenceNumbers))
	}

	scanner := bufio.NewScanner(eventListenerConnection)
	for scanner.Scan() && scanner.Text() != "MATCH BEGIN" {
		fmt.Println("Waiting for EVENT LISTINER to be ready...")
	}
	fmt.Println("MATCH BEGIN")
	fmt.Println("Sending match events...")

	for _, sequenceNumber := range matchSequenceNumbers {
		_, err := eventListenerConnection.Write([]byte(fmt.Sprintf("%d\n", sequenceNumber)))
		if err != nil {
			fmt.Printf("Failed while sending sequence number: %d", sequenceNumber)
			return "", err
		}
	}

	fmt.Println("All match events sent.")

	// Scan for success or failure message from the EVENT LISTENER.
	matchEndMessage := "Unknown"
	for scanner.Scan() {
		matchEndMessage = scanner.Text()
	}

	return matchEndMessage, nil
}

func CreateConnection(protocol string, url string, port uint) (net.Conn, error) {
	newConnection, err := net.Dial(protocol, fmt.Sprintf("%s:%d", url, port))
	if err != nil {
		return nil, err
	}
	return newConnection, err
}
