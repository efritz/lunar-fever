package gameplay

type HealthComponent struct {
	Health    float32
	MaxHealth float32
}

type HealthComponentType struct{}

var healthComponentType = HealthComponentType{}

func (c *HealthComponent) ComponentType() HealthComponentType {
	return healthComponentType
}
