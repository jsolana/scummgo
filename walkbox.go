package scummgo

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Walkbox refers to a convex polygonal area that defines the walkable space for actors.
type WalkBox struct {
	walkBoxID string
	enabled   bool
	vertices  [4]Positionf
	scale     float32
	room      *Room
}

// NewWalkBox creates a new WalkBox with the given ID and vertices.
// It ensures the polygon formed by the vertices is convex. If not, it will cause a panic.
// Why convex? Because you can draw a straight line/path between any two vertices inside the polygon
// without needing to implement complex pathfinding algorithms.
func NewWalkBox(id string, vertices [4]Position, scale float32) *WalkBox {
	var verticesf [4]Positionf
	for i, v := range vertices {
		verticesf[i] = v.ToPosf()
	}
	w := &WalkBox{
		walkBoxID: id,
		vertices:  verticesf,
		enabled:   true,
		scale:     quantizeScale(scale),
	}

	if !w.isConvex() {
		log.Panicf("walkbox must be a convex polygon: %v", vertices)
	}
	return w
}

// scaleSteps is the number of steps of a whole unit of scale. Scales are limited to multiples of
// 1/scaleSteps so scaled pixel art keeps a regular look.
const scaleSteps = 16

// quantizeScale rounds s to the nearest multiple of 1/scaleSteps, with a minimum of one step.
func quantizeScale(s float32) float32 {
	steps := math.Round(float64(s) * scaleSteps)
	return float32(math.Max(steps, 1)) / scaleSteps
}

// Scale returns the scale factor applied to the actors standing in the WalkBox. It is always a
// multiple of 1/16.
func (w *WalkBox) Scale() float32 {
	return w.scale
}

// yRange returns the top and bottom of the WalkBox.
func (w *WalkBox) yRange() (top, bottom float32) {
	top, bottom = w.vertices[0].Y, w.vertices[0].Y
	for _, v := range w.vertices {
		top, bottom = min(top, v.Y), max(bottom, v.Y)
	}
	return top, bottom
}

// Draws the edges of the WalkBox.
func (w *WalkBox) draw() {
	numVertices := len(w.vertices)
	for i := 0; i < numVertices; i++ {
		p1 := w.vertices[i]
		p2 := w.vertices[(i+1)%numVertices]
		rl.DrawLineEx(p1.toRaylib(), p2.toRaylib(), 1.2, rl.NewColor(0x55, 0xFF, 0x55, 0x7D))
	}
}

// isConvex check if the current WalkBox is a convex poligon.
func (w *WalkBox) isConvex() bool {
	numVertices := len(w.vertices)

	var totalCrossProduct float32
	var polygonDirection bool // true if clockwise, false if counter-clockwise
	for i := 0; i < numVertices; i++ {
		// Get three consecutive vertices (cyclically)
		p1 := w.vertices[i]
		p2 := w.vertices[(i+1)%numVertices]
		p3 := w.vertices[(i+2)%numVertices]

		cp := p1.CrossProduct(p2, p3)

		if cp == 0 {
			continue // Skip collinear vertices
		}

		totalCrossProduct += cp

		if i == 0 {
			polygonDirection = cp > 0
		} else {
			if (cp > 0) != polygonDirection {
				return false // If direction changes, the polygon is not convex
			}
		}
	}
	return totalCrossProduct != 0
}

// containsPoint check if the provided position is in the boundaries defined by the WalkBox.
func (w *WalkBox) containsPoint(p Positionf) bool {
	// Check if the position is one of the vertices
	for _, vertex := range w.vertices {
		if p.Equals(vertex) {
			return true
		}
	}
	numberOfIntersections := 0
	numVertices := len(w.vertices)

	for i := 0; i < numVertices; i++ {
		p1 := w.vertices[i]
		p2 := w.vertices[(i+1)%numVertices]

		if p.IsIntersecting(p1, p2) {
			numberOfIntersections++
		}
	}

	return numberOfIntersections%2 == 1 // Odd count means inside
}

// isAdjacent checks if two WalkBoxes are adjacent. It returns false if either WalkBox is disabled.
func (w *WalkBox) isAdjacent(otherWalkBox *WalkBox) bool {
	if w.enabled && otherWalkBox.enabled {
		for _, vertex := range otherWalkBox.vertices {
			if w.containsPoint(vertex) {
				return true
			}
		}

		// two-way verification
		for _, vertex := range w.vertices {
			if otherWalkBox.containsPoint(vertex) {
				return true
			}
		}
	}
	return false
}

