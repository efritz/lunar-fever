package maps

import (
	"github.com/efritz/lunar-fever/internal/common/math"
)

// mergeBounds merges pairs of bounds that share a common edge sequence.
// This returns a slice of bounds with the minimum number of bounds. The
// input bounds may be mutated.
func mergeBounds(bounds []Bound, obstacles []Edge) []Bound {
outer:
	for {
		for i := 0; i < len(bounds); i++ {
			for j := i + 1; j < len(bounds); j++ {
				merged, ok := tryMergeBounds(bounds[i], bounds[j], obstacles)
				if !ok {
					continue
				}

				// Replace first bound with merged bound
				bounds[i] = merged

				// Remove second bound from slice
				n := len(bounds)
				bounds[j] = bounds[n-1]
				bounds = bounds[:n-1]

				// We've just changed our indexes; start over
				continue outer
			}
		}

		break
	}

	return bounds
}

// tryMergeBounds attempts to merge two bounds that share a common edge
// sequence. If the bounds share a common edge sequence, the merged bound
// is returned along with a true-valued flag.
func tryMergeBounds(b1, b2 Bound, obstacles []Edge) (Bound, bool) {
	n := len(b1.Vertices)
	m := len(b2.Vertices)

	b1Start, b1End, b2Start, b2End, ok := sharedEdgeVertexIndexes(b1, b2, obstacles)
	if !ok {
		return Bound{}, false
	}

	var combinedVertices []math.Vector
	for i := 0; i < n; i++ {
		// This is the start of the shared edge sequence
		if i != b1Start {
			// This is a non-shared vertex from b1; add it normally
			combinedVertices = append(combinedVertices, b1.Vertices[i])
		} else {
			// We're now at the start of the shared edge sequence in b1.
			//
			// (1) We'll start by adding the vertex that ends the shared edge sequence.
			// (2) Then, we'll trace the entirety of the non-shared portion of b2.
			// (3) Finally, we'll add the vertex that starts the shared edge sequence in
			//     b2, which puts us back on the non-shared portion of b1.

			for j := b2End; ; j = nextVertexIndex(j, m) {
				combinedVertices = append(combinedVertices, b2.Vertices[j])

				if j == b2Start {
					break
				}
			}

			// Skip over the shared edge sequence completely in the outer loop. On the
			// following iteration, we'll see the vertex directly after the shared edge
			// sequence in b1.

			i = b1End
		}
	}

	return newBound(combinedVertices...), true
}

// sharedEdgeVertexIndexes returns the the indices of a unique contiguous edge
// sequence shared between the give bounds. If there are additional disjoint
// shared features, no shared edge sequence is considered valid.
func sharedEdgeVertexIndexes(b1, b2 Bound, obstacles []Edge) (b1Start, b1End, b2Start, b2End int, _ bool) {
	n := len(b1.Vertices)
	m := len(b2.Vertices)

	var b1Indexes []int
	var b2Indexes []int

	for i, v1 := range b1.Vertices {
		for j, v2 := range b2.Vertices {
			if v1.Equal(v2) {
				b1Indexes = append(b1Indexes, i)
				b2Indexes = append(b2Indexes, j)
			}
		}
	}

	// Ensure we have at least two vertex indexes that can form an edge.
	// Otherwise we have two bounds that are non-adjacent or trivially
	// share a single vertex.

	if len(b1Indexes) < 2 {
		return 0, 0, 0, 0, false
	}

	// Ensure that each of the shared vertex indexes are in order and contiguous.
	// If this is not the case then there are disjoint shared features that would
	// likely create a hole when combined into a single polygon.

	for i, index := range b1Indexes {
		// Ensure next(previous) = index
		if i > 0 && index != nextVertexIndex(b1Indexes[i-1], n) {
			return 0, 0, 0, 0, false
		}
	}

	for i, index := range b2Indexes {
		// Ensure previous = next(index)
		//
		// By construction, tile-based vertex sequences have a clockwise orientation.
		// This is true for trivial single-tile bounds and is maintained for merged
		// bounds. This means that the shared vertex sequences in two adjacent bounds
		// will approach each other in opposite directions.

		if i > 0 && b2Indexes[i-1] != nextVertexIndex(index, m) {
			return 0, 0, 0, 0, false
		}
	}

	// Lastly, we'll test each of the edges in the shared edge sequence against all of
	// the obstacles. If any of the edges overlap with any obstacle, then we don't want
	// to merge the two bounds together, or we'll be erasing level information.
	//
	// Note that we don't do the usual wraparound iteration here as this is a shared edge
	// subsequence, not a full polygon.

	for i := 0; i < len(b1Indexes)-1; i++ {
		a := b1.Vertices[(b1Indexes[0]+i)%n]
		b := b1.Vertices[(b1Indexes[0]+i+1)%n]

		for _, obstacle := range obstacles {
			if axisAlignedSegmentsOverlap(a, b, obstacle.From, obstacle.To) {
				return 0, 0, 0, 0, false
			}
		}
	}

	return b1Indexes[0], b1Indexes[len(b1Indexes)-1], b2Indexes[len(b1Indexes)-1], b2Indexes[0], true
}

func axisAlignedSegmentsOverlap(i, ii, j, jj math.Vector) bool {
	if i.X == ii.X {
		if j.X == jj.X && i.X == j.X {
			minI := math.Min(i.Y, ii.Y)
			maxI := math.Max(i.Y, ii.Y)
			minJ := math.Min(j.Y, jj.Y)
			maxJ := math.Max(j.Y, jj.Y)

			return minI < maxJ && maxI > minJ
		}
	}

	if i.Y == ii.Y {
		if j.Y == jj.Y && i.Y == j.Y {
			minI := math.Min(i.X, ii.X)
			maxI := math.Max(i.X, ii.X)
			minJ := math.Min(j.X, jj.X)
			maxJ := math.Max(j.X, jj.X)

			return minI < maxJ && maxI > minJ
		}
	}

	return false
}
