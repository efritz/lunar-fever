package physics

import (
	stdmath "math"

	"github.com/efritz/lunar-fever/internal/common/math"
)

type Contact struct {
	fixture1        Fixture
	body1           *Body
	fixture2        Fixture
	body2           *Body
	contacts        []math.Vector
	penetration     float32
	normal          math.Vector
	restitution     float32
	staticFriction  float32
	dynamicFriction float32
}

type penetrationQueryResult struct {
	index    int
	distance float32
}

// TODO - attach body to fixtures

func NewContact(fixture1 Fixture, body1 *Body, fixture2 Fixture, body2 *Body) *Contact {
	penetration1 := queryFaceDirections(fixture1, body1, fixture2, body2)
	if penetration1.distance >= 0 {
		return nil
	}

	penetration2 := queryFaceDirections(fixture2, body2, fixture1, body1)
	if penetration2.distance >= 0 {
		return nil
	}

	if greaterThan(penetration1.distance, penetration2.distance) {
		return newContact(fixture1, body1, fixture2, body2, penetration1.index, false)
	} else {
		return newContact(fixture2, body2, fixture1, body1, penetration2.index, true)
	}
}

func greaterThan(a, b float32) bool {
	const (
		biasRelative = 0.95
		biasAbsolute = 0.01
	)

	// TODO - .... why?
	return a >= b*biasRelative+a*biasAbsolute
}

func newContact(ref Fixture, refBody *Body, inc Fixture, incBody *Body, refIndex int, flipped bool) *Contact {
	incIndex := findIncidentFace(inc, incBody, ref.NormalInWorldSpace(refBody, refIndex))

	refv1 := ref.VertexInWorldSpace(refBody, refIndex)
	refv2 := ref.VertexInWorldSpace(refBody, (refIndex+1)%len(ref.Vertices))

	incv1 := inc.VertexInWorldSpace(incBody, incIndex)
	incv2 := inc.VertexInWorldSpace(incBody, (incIndex+1)%len(inc.Vertices))

	refFaceNormal := ref.NormalInWorldSpace(refBody, refIndex)
	sideNormal1 := refFaceNormal.Orthogonalize()
	sideNormal2 := refFaceNormal.Orthogonalize().Neg()

	incv1, incv2, clip1 := clip(sideNormal2, refv1, incv1, incv2)
	incv1, incv2, clip2 := clip(sideNormal1, refv2, incv1, incv2)

	if !clip1 || !clip2 {
		return nil
	}

	penetration := float32(0)
	var contacts []math.Vector

	d1 := refFaceNormal.Dot(incv1) - refFaceNormal.Dot(refv1)
	d2 := refFaceNormal.Dot(incv2) - refFaceNormal.Dot(refv2)

	if d1 <= 0 {
		penetration += -d1
		contacts = append(contacts, incv1)
	}

	if d2 <= 0 {
		penetration += -d2
		contacts = append(contacts, incv2)
	}

	if len(contacts) == 0 {
		return nil
	}

	if !flipped {
		return &Contact{
			fixture1:        ref,
			body1:           refBody,
			fixture2:        inc,
			body2:           incBody,
			contacts:        contacts,
			penetration:     penetration / float32(len(contacts)),
			normal:          refFaceNormal,
			restitution:     math.Max(inc.restitution, ref.restitution),
			staticFriction:  math.Sqrt32(inc.staticFriction * ref.staticFriction),
			dynamicFriction: math.Sqrt32(inc.dynamicFriction * ref.dynamicFriction),
		}
	} else {
		return &Contact{
			fixture1:        inc,
			body1:           incBody,
			fixture2:        ref,
			body2:           refBody,
			contacts:        contacts,
			penetration:     penetration / float32(len(contacts)),
			normal:          refFaceNormal.Neg(),
			restitution:     math.Max(inc.restitution, ref.restitution),
			staticFriction:  math.Sqrt32(inc.staticFriction * ref.staticFriction),
			dynamicFriction: math.Sqrt32(inc.dynamicFriction * ref.dynamicFriction),
		}
	}
}

func queryFaceDirections(fixture1 Fixture, body1 *Body, fixture2 Fixture, body2 *Body) penetrationQueryResult {
	bestIndex := 0
	bestDistance := float32(-stdmath.MaxFloat32)

	for i := range fixture1.Vertices {
		face := fixture1.NormalInWorldSpace(body1, i)
		support := fixture2.VertexInWorldSpace(body2, getSupportPoint(fixture2, body2, face.Neg()))
		distance := face.Dot(support) - face.Dot(fixture1.VertexInWorldSpace(body1, i))

		if distance > bestDistance {
			bestIndex = i
			bestDistance = distance
		}
	}

	return penetrationQueryResult{
		index:    bestIndex,
		distance: bestDistance,
	}
}