// distance calculates the shortest distance from the WalkBox to the given position.
func (wb *WalkBox) distance(p Positionf) float32 {
	numVertices := len(wb.vertices)

	var minDistance float32 = math.MaxFloat32
	for i := 0; i < numVertices; i++ {
		p1 := wb.vertices[i]
		p2 := wb.vertices[(i+1)%numVertices]
		closestPoint := p.ClosestPointOnSegment(p1, p2)
		currentDistance := p.Distance(closestPoint)

		if currentDistance < minDistance {
			minDistance = currentDistance
		}
	}

	return minDistance
}

// WalkBoxMatrix represents a collection of WalkBoxes and their adjacency relationships.
type WalkBoxMatrix struct {
	walkBoxes       []*WalkBox
	itineraryMatrix [][]int
}

const (
	// infinityDistance represents a maximum distance value used for unconnected paths.
	infinityDistance = 255
	// InvalidWalkBox indicates an invalid WalkBox ID, typically used to signify non-existence.
	InvalidWalkBox = -1
)

// NewWalkBoxMatrix creates and returns a new WalkBoxMatrix instance
func NewWalkBoxMatrix(walkboxes []*WalkBox) *WalkBoxMatrix {
	wm := &WalkBoxMatrix{
		walkBoxes: walkboxes,
	}

	wm.resetItinerary()
	return wm
}

// WalkBoxes draw walkable boxes of the WalkBoxMatrix.
func (wm *WalkBoxMatrix) Draw() {
	for _, wb := range wm.walkBoxes {
		if wb.enabled {
			wb.draw()
		}
	}
}

// resetItinerary computes the shortest paths between WalkBoxes and returns the resulting itinerary matrix.
func (wm *WalkBoxMatrix) resetItinerary() {
	numBoxes := len(wm.walkBoxes)
	distanceMatrix := make([][]int, numBoxes)
	itineraryMatrix := make([][]int, numBoxes)

	for i, walkbox := range wm.walkBoxes {
		itineraryMatrix[i] = make([]int, numBoxes)
		distanceMatrix[i] = make([]int, numBoxes)

		// Initialize the distance matrix: each box has distance 0 to itself,
		// and distance 1 to its direct neighbors. Initially, it has distance
		// 255 (= infinityDistance) to all other boxes.
		for j, otherWalkBox := range wm.walkBoxes {
			if i == j {
				distanceMatrix[i][j] = 0
				itineraryMatrix[i][j] = i
			} else if walkbox.isAdjacent(otherWalkBox) {
				distanceMatrix[i][j] = 1
				itineraryMatrix[i][j] = j
			} else {
				distanceMatrix[i][j] = infinityDistance
				itineraryMatrix[i][j] = InvalidWalkBox
			}
		}
	}

	// Compute the shortest routes between boxes via Kleene's algorithm.
	for i := range wm.walkBoxes {
		for j := range wm.walkBoxes {
			for k := range wm.walkBoxes {
				distIK := distanceMatrix[i][k]
				distKJ := distanceMatrix[k][j]
				if distanceMatrix[i][j] > distIK+distKJ {
					distanceMatrix[i][j] = distIK + distKJ
					itineraryMatrix[i][j] = k
				}
			}
		}
	}

	wm.itineraryMatrix = itineraryMatrix
}

// EnableWalkBox enables or disables the specified walk box and recalculates the itinerary matrix.
func (wm *WalkBoxMatrix) EnableWalkBox(id string, enabled bool) {
	for _, w := range wm.walkBoxes {
		if w.walkBoxID == id {
			w.enabled = enabled
			wm.resetItinerary()
			break
		}
	}
}

// WayPoint represents a point and its associated WalkBox.
type WayPoint struct {
	Walkbox  *WalkBox
	Position Position
}

// FindPath calculates and returns a path as a sequence of waypoints from the
// starting point 'from' to the destination 'to' within the walk box matrix.
// The path is returned as a slice of waypoints.
func (wm *WalkBoxMatrix) FindPath(from, to Position) []*WayPoint {
	var path []*WayPoint
	fromf := from.ToPosf()
	tof := to.ToPosf()
	current, _ := wm.walkBoxAt(fromf)
	target, _ := wm.walkBoxAt(tof)

	path = append(path, &WayPoint{Walkbox: wm.walkBoxes[current], Position: from})
	for current != target {
		next := wm.nextWalkBox(current, target)
		if next == InvalidWalkBox {
			break
		}
		nextPosition := wm.closestPositionToWalkBox(fromf, next)
		path = append(path, &WayPoint{Walkbox: wm.walkBoxes[next], Position: nextPosition.ToPos()})
		current = next
		fromf = nextPosition
	}

	path = append(path, &WayPoint{Walkbox: wm.walkBoxes[current], Position: wm.closestPositionOnWalkBox(current, tof).ToPos()})
	return path
}

