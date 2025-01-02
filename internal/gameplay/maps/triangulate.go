package maps

import (
	stdmath "math"
	"slices"

	"github.com/efritz/lunar-fever/internal/common/math"
)

// triangulate decomposes a set of axis-aligned polygons into a set of
// triangles that cover the same area. The input bounds may be mutated.
func triangulate(bounds []Bound) []Bound {
	var triangles []Bound
	for _, b := range bounds {
		for _, vertices := range triangulatePolygon(b.Vertices) {
			triangles = append(triangles, newBound(vertices...))
		}
	}

	return triangles
}

func triangulatePolygon(vertices []math.Vector) [][]math.Vector {
	n := len(vertices)
	if n < 3 {
		panic("too few vertices")
	}
	if n == 3 {
		return [][]math.Vector{vertices}
	}

	for i, curr := range vertices {
		prev := vertices[prevVertexIndex(i, n)]
		next := vertices[nextVertexIndex(i, n)]

		if isClippableEar(vertices, prev, curr, next) {
			// Extract clipped ear and recursively triangulate remaining polygon
			ear := []math.Vector{prev, curr, next}
			rest := triangulatePolygon(slices.Delete(vertices, i, i+1))
			return append(rest, ear)
		}
	}

	panic("failed to triangulate")
}

func isClippableEar(vertices []math.Vector, a, b, c math.Vector) bool {
	return isConvexAngle(a, b, c) && !isAxisAlignedColinearLine(a, b, c) && !anyVertexInTriangle(vertices, a, b, c)
}

func isConvexAngle(a, b, c math.Vector) bool {
	v1 := c.Sub(b)
	v2 := a.Sub(b)

	angle := math.Atan232(v1.Cross(v2), v1.Dot(v2))
	if angle < 0 {
		angle += 2 * stdmath.Pi
	}

	return angle <= stdmath.Pi
}

func anyVertexInTriangle(vertices []math.Vector, a, b, c math.Vector) bool {
	for _, check := range vertices {
		if !check.Equal(a) && !check.Equal(b) && !check.Equal(c) && PointInTriangle(a, b, c, check) {
			return true
		}
	}

	return false
}

func PointInTriangle(a, b, c, p math.Vector) bool {
	v0 := c.Sub(a)
	v1 := b.Sub(a)
	v2 := p.Sub(a)

	dot00 := v0.Dot(v0)
	dot01 := v0.Dot(v1)
	dot02 := v0.Dot(v2)
	dot11 := v1.Dot(v1)
	dot12 := v1.Dot(v2)

	denom := (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) / denom
	v := (dot00*dot12 - dot01*dot02) / denom

	return (u >= 0) && (v >= 0) && (u+v < 1)
}
