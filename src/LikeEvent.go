package main

type LikeEvent struct {
	// uint64 is used since it is possible for the system to generate thousands of
	// like events
	SequenceNum uint64 // a negative sequence number does not make sense in the application
	LikeType string
	// TRADE-OFF: using strings allows user IDs to be something other than numbers
	// such as a GUID, hash value, email address, or user handle. However, based on
	// the example input, the natural type of the ID is a positive integer.
	// If mathematical operations are to be supported on the user IDs, then
	// an unsigned integer type makes sense. From a security standpoint,
	// integer values may allow a type of user enumeration to take place
	// should an attacker gain access to the values. A pattern such as
	// when a user joined may emerge as, for example, older users would
	// be assigned a user ID number that is closer to zero. This could
	// have the unwanted situation of identifying administrator users as
	// they have a higher chance of entering the system at an earlier date
	// than other users.
	FromUserId string
	ToUserId string
}
