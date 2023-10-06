package collision

import (
	"testing"
)

func TestAddSegment(t *testing.T) {
	manager := GetLineSegmentManagerInstance()
	manager.AddSegment(Segment{Point{5, 0}, Point{5, 10}})
	if len(manager.events) != 2 {
		t.Errorf("AddSegment failed, expected %v, got %v", 2, len(manager.events))
	}
}

func TestCheckIntersection(t *testing.T) {
	manager := GetLineSegmentManagerInstance()
	manager.AddSegment(Segment{Point{5, 0}, Point{5, 10}})
	target := Segment{Point{0, 5}, Point{10, 5}}
	intersectingSegment := manager.CheckIntersection(target)
	if intersectingSegment == nil {
		t.Errorf("CheckIntersection failed, expected %v, got %v", "not nil", "nil")
	}
}
