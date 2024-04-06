package physics

import (
	"slices"

	"github.com/efritz/lunar-fever/internal/common/math"
)

type Fixture struct {
	Vertices        []math.Vector
	normals         []math.Vector
	density         float32
	restitution     float32
	staticFriction  float32
	dynamicFriction float32
}

func NewBasicFixture(x, y, w, h, density, restitution, staticFriction, dynamicFriction float32) Fixture {
	vertices := []math.Vector{
		{X: x - w, Y: y - h},
		{X: x + w, Y: y - h},
		{X: x - w, Y: y + h},
		{X: x + w, Y: y + h},
	}

	return NewFixture(vertices, density, restitution, staticFriction, dynamicFriction)
}

func NewFixture(vertices []math.Vector, density, restitution, staticFriction, dynamicFriction float32) Fixture {
	sortByPolarAngle(vertices)

	var normals []math.Vector
	for i, v1 := range vertices {
		v2 := vertices[(i+1)%len(vertices)]
		normals = append(normals, v2.Sub(v1).Normalize().Orthogonalize())
	}

	return Fixture{
		Vertices:        vertices,
		normals:         normals,
		density:         density,
		restitution:     restitution,
		staticFriction:  staticFriction,
		dynamicFriction: dynamicFriction,
	}
}

func (f Fixture) VertexInWorldSpace(body *Body, index int) math.Vector {
	return body.Rotation.Mul(f.Vertices[index]).Add(body.Position)
}

func (f Fixture) NormalInWorldSpace(body *Body, index int) math.Vector {
	return body.Rotation.Mul(f.normals[index])
}

func sortByPolarAngle(vertices []math.Vector) {
	centroid := centroid(vertices)

	polarAngleByVector := map[math.Vector]float32{}
	for _, v := range vertices {
		polarAngleByVector[v] = math.Atan232(v.Y-centroid.Y, v.X-centroid.X)
	}

	slices.SortFunc(vertices, func(a, b math.Vector) int {
		return math.Sign(polarAngleByVector[a] - polarAngleByVector[b])
	})
}

func centroid(vertices []math.Vector) math.Vector {
	centroid := math.Vector{}
	for _, v := range vertices {
		centroid = centroid.Add(v)
	}

	return centroid.Divs(float32(len(vertices)))
}
