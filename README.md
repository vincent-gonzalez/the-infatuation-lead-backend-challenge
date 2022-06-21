# Lead Backend Challenge Solution - Vincent Gonzalez

## How to run
A makefile has been included in the base directory. Use `make buildRun` to build and execute the program. Use `make run` to simply run the program after it has already been compiled.

As a courtesy, a compiled binary has been provided as a part of this repo. So, you may use the `make run` command to immediately execute the solution program.

## Design
The basic algorithm of the program flows as follows:
1. Connect to the EVENT SOURCE server over TCP on localhost:9090.
2. Receive the event messages that the EVENT SOURCE sends.
    a. Wait for the EVENT BEGIN message.
    b. Scan the incoming data stream for a message and parse it into a struct (LikeEvent) to contain the event data fields.
    c. Append the message to a slice ([]LikeEvent).
    d. Loop until the EVENT END message is received.
3. Sort the the received events into sequential order.
4. Find matching like events.
    a. A map of map\[string\][]LikeEvent is used to keep track of all LIKE_LIKED messages.
    b. Iterate over the slice of all event messages.
    c. If the like type of the message is not "LIKE_LIKED", then skip over the message and move on to the next one.
    d. If the like type of the message is "LIKE_LIKED", then check the map of LIKE_LIKED events for a item matching the current message's ToUserId.
    e. If a matching user ID is found, get that entry's slice of LikeEvents from the map. Iterate through that list to see if that user has also liked the current message's user.
    f. If a match is found, take the later sequence number and add it to a slice of match event sequence numbers.
    g. Add the current LIKE_LIKED message to the map using the message's FromUserId as the key.
    h. Return all match event sequence numbers as a slice.
5. Send the match event sequence numbers to the EVENT LISTENER.
    a. Connect to the EVENT LISTENER over TCP and using localhost:9099.
    b. Wait for the EVENT LISTENER to send the MATCH BEGIN message.
    c. Iterate over the slice of sequence numbers and send the sequence numbers one by one.
    d. Wait and listen for the EVENT LISTENER to send either the MATCH END - OK or MATCH END - ERROR message.
6. Exit the application.

## Files
### main.go
The main program logic originates here.

### EventFunctions.go
o

### EventFunctions_test.go
o

### LikeEvent.go
o

### Makefile
o
