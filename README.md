# Lead Backend Challenge Solution - Vincent Gonzalez

## Runtime Environment
The minimum version of Go that is required to run this solution is `1.15`. This version is also defined in the `go.mod` file of this solution. However, I have compiled and run this solution using version `1.17.11` of Go.

## How to run
A makefile has been included in the base directory. Use `make buildRun` to build and execute the program. Use `make run` to simply run the program after it has already been compiled.

As a courtesy, a compiled binary file of this solution has been provided as a part of this repo. So, you may use the `make run` command to immediately execute the solution program. If the binary file gives you trouble, feel free to build the program from this repo using the `make build` command.

## Design
The basic algorithm of the program flows as follows:
1. Connect to the EVENT SOURCE server over TCP on localhost:9090.
2. Receive the event messages that the EVENT SOURCE sends.
    1. Wait for the EVENT BEGIN message.
    2. Scan the incoming data stream for a message and parse it into a struct (LikeEvent) to contain the event data fields.
    3. Append the message to a slice ([]LikeEvent).
    4. Loop until the EVENT END message is received.
3. Sort the the received events into sequential order.
4. Find matching like events.
    1. A map of map\[string\][]LikeEvent is used to keep track of all LIKE_LIKED messages.
    2. Iterate over the slice of all event messages.
    3. If the like type of the message is not "LIKE_LIKED", then skip over the message and move on to the next one.
    4. If the like type of the message is "LIKE_LIKED", then check the map of LIKE_LIKED events for a item matching the current message's ToUserId.
    5. If a matching user ID is found, get that entry's slice of LikeEvents from the map. Iterate through that list to see if that user has also liked the current message's user.
    6. If a match is found, take the later sequence number and add it to a slice of match event sequence numbers.
    7. Add the current LIKE_LIKED message to the map using the message's FromUserId as the key.
    8. Return all match event sequence numbers as a slice.
5. Send the match event sequence numbers to the EVENT LISTENER.
    1. Connect to the EVENT LISTENER over TCP and using localhost:9099.
    2. Wait for the EVENT LISTENER to send the MATCH BEGIN message.
    3. Iterate over the slice of sequence numbers and send the sequence numbers one by one.
    4. Wait and listen for the EVENT LISTENER to send either the MATCH END - OK or MATCH END - ERROR message.
6. Exit the application.

## Thoughts on improvements in future iterations
### Enumerated types
A helpful improvement would be to use enumeration types instead of hard-coded strings for things like *EVENT BEGIN* and *LIKE_LIKED*. Go's support of enumerations is different from a language like C#, and is less helpful when it comes to string enumerations as the type checking is not as strict with string enumerations in Go. Constants could be a middle ground in the meantime.
### Environment variables/configuration file
Several items in this program would do well when defined by the running environment through environment variables or configuration files. For example, the current message field delimiter is the pipe (|) character. This could change in the future, and instead of changing the code, we could simply get the new message delimiter from an environment variable and continue processing as normal.

## Files
### main.go
The main program logic originates here. The

### EventFunctions.go
Contains functions that operate on each event.
#### ReceiveEvents(protocol string, eventSourceURL string, port uint) ([]LikeEvent, error)
Connects to the EVENT SOURCE server, reads the like event messages that it sends, and parses each message into a LikeEvent object.

Parameters:
- protocol `string` - the network protocol to connect to the EVENT SOURCE.
- eventSourceURL `string` - the server address of the EVENT SOURCE.
- port `uint` - the port that the EVENT SOURCE is listening on. Negative integers are not valid port values, so an unsigned integer is used.

Return values:
- `[]LikeEvent` - a slice of all parsed messages received from EVENT SOURCE
- `error` - `nil` or any failure to receive messages or failure to parse messages.

Logic:
- The function opens a (TCP) connection to the EVENT SOURCE.
- After opening a connection, a scanner to read incoming messages is created using the open connection.
- The function waits for the EVENT SOURCE to send the EVENT BEGIN message.
- The function then reads every message that is received over the open connection until the EVENT END message is received.
- For each event message that is received, the message is parsed into a LikeEvent object using the ParseEvent method.
- After creating a LikeEvent from a message, the LikeEvent is stored in a slice of LikeEvents.
- Once the function is finished reading all the messages sent from the EVENT SOURCE, the connection to the EVENT SOURCE is closed, and the slice of LikeEvents is returned to the caller.

