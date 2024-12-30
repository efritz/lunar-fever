package gameplay

type InteractionComponent struct {
	Interacting   bool
	CooldownTimer int64
}

type InteractionComponentType struct{}

var interactionComponentType = InteractionComponentType{}

func (c *InteractionComponent) ComponentType() InteractionComponentType {
	return interactionComponentType
}
