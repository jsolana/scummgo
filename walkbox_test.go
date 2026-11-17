package scummgo_test

import (
	"testing"

	"github.com/jsolana/scummgo"
	"github.com/stretchr/testify/assert"
)

const (
	DefaultWalkBoxID = "walkbox"
	DefaultScale     = 1.0
)

func TestNewWalkBox(t *testing.T) {
	testCases := []struct {
		name        string
		vertices    [4]scummgo.Position
		shouldPanic bool
		message     string
	}{
		{
			name:        "Concave polygon should panic",
			vertices:    [4]scummgo.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 2, Y: 1}, {X: 4, Y: 4}},
			shouldPanic: true,
			message:     "Expected panic because vertices form a concave polygon!",
		},
		{
			name:        "Collinear vertices should panic",
			vertices:    [4]scummgo.Position{{X: 1, Y: 2}, {X: 2, Y: 4}, {X: 3, Y: 6}, {X: 4, Y: 8}},
			shouldPanic: true,
			message:     "Expected panic because vertices are collinear!",
		},
		{
			name:        "Should successfully create a valid WalkBox with a convex polygon",
			vertices:    [4]scummgo.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}},
			shouldPanic: false,
			message:     "Expected create a valid WalkBox, vertices form a convex polygon!",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.shouldPanic {
				assert.Panics(t, func() {
					scummgo.NewWalkBox(DefaultWalkBoxID, testCase.vertices, DefaultScale)
				}, testCase.message)
			} else {
				assert.NotPanics(t, func() {
					scummgo.NewWalkBox(DefaultWalkBoxID, testCase.vertices, DefaultScale)
				}, testCase.message)
			}
		})
	}

}

func TestWalkBoxScale(t *testing.T) {
	vertices := [4]scummgo.Position{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4}}

	testCases := []struct {
		name     string
		scale    float32
		expected float32
	}{
		{name: "Full size is kept", scale: 1, expected: 1},
		{name: "Multiples of 1/16 are kept", scale: 13.0 / 16, expected: 13.0 / 16},
		{name: "Scale is rounded down to the nearest 1/16", scale: 0.95, expected: 15.0 / 16},
		{name: "Scale is rounded up to the nearest 1/16", scale: 0.3, expected: 5.0 / 16},
		{name: "Scale never goes below 1/16", scale: 0, expected: 1.0 / 16},
		{name: "Scale can be bigger than one", scale: 1.5, expected: 1.5},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			wb := scummgo.NewWalkBox(DefaultWalkBoxID, vertices, testCase.scale)
			assert.Equal(t, testCase.expected, wb.Scale())
		})
	}
}

func TestWalkBoxMatrixScaleAt(t *testing.T) {
	far := scummgo.NewWalkBox("far", [4]scummgo.Position{{X: 0, Y: 0}, {X: 10, Y: 0}, {X: 10, Y: 30}, {X: 0, Y: 30}}, 4.0/16)
	near := scummgo.NewWalkBox("near", [4]scummgo.Position{{X: 0, Y: 30}, {X: 10, Y: 30}, {X: 10, Y: 60}, {X: 0, Y: 60}}, 1)
	matrix := scummgo.NewWalkBoxMatrix([]*scummgo.WalkBox{far, near})

	testCases := []struct {
		name     string
		y        float32
		expected float32
	}{
		{name: "Far from the shared edge the scale of the box is used", y: 45, expected: 1},
		{name: "Far from the shared edge in the other box", y: 5, expected: 4.0 / 16},
		{name: "On the shared edge the scale is the middle value", y: 30.5, expected: 10.0 / 16},
		{name: "Close to the shared edge the scale is blended", y: 35, expected: 13.0 / 16},
		{name: "At the end of the blending area the scale of the box is used", y: 40, expected: 1},
		{name: "Blending also happens in the other box", y: 25, expected: 7.0 / 16},
		{name: "The edges not shared with other boxes do not blend", y: 60, expected: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, matrix.ScaleAt(scummgo.NewPosf(5, testCase.y)))
		})
	}
}

