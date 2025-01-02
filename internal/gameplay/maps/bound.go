package maps

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type Bound struct {
	ID       int
	Vertices []math.Vector
	Color    rendering.Color
}

var boundID = 0

func newBound(vertices ...math.Vector) Bound {
	boundID++

	return Bound{
		ID:       boundID,
		Vertices: vertices,
		Color:    randomColor(),
	}
}

//
//

type Edge struct {
	From math.Vector
	To   math.Vector
}

func newEdge(from, to math.Vector) Edge {
	return Edge{
		From: from,
		To:   to,
	}
}

//
//

func randomColor() rendering.Color {
	return rendering.Color{
		R: math.Random(0, 1),
		G: math.Random(0, 1),
		B: math.Random(0, 1),
		A: 0.25,
	}
}
