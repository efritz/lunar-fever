package maps

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/engelsjk/polygol"
)

const obstacleExtents = float32(32)

func expandObstacleEdge(obstacle Edge) Bound {
	if obstacle.To.X == obstacle.From.X {
		if obstacle.From.Y > obstacle.To.Y {
			// sanity check
			panic("wall is not top to bottom")
		}

		return newBound(
			math.Vector{obstacle.From.X - obstacleExtents, obstacle.From.Y},
			math.Vector{obstacle.To.X + obstacleExtents, obstacle.From.Y},
			math.Vector{obstacle.To.X + obstacleExtents, obstacle.To.Y},
			math.Vector{obstacle.From.X - obstacleExtents, obstacle.To.Y},
		)
	} else if obstacle.To.Y == obstacle.From.Y {
		if obstacle.From.X > obstacle.To.X {
			// sanity check
			panic("wall is not left to right")
		}

		return newBound(
			math.Vector{obstacle.From.X, obstacle.From.Y - obstacleExtents},
			math.Vector{obstacle.To.X, obstacle.From.Y - obstacleExtents},
			math.Vector{obstacle.To.X, obstacle.To.Y + obstacleExtents},
			math.Vector{obstacle.From.X, obstacle.To.Y + obstacleExtents},
		)
	}

	panic("malformed edge")
}

func subtract(bounds []Bound, obstacles []Edge, fixtures []Bound) []Bound {
	var obstacleBounds []Bound
	for _, obstacle := range obstacles {
		obstacleBounds = append(obstacleBounds, expandObstacleEdge(obstacle))
	}

	for _, fixture := range fixtures {
		minX, minY := float32(+stdmath.MaxFloat32), float32(+stdmath.MaxFloat32)
		maxX, maxY := float32(-stdmath.MaxFloat32), float32(-stdmath.MaxFloat32)

		for _, vertex := range fixture.Vertices {
			minX = math.Min(minX, vertex.X)
			minY = math.Min(minY, vertex.Y)
			maxX = math.Max(maxX, vertex.X)
			maxY = math.Max(maxY, vertex.Y)
		}

		obstacleBounds = append(obstacleBounds, newBound(
			math.Vector{minX - obstacleExtents, minY - obstacleExtents},
			math.Vector{maxX + obstacleExtents, minY - obstacleExtents},
			math.Vector{maxX + obstacleExtents, maxY + obstacleExtents},
			math.Vector{minX - obstacleExtents, maxY + obstacleExtents},
		))
	}

	for _, obstacle := range obstacleBounds {
		var newBounds []Bound
		for _, bound := range bounds {
			newBounds = append(newBounds, diff(bound, obstacle)...)
		}

		bounds = newBounds
	}

	// return obstacleBounds
	return bounds
}

func diff(a, b Bound) []Bound {
	ax := [][]float64{}
	for _, v := range a.Vertices {
		ax = append(ax, []float64{float64(v.X), float64(v.Y)})
	}
	ax = append(ax, ax[0])

	bx := [][]float64{}
	for _, v := range b.Vertices {
		bx = append(bx, []float64{float64(v.X), float64(v.Y)})
	}
	bx = append(bx, bx[0])

	difference, err := polygol.Difference([][][][]float64{{ax}}, [][][][]float64{{bx}})
	if err != nil {
		panic(err.Error())
	}
	if len(difference) == 0 {
		return nil
	}

	var bs []Bound
	for _, d := range difference {
		for _, d2 := range d {
			var vx []math.Vector
			for _, v := range d2 {
				vx = append(vx, math.Vector{float32(v[0]), float32(v[1])})
			}

			bs = append(bs, newBound(vx...))
		}
	}

	return bs
}
