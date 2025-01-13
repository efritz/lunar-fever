package maps

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/engelsjk/polygol"
)

const wallExtents = float32(32)

func expandObstacleEdge(obstacle Edge) Bound {
	if obstacle.To.X == obstacle.From.X {
		if obstacle.From.Y > obstacle.To.Y {
			// sanity check
			panic("wall is not top to bottom")
		}

		return newBound(
			math.Vector{obstacle.From.X - wallExtents, obstacle.From.Y},
			math.Vector{obstacle.To.X + wallExtents, obstacle.From.Y},
			math.Vector{obstacle.To.X + wallExtents, obstacle.To.Y},
			math.Vector{obstacle.From.X - wallExtents, obstacle.To.Y},
		)
	} else {
		if obstacle.From.X > obstacle.To.X {
			// sanity check
			panic("wall is not left to right")
		}

		return newBound(
			math.Vector{obstacle.From.X, obstacle.From.Y - wallExtents},
			math.Vector{obstacle.To.X, obstacle.From.Y - wallExtents},
			math.Vector{obstacle.To.X, obstacle.To.Y + wallExtents},
			math.Vector{obstacle.From.X, obstacle.To.Y + wallExtents},
		)
	}

	panic("malformed edge")
}

func subtract(bounds []Bound, obstacles []Edge) []Bound {
	var obstacleBounds []Bound
	for _, obstacle := range obstacles {
		obstacleBounds = append(obstacleBounds, expandObstacleEdge(obstacle))
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