// nextWalkBox returns the next walk box in the path from the source to the destination.
func (wm *WalkBoxMatrix) nextWalkBox(from, to int) int {
	if from < 0 || from >= len(wm.walkBoxes) || to < 0 || to >= len(wm.walkBoxes) {
		return InvalidWalkBox
	}

	if wm.itineraryMatrix[from][to] == to {
		return to
	}

	next := wm.itineraryMatrix[from][to]
	return wm.nextWalkBox(from, next)
}

// walkBoxAt returns the walk box identifier at the given position or the closest one,
// along with a boolean indicating inclusion. If the point is located between two or more
// boxes, it returns the lowest walk box ID among them.
func (wm *WalkBoxMatrix) walkBoxAt(p Positionf) (id int, included bool) {
	var minDistance float32 = math.MaxFloat32
	id = InvalidWalkBox
	for i, wb := range wm.walkBoxes {
		if included = wb.containsPoint(p); included {
			return i, included
		}

		current := wb.distance(p)
		if current < minDistance {
			id = i
			minDistance = current
		}
	}

	return id, false
}

// scaleBlend is the fraction of the height of a WalkBox, measured from each edge it shares with
// another WalkBox, where its scale blends with the scale of the neighbour.
const scaleBlend = 1.0 / 3

// ScaleAt returns the scale of an actor standing at p. It is the scale of the WalkBox at p, except
// near the edges shared with the adjacent WalkBoxes above and below: there the scale moves
// gradually towards the middle value between both boxes, so crossing from one box to the other
// does not make the actor jump in size. The result is a multiple of 1/16.
func (wm *WalkBoxMatrix) ScaleAt(p Positionf) float32 {
	id, _ := wm.walkBoxAt(p)
	if id == InvalidWalkBox {
		return 1
	}
	wb := wm.walkBoxes[id]
	top, bottom := wb.yRange()
	margin := (bottom - top) * scaleBlend

	// Average scale of the neighbours above and below the box.
	var sum [2]float32
	var count [2]int
	for i, other := range wm.walkBoxes {
		if i == id || !wb.isAdjacent(other) {
			continue
		}
		otop, obottom := other.yRange()
		side := 0 // above
		if otop+obottom > top+bottom {
			side = 1 // below
		}
		sum[side] += other.scale
		count[side]++
	}

	scale := wb.scale
	blend := func(side int, distance float32) {
		if count[side] == 0 || distance >= margin {
			return
		}
		edge := (wb.scale + sum[side]/float32(count[side])) / 2
		t := max(distance, 0) / margin
		scale = edge + (wb.scale-edge)*t
	}
	blend(0, p.Y-top)
	blend(1, bottom-p.Y)
	return quantizeScale(scale)
}

// closestPositionOnWalkBox returns the closest point on the walkbox at a given position.
func (wm *WalkBoxMatrix) closestPositionOnWalkBox(from int, p Positionf) Positionf {
	wb := wm.walkBoxes[from]
	if wb.containsPoint(p) {
		return p
	}

	return wm.closestPositionToWalkBox(p, from) // looking for the closest edge point
}

// closestPositionToWalkBox returns the nearest point on the edge of the specified walkbox
// from a given position p.
func (wm *WalkBoxMatrix) closestPositionToWalkBox(p Positionf, to int) Positionf {
	var minDistance float32 = math.MaxFloat32
	var closestPoint Positionf
	wb := wm.walkBoxes[to]
	numVertices := len(wb.vertices)

	for i := 0; i < numVertices; i++ {
		p1 := wb.vertices[i]
		p2 := wb.vertices[(i+1)%numVertices]
		currentPoint := p.ClosestPointOnSegment(p1, p2)
		currentDistance := p.Distance(currentPoint)

		if currentDistance < minDistance {
			minDistance = currentDistance
			closestPoint = currentPoint
		}
	}

	return closestPoint
}
