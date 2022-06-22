# Lead Backend Challenge Solution - Vincent Gonzalez

## How to run
A makefile has been included in the base directory. Use `make buildRun` to build and execute the program. Use `make run` to simply run the program after it has already been compiled.

As a courtesy, a compiled binary has been provided as a part of this repo. So, you may use the `make run` command to immediately execute the solution program. If the binary file gives you trouble, feel free to build the program from this repo using the `make build` command.

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
### Concurrency
## Files
### main.go
The main program logic originates here. The

### EventFunctions.go
Contains functions that operate on each event.

### EventFunctions_test.go
Contains unit tests for the event functions.

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
