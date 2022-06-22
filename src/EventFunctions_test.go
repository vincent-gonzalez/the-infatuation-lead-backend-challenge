package main

import (
	"testing"
)

// Test that a valid message is parsed correctly.
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
// Test that a malformed message is not parsed.
func TestParseEventBadEvent (t *testing.T) {
	testEvent := "5,LIKE_LIKED,300,100"

	got, err := ParseEvent(testEvent, "|", 4)

	if got != nil && err != nil {
		t.Errorf("Expected nil and err != nil but got %v", got)
	}
}

// Test that matches are found with known good input that contains match events.
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
}

// Test that match events are not found when the input is known to not have match events.
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

	if len(got) != 0 {
		t.Errorf("Expected length of matches to be zero. Got: %d", len(got))
	}
}

// Helper function for uint64 slices.
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
