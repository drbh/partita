package collision

import (
	"sort"
	"sync"
)

type Point struct {
	x, y float64
}

type Segment struct {
	start, end Point
}

type Event struct {
	point   Point
	seg     *Segment
	isStart bool
}

type LineSegmentManager struct {
	events []Event
	mu     sync.Mutex
}

var instance *LineSegmentManager
var once sync.Once

func GetLineSegmentManagerInstance() *LineSegmentManager {
	once.Do(func() {
		instance = &LineSegmentManager{}
	})
	return instance
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

func NewSegment(start, end Point) Segment {
	return Segment{start, end}
}

func NewSegmentFromCoords(x1, y1, x2, y2 float64) Segment {
	return Segment{Point{x1, y1}, Point{x2, y2}}
}

// ClearAllSegments clears all segments from the LineSegmentManager
func (lsm *LineSegmentManager) ClearAllSegments() {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()
	lsm.events = []Event{}
}

func (lsm *LineSegmentManager) AddSegment(s Segment) {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()
	lsm.events = append(lsm.events, Event{point: s.start, seg: &s, isStart: true}, Event{point: s.end, seg: &s, isStart: false})
}

func orientation(p, q, r Point) int {
	val := (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
	if val == 0 {
		return 0
	}
	if val > 0 {
		return 1
	}
	return 2
}

func onSegment(p, q, r Point) bool {
	return q.x <= max(p.x, r.x) && q.x >= min(p.x, r.x) && q.y <= max(p.y, r.y) && q.y >= min(p.y, r.y)
}

func doIntersect(p1, q1, p2, q2 Point) bool {
	o1, o2, o3, o4 := orientation(p1, q1, p2), orientation(p1, q1, q2), orientation(p2, q2, p1), orientation(p2, q2, q1)
	if o1 != o2 && o3 != o4 {
		return true
	}
	return (o1 == 0 && onSegment(p1, p2, q1)) || (o2 == 0 && onSegment(p1, q2, q1)) || (o3 == 0 && onSegment(p2, p1, q2)) || (o4 == 0 && onSegment(p2, q1, q2))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (lsm *LineSegmentManager) CheckIntersection(target Segment) *Segment {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	segments := append([]Event(nil), lsm.events...)
	segments = append(segments, Event{point: target.start, seg: &target, isStart: true}, Event{point: target.end, seg: &target, isStart: false})

	sort.Slice(segments, func(i, j int) bool {
		if segments[i].point.x == segments[j].point.x {
			return segments[i].point.y < segments[j].point.y
		}
		return segments[i].point.x < segments[j].point.x
	})

	active := make(map[*Segment]struct{})
	for _, e := range segments {
		if e.isStart {
			for seg := range active {

				// log.Printf("Checking intersection between %v and %v\n", seg, e.seg)

				if doIntersect(seg.start, seg.end, e.seg.start, e.seg.end) {
					if e.seg == &target || seg == &target {
						return seg
					}
				}
			}
			active[e.seg] = struct{}{}
		} else {
			delete(active, e.seg)
		}
	}
	return nil
}