func TestFindPath(t *testing.T) {
	/*
		Polygons disposition:

		          +-------+
		          |       |
		  +-------|       |+-------
		  |       |       |       |
		  | box6  | box7  | box8  |
		  |       |       |       |
		  +-------|       |+-------
		          |       |
		  +-------+-------+-------+
		  |       |       |       |
		  | box3  | box4  | box5  |
		  |       |       |       |
		  +-------+-------+-------+
		  |       |       |       |
		  | box0  | box1  | box2  |
		  |       |       |       |
		  +-------+-------+-------+

		Each box represents a square, except box7, which is three times taller.
		- box0 is adjacent to box1, box3
		- box1 is adjacent to box0, box2, box3, box4
		- box2 is adjacent to box1, box5
		- box3 is adjacent to box0, box1, box4
		- box4 is adjacent to box1, box3, box5, box7
		- box5 is adjacent to box2, box4
		- box6 is adjacent to box7 (positioned above box3 but not connected)
		- box7 is adjacent to box4, box6 (taller and positioned above box4)
		- box8 is adjacent to box7 (positioned above box5 but not connected)
	*/

	testCases := []struct {
		name       string
		from       scummgo.Position
		to         scummgo.Position
		expectedTo scummgo.Position
		assertFunc func(t *testing.T, path []*scummgo.WayPoint, expectedTo scummgo.Position)
	}{
		{
			name:       "Should return a valid path when 'from' and 'to' are inside walk boxes",
			from:       scummgo.Position{X: 0, Y: 0}, // inside box0
			to:         scummgo.Position{X: 1, Y: 2}, // inside box7
			expectedTo: scummgo.Position{X: 1, Y: 2}, // expected return to
			assertFunc: func(t *testing.T, path []*scummgo.WayPoint, expectedTo scummgo.Position) {
				assert.True(t, path[len(path)-1].Position.Equals(expectedTo))
			},
		},
		{
			name:       "Should return the closest point when 'to' is outside the closest walk box",
			from:       scummgo.Position{X: 0, Y: 0}, // inside box0
			to:         scummgo.Position{X: 3, Y: 1}, // outside all boxes, close to box5
			expectedTo: scummgo.Position{X: 3, Y: 1}, // expected to return the closest point inside box5
			assertFunc: func(t *testing.T, path []*scummgo.WayPoint, expectedTo scummgo.Position) {
				assert.True(t, path[len(path)-1].Position.Equals(expectedTo))
			},
		},
	}

	box0 := scummgo.NewWalkBox("walkbox0", [4]scummgo.Position{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, DefaultScale)
	box1 := scummgo.NewWalkBox("walkbox1", [4]scummgo.Position{{1, 0}, {2, 0}, {2, 1}, {1, 1}}, DefaultScale)
	box2 := scummgo.NewWalkBox("walkbox2", [4]scummgo.Position{{2, 0}, {3, 0}, {3, 1}, {2, 1}}, DefaultScale)
	box3 := scummgo.NewWalkBox("walkbox3", [4]scummgo.Position{{0, 1}, {1, 1}, {1, 2}, {0, 2}}, DefaultScale)
	box4 := scummgo.NewWalkBox("walkbox4", [4]scummgo.Position{{1, 1}, {2, 1}, {2, 2}, {1, 2}}, DefaultScale)
	box5 := scummgo.NewWalkBox("walkbox5", [4]scummgo.Position{{2, 1}, {3, 1}, {3, 2}, {2, 2}}, DefaultScale)
	box6 := scummgo.NewWalkBox("walkbox6", [4]scummgo.Position{{1, 4}, {0, 4}, {0, 3}, {1, 3}}, DefaultScale) // starts in top right vertex on purpose
	box7 := scummgo.NewWalkBox("walkbox7", [4]scummgo.Position{{1, 2}, {2, 2}, {2, 5}, {1, 5}}, DefaultScale)
	box8 := scummgo.NewWalkBox("walkbox8", [4]scummgo.Position{{2, 3}, {3, 3}, {3, 4}, {2, 4}}, DefaultScale)

	walkBoxMatrix := scummgo.NewWalkBoxMatrix([]*scummgo.WalkBox{box0, box1, box2, box3, box4, box5, box6, box7, box8})

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			path := walkBoxMatrix.FindPath(testCase.from, testCase.to)
			testCase.assertFunc(t, path, testCase.expectedTo)
		})
	}
}
