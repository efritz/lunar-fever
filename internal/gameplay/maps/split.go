package maps

import "github.com/efritz/lunar-fever/internal/common/math"

// splitBoundsAtIntersections adds additional vertices where bounds intersect with doors
// or other bounds in the same set. This is necessary to ensure that after we triangulate
// the bounds that all adjacent edges are equivalent (and not just partially overlapping).
func splitBoundsAtIntersections(bounds []Bound, doors []Edge) []Bound {
	for i := range bounds {
		bounds[i] = splitBoundAtIntersections(bounds, bounds[i], doors)
	}

	return bounds
}

func splitBoundAtIntersections(bounds []Bound, bound Bound, doors []Edge) Bound {
	var queue []Edge

	// queue = append(queue, doors...)
	for _, door := range doors {
		vx := expandDoorEdge(door)
		for i, v := range vx.Vertices {
			queue = append(queue, newEdge(v, vx.Vertices[nextVertexIndex(i, len(vx.Vertices))]))
		}
	}

	for _, other := range bounds {
		if bound.ID != other.ID {
			for i, v := range other.Vertices {
				queue = append(queue, newEdge(v, other.Vertices[nextVertexIndex(i, len(other.Vertices))]))
			}
		}
	}

	for len(queue) > 0 {
		edge := queue[0]
		queue = queue[1:]

		if vertices, ok := splitPolygon(bound.Vertices, edge); ok {
			bound.Vertices = vertices

			// Add the edge back to the queue after a successful split. We may have
			// added only one vertex of two necessary vertices to the polygon.
			queue = append(queue, edge)
		}
	}

	return bound
}

// splitPolygon returns a new set of vertices representing the same polygon but with
// an extra vertex added where the bound intersects the given edge. If no such vertex
// was inserted, an empty vertex set and false-valued flag is returned.
//
// Note that this method inserts at most ONE vertex, even if the edge intersects the
// polygon multiple times. It's up to the caller to invoke this function a appropriate
// number of times to create all the necessary vertices.
func splitPolygon(vertices []math.Vector, edge Edge) ([]math.Vector, bool) {
	n := len(vertices)

	for i, v1 := range vertices {
		v2 := vertices[nextVertexIndex(i, n)]

		if vertices, ok := splitPolygonEdge(vertices, v1, v2, i, edge); ok {
			return vertices, true
		}
	}

	return nil, false
}

func splitPolygonEdge(vertices []math.Vector, v1, v2 math.Vector, i int, edge Edge) ([]math.Vector, bool) {
	if !isAxisAlignedColinearLine(v1, v2, edge.From, edge.To) {
		return nil, false
	}

	for _, intersection := range []math.Vector{v1, v2, edge.From, edge.To} {
		if alteredVertices, ok := insertVertex(vertices, v1, v2, i, intersection); ok {
			return alteredVertices, true
		}
	}

	return nil, false
}

func insertVertex(vertices []math.Vector, v1, v2 math.Vector, i int, newV math.Vector) ([]math.Vector, bool) {
	splitX := v1.Y == v2.Y && math.Min(v1.X, v2.X) < newV.X && newV.X < math.Max(v1.X, v2.X)
	splitY := v1.X == v2.X && math.Min(v1.Y, v2.Y) < newV.Y && newV.Y < math.Max(v1.Y, v2.Y)

	if !splitX && !splitY {
		return nil, false
	}

	return append(append([]math.Vector{newV}, vertices[i+1:]...), vertices[:i+1]...), true
}
