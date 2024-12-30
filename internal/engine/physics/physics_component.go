package physics

type PhysicsComponent struct {
	Body               *Body
	CollisionsDisabled bool
}

type PhysicsComponentType struct{}

var physicsComponentType = PhysicsComponentType{}

func (c *PhysicsComponent) ComponentType() PhysicsComponentType {
	return physicsComponentType
}
