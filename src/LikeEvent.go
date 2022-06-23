package main

type LikeEvent struct {
	SequenceNum uint64
	LikeType    string
	FromUserId  string
	ToUserId    string
}
