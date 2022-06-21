package main

import (
	// "fmt"
	"testing"
)

func TestParseEvent(t *testing.T) {
	testEvent := "5|LIKE_LIKED|300|100"

	got, _ := ParseEvent(testEvent, "|", 4)
	want := LikeEvent{
		SequenceNum: 5,
		LikeType: "LIKE_LIKED",
		FromUserId: "300",
		ToUserId: "100",
	}

	if *got != want {
		t.Errorf("Expected %v but got %v", want, got)
	}
}

func TestFindMatchEvents(t *testing.T) {
	testMatches := []LikeEvent {
		{
			SequenceNum: 1,
			LikeType: "LIKE_UNSPECIFIED",
			FromUserId: "300",
			ToUserId: "100",
		},
		{
			SequenceNum: 2,
			LikeType: "LIKE_LIKED",
			FromUserId: "100",
			ToUserId: "200",
		},
		{
			SequenceNum: 3,
			LikeType: "LIKE_NOT_LIKED",
			FromUserId: "300",
			ToUserId: "200",
		},
		{
			SequenceNum: 4,
			LikeType: "LIKE_LIKED",
			FromUserId: "200",
			ToUserId: "100",
		},
		{
			SequenceNum: 5,
			LikeType: "LIKE_LIKED",
			FromUserId: "300",
			ToUserId: "100",
		},
		{
			SequenceNum: 6,
			LikeType: "LIKE_UNSPECIFIED",
			FromUserId: "200",
			ToUserId: "300",
		},
		{
			SequenceNum: 7,
			LikeType: "LIKE_LIKED",
			FromUserId: "100",
			ToUserId: "300",
		},
	}

	got, _ := FindMatchEvents(testMatches);
	want := []uint64 {4, 7};

	if !testSlicesAreEqual(want, got) {
		t.Errorf("Expected %v, but got %v\n", want, got);
	}
	//  else {
	// 	fmt.Printf("Expected %v and got %v\n", want, got);
	// }
}

func TestNoMatchEventsFound(t *testing.T) {
	matches := []LikeEvent {
		{
			SequenceNum: 1,
			LikeType: "LIKE_UNSPECIFIED",
			FromUserId: "300",
			ToUserId: "100",
		},
		{
			SequenceNum: 2,
			LikeType: "LIKE_NOT_LIKED",
			FromUserId: "100",
			ToUserId: "200",
		},
		{
			SequenceNum: 3,
			LikeType: "LIKE_NOT_LIKED",
			FromUserId: "300",
			ToUserId: "200",
		},
		{
			SequenceNum: 4,
			LikeType: "LIKE_LIKED",
			FromUserId: "200",
			ToUserId: "100",
		},
		{
			SequenceNum: 5,
			LikeType: "LIKE_NOT_LIKED",
			FromUserId: "300",
			ToUserId: "100",
		},
		{
			SequenceNum: 6,
			LikeType: "LIKE_UNSPECIFIED",
			FromUserId: "200",
			ToUserId: "300",
		},
		{
			SequenceNum: 7,
			LikeType: "LIKE_LIKED",
			FromUserId: "100",
			ToUserId: "300",
		},
	}

	got, _ := FindMatchEvents(matches)
	// want := []uint{}

	if len(got) != 0 {
		t.Errorf("Expected length of matches to be zero. Got: %d", len(got))
	}
}

func testSlicesAreEqual(a, b []uint64) bool {
	if len(a) != len(b) {
		return false;
	}

	for i, v := range a {
		if v != b[i] {
			return false;
		}
	}

	return true;
}