func getSupportPoint(fixture Fixture, body *Body, direction math.Vector) int {
	bestIndex := 0
	bestDistance := float32(-stdmath.MaxFloat32)

	for i := range fixture.Vertices {
		distance := fixture.VertexInWorldSpace(body, i).Dot(direction)

		if distance > bestDistance {
			bestIndex = i
			bestDistance = distance
		}
	}

	return bestIndex
}

func findIncidentFace(fixture Fixture, body *Body, normal math.Vector) int {
	bestIndex := 0
	bestDistance := float32(stdmath.MaxFloat32)

	for i := range fixture.Vertices {
		distance := normal.Dot(fixture.NormalInWorldSpace(body, i))

		if distance < bestDistance {
			bestIndex = i
			bestDistance = distance
		}
	}

	return bestIndex
}

func clip(normal, point, v1, v2 math.Vector) (newV1, newV2 math.Vector, _ bool) {
	d1 := normal.Dot(v1) - normal.Dot(point)
	d2 := normal.Dot(v2) - normal.Dot(point)

	if d1 < 0 && d2 < 0 {
		return v1, v2, false
	}
	if d1 >= 0 && d2 >= 0 {
		return v1, v2, true
	}

	if d2 >= 0 {
		v1, v2 = v2, v1
		d1, d2 = d2, d1
	}

	v2 = v2.Sub(v1).Muls(d1 / (d1 - d2)).Add(v1)
	return v1, v2, true
}

//
//

const (
	penetrationAllowance  = float32(0.01)
	penetrationCorrection = float32(0.8)
)

func (c Contact) ApplyImpulse() {
	for _, contact := range c.contacts {
		c.applyImpulse(contact)
	}
}

func (c Contact) applyImpulse(contact math.Vector) {
	r1 := contact.Sub(c.body1.Position)
	r2 := contact.Sub(c.body2.Position)

	vn := getRelativeVelocity(c.body1, c.body2, r1, r2).Dot(c.normal)
	if vn > 0 {
		return
	}

	r1n := r1.Cross(c.normal)
	r2n := r2.Cross(c.normal)
	invMassSum := c.body1.inverseMass + c.body2.inverseMass + (r1n*r1n)*c.body1.inverseInertia + (r2n*r2n)*c.body2.inverseInertia

	magnitude1 := -vn * (c.restitution + 1)
	magnitude1 /= invMassSum
	magnitude1 /= float32(len(c.contacts))

	c.body1.ApplyImpulse(c.normal.Muls(magnitude1), r1, true)
	c.body2.ApplyImpulse(c.normal.Muls(magnitude1), r2, false)

	rv := getRelativeVelocity(c.body1, c.body2, r1, r2)
	tangent := rv.Add(c.normal.Muls(-rv.Dot(c.normal))).Normalize()

	magnitude2 := -rv.Dot(tangent)
	magnitude2 /= invMassSum
	magnitude2 /= float32(len(c.contacts))

	if math.Equal(magnitude2, 0) {
		return
	}

	if math.Abs32(magnitude2) >= magnitude1*c.staticFriction {
		magnitude2 = magnitude1 * -c.dynamicFriction
	}

	c.body1.ApplyImpulse(tangent.Muls(magnitude2), r1, true)
	c.body2.ApplyImpulse(tangent.Muls(magnitude2), r2, false)
}

func (c Contact) Correct() {
	if c.penetration < penetrationAllowance {
		return
	}

	correction := c.penetration / (c.body1.inverseMass + c.body2.inverseMass) * penetrationCorrection
	c.body1.Position = c.body1.Position.Add(c.normal.Muls(-c.body1.inverseMass * correction))
	c.body2.Position = c.body2.Position.Add(c.normal.Muls(+c.body2.inverseMass * correction))
}

func getRelativeVelocity(body1, body2 *Body, r1, r2 math.Vector) math.Vector {
	r1c := r1.Crosss(body1.AngularVelocity)
	r2c := r2.Crosss(body2.AngularVelocity)

	return body2.LinearVelocity.Add(r2c).Sub(body1.LinearVelocity).Sub(r1c)
}
