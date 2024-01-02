package physics

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
)

type Body struct {
	Fixtures        []Fixture
	inverseMass     float32
	inverseInertia  float32
	Position        math.Vector
	LinearVelocity  math.Vector
	AngularVelocity float32
	Orient          float32
	rotation        math.Matrix32
	force           math.Vector
	torque          float32
}

func NewBody(fixtures []Fixture) *Body {
	mass := float32(0)
	inertia := float32(0)
	centroid := math.Vector{}

	for _, fixture := range fixtures {
		fixtureCentroid := math.Vector{}
		fixtureArea := float32(0)
		fixtureInertia := float32(0)

		for i, p1 := range fixture.Vertices {
			p2 := fixture.Vertices[(i+1)%len(fixture.Vertices)]

			triangleArea := p1.X*p2.Y - p2.X*p1.Y

			fixtureArea += triangleArea
			fixtureCentroid = fixtureCentroid.Add(p1.Muls(triangleArea / 6))
			fixtureCentroid = fixtureCentroid.Add(p2.Muls(triangleArea / 6))

			x2 := p1.X*p1.X + p1.X*p2.X + p2.X*p2.X
			y2 := p1.Y*p1.Y + p1.Y*p2.Y + p2.Y*p2.Y
			fixtureInertia += (x2 + y2) * triangleArea
		}

		fixtureArea /= 2
		fixtureInertia /= 12

		fixtureCentroid = fixtureCentroid.Divs(fixtureArea)

		fixtureMass := fixtureArea * fixture.density
		fixtureInertia = fixtureInertia * fixture.density

		centroid = centroid.Add(fixtureCentroid.Muls(fixtureMass))

		mass += fixtureMass
		inertia += fixtureInertia
	}

	if mass != 0 {
		centroid = centroid.Divs(mass)
	}

	inverseMass := float32(0)
	if mass != 0 {
		inverseMass = 1 / mass
	}

	inverseInertia := float32(0)
	if inertia != 0 {
		inverseInertia = 1 / inertia
	}

	for _, fixture := range fixtures {
		for i, v := range fixture.Vertices {
			fixture.Vertices[i] = v.Sub(centroid)
		}
	}

	return &Body{
		Fixtures:       fixtures,
		inverseMass:    inverseMass,
		inverseInertia: inverseInertia,
		rotation:       math.Matrix32{1, 0, 0, 1},
	}
}

func (b *Body) SetOrient(radians float32) {
	c := math.Cos32(radians)
	s := math.Sin32(radians)

	b.Orient = radians
	b.rotation = math.Matrix32{c, -s, s, c}
}

// TODO - unused?
func (b *Body) ApplyForce(force math.Vector) {
	b.force = b.force.Add(force)
	// TODO: if applied at a specific point
	// this.torque += (point.x - position.x) * force.y - (point.y - position.y) * force.x;
}

// TODO - unused?
func (b *Body) ApplyTorque(torque float32) {
	b.torque += torque
}

func (b *Body) ClearForces() {
	b.force = math.Vector{}
	b.torque = 0
}

func (b *Body) ApplyImpulse(impulse math.Vector, contact math.Vector, negative bool) {
	if negative {
		impulse = impulse.Neg()
	}

	b.LinearVelocity = b.LinearVelocity.Add(impulse.Muls(b.inverseMass))
	b.AngularVelocity += b.inverseInertia * contact.Cross(impulse)
}

func (b *Body) NonorientedBound() (x1, y1, x2, y2 float32) {
	minX := float32(+stdmath.MaxFloat32)
	maxX := float32(-stdmath.MaxFloat32)
	minY := float32(+stdmath.MaxFloat32)
	maxY := float32(-stdmath.MaxFloat32)

	for _, fixture := range b.Fixtures {
		for _, vertex := range fixture.Vertices {
			v := vertex.Add(b.Position)

			minX = math.Min(minX, v.X)
			maxX = math.Max(maxX, v.X)
			minY = math.Min(minY, v.Y)
			maxY = math.Max(maxY, v.Y)
		}
	}

	return minX, minY, maxX, maxY
}

func (b *Body) CoverBound() (x1, y1, x2, y2 float32) {
	minX := float32(+stdmath.MaxFloat32)
	maxX := float32(-stdmath.MaxFloat32)
	minY := float32(+stdmath.MaxFloat32)
	maxY := float32(-stdmath.MaxFloat32)

	for _, fixture := range b.Fixtures {
		for i := range fixture.Vertices {
			v := fixture.VertexInWorldSpace(b, i)

			minX = math.Min(minX, v.X)
			maxX = math.Max(maxX, v.X)
			minY = math.Min(minY, v.Y)
			maxY = math.Max(maxY, v.Y)
		}
	}

	return minX, minY, maxX, maxY
}
