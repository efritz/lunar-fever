package maps

import "github.com/efritz/lunar-fever/internal/common/math"

// simplifyBounds removes collinear vertices from the given bounds.
// This mutates the input bounds in place.
func simplifyBounds(bounds []Bound) []Bound {
	for i, b := range bounds {
		bounds[i].Vertices = removeCollinearVertices(b.Vertices)
	}

	return bounds
}

func removeCollinearVertices(vertices []math.Vector) []math.Vector {
	if len(vertices) <= 2 {
		return vertices
	}

	var simplified []math.Vector
	for _, v := range vertices {
		if n := len(simplified); n >= 2 {
			// If the previous two vertices and the new vertex would create a colllinear line,
			// then remove the middle vertex (the vertex at the top of the simplified stack).
			if isAxisAlignedColinearLine(simplified[n-2], simplified[n-1], v) {
				simplified = simplified[:n-1]
			}
		}

		simplified = append(simplified, v)
	}

	if n := len(simplified); n >= 3 {
		// Check if there is a collinear line that wraps around the vertex list. If so, remove
		// the middle vertex, just as we do above. This can happen in two distinct, mutually
		// esxclusive scenarios:
		//
		// (1) The last two vertices and the first vertex form a line, or
		// (2) The last vertex, the first two vertices form a line.
		//
		// Each scenario has a different middle vertex to remove.

		if isAxisAlignedColinearLine(simplified[n-2], simplified[n-1], simplified[0]) {
			simplified = simplified[:n-1]
		} else if isAxisAlignedColinearLine(simplified[n-1], simplified[0], simplified[1]) {
			simplified = simplified[1:]
		}
	}

	return simplified
}
