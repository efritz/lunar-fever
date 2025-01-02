package maps

import "github.com/efritz/lunar-fever/internal/common/math"

func nextVertexIndex(i, n int) int {
	return (i + 1) % n
}

func prevVertexIndex(i, n int) int {
	return (i - 1 + n) % n
}

// isAxisAlignedColinearLine returns true if all of the given vertices have
// either the same X or the same Y component.
func isAxisAlignedColinearLine(vs ...math.Vector) bool {
	sameX := true
	sameY := true

	for _, v := range vs {
		sameX = sameX && v.X == vs[0].X
		sameY = sameY && v.Y == vs[0].Y
	}

	return sameX || sameY
}