#### ParseEvent(eventMessage, dataDelimiter string, numberOfFields int) (*LikeEvent, error)
Takes in a event message string and parses it into a LikeEvent object using the dataDelimiter to determine each field location.

Parameters:
- eventMessage `string` - the event message to be parsed
- dataDelimiter `string` - the delimiter character that separates each data field in the event message.
- numberOfFields `int` - the number of data fields that the event message is expected to contain.

Return values:
- `*LikeEvent` - a pointer to a newly created and parsed LikeEvent
- `error` - `nil` when a message is successfully parsed. Otherwise a failure due to receiving an invalid message string or a failure to convert the sequence number to a `uint64`.

Logic:
- The data contained in the eventMessage parameter is split into multiple fields using the delimiter character found in the dataDelimiter parameter.
- If the number of fields obtained in the split is less than the number of expected fields from the numberOfFields parameter, then the message is invalid and the function returns a nil and error.
- If the message is valid, then the function continues to parse out the data into a LikeEvent.
- The sequence number field is converted from a string to a uint64 value. If this field cannot be converted, then the message is invalid (as the sequence number is not a number) and a nil and error is returned to the caller.
- Each message field is assigned to the properties of a new LikeEvent object which is returned along with a nil error (no error occurred).

#### FindMatchEvents(likeEvents []LikeEvent) ([]uint64, error)
This function iterates over a sorted slice of LikeEvents to find match events and returns a slice of match event sequence numbers.

Parameters:
- likeEvents `[]LikeEvent` - a slice of sorted LikeEvent objects

Return values:
- `[]uint64` - a slice that contains all the match events that were found in the likeEvents slice.
- `error` - If the passed in likeEvents slice is empty, then an error is returned, otherwise `nil` is returned.

Logic:

#### SendMatchEvents(matchSequenceNumbers []uint64, protocol string, eventListenerURL string, port uint) error

### EventFunctions_test.go
Contains unit tests for the event functions.

#### TestParseEvent(t *testing.T)
Tests that a valid message is parsed correctly.

#### TestParseEventBadEvent (t *testing.T)
Tests that a malformed message is not parsed.

#### TestFindMatchEvents(t *testing.T)
Tests that matches are found with known good input that contains match events.

#### TestNoMatchEventsFound(t *testing.T)
Tests that match events are not found when the input is known to not have match events.

#### testSlicesAreEqual(a, b []uint64) bool
A helper function for testing the equality of uint64 slices.

### LikeEvent.go
Contains the data type that represents the like events sent from the EVENT SOURCE.
There are 4 fields:
- SequenceNum `uint64` - contains the sequence number of the message. Values from the message will likely need to be converted from `string` to `uint64`.
- LikeType `string` - contains the like type of a message
- FromUserId `string` - contains the user ID of the user sending the like
- ToUserId `string` - contains the user ID of the user that is receiving the like.

#### Thoughts behind the SequenceNum typing
Negative numbers don't make sense on an identifier field like SequenceNum. It is also possible that tens of thousands of event messages can be generated, so a large variable type is required to hold larger sequence number values. In addition, while the sequence number could be typed as a `string`, the natural type of a sequence number is an integer.

#### Thoughts behind the FromUserId and ToUserId typing
There is a trade-off to using strings for the user IDs instead of integers. Strings allow user IDs to be something other than numbers
such as a GUID, hash value, email address, or user handle. However, based on
the example input, the natural type of the ID is a positive integer.
If mathematical operations are to be supported on the user IDs, then
an unsigned integer type makes sense. From a security standpoint,
integer values may allow a type of user enumeration to take place
should an attacker gain access to the values. A pattern such as
when a user joined may emerge as, for example, older users would
be assigned a user ID number that is closer to zero. This could
have the unwanted situation of identifying administrator users as
they have a higher chance of entering the system at an earlier date
than other users.

### Makefile
Contains the CLI commands for the project. It also contains some environment variables that are used within the CLI commands.

The file includes the following commands:
- `build` - changes to the source directory, compiles the program, and places the executable into a build directory.
- `buildRun` - combines the `build` and `run` commands into one.
- `run` - executes the compiled program.
- `test` - executes the unit tests in the project.
