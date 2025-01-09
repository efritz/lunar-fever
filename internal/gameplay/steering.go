package gameplay

import "github.com/efritz/lunar-fever/internal/common/math"

type SteeringBehavior interface {
	Calculate() math.Vector
}

// WeightedBehavior pairs a behavior with a weight.
type WeightedBehavior struct {
	Behavior SteeringBehavior
	Weight   float32
}

// SteeringManager combines multiple behaviors.
type SteeringManager struct {
	Behaviors []WeightedBehavior
	MaxForce  float32
}

func NewSteeringManager(maxForce float32) *SteeringManager {
	return &SteeringManager{
		Behaviors: []WeightedBehavior{},
		MaxForce:  maxForce,
	}
}

// AddBehavior appends a new weighted behavior.
func (sm *SteeringManager) AddBehavior(b SteeringBehavior, weight float32) {
	sm.Behaviors = append(sm.Behaviors, WeightedBehavior{
		Behavior: b,
		Weight:   weight,
	})
}

// Calculate produces the final steering vector from all behaviors.
func (sm *SteeringManager) Calculate() math.Vector {
	var steering math.Vector

	for _, wb := range sm.Behaviors {
		steer := wb.Behavior.Calculate().Muls(wb.Weight)
		steering = steering.Add(steer)
	}

	// Optionally, clamp the steering magnitude to MaxForce
	if sm.MaxForce > 0 && steering.Len() > sm.MaxForce {
		steering = steering.Normalize().Muls(sm.MaxForce)
	}

	return steering
}

// Seek tries to move the entity toward a target at a given desired speed.
type Seek struct {
	Position     *math.Vector // Current position of the entity
	Velocity     *math.Vector // Current velocity of the entity
	Target       math.Vector  // Target to seek
	DesiredSpeed float32
}

func (s *Seek) Calculate() math.Vector {
	desired := s.Target.Sub(*s.Position)
	desired = desired.Normalize().Muls(s.DesiredSpeed) // desired velocity
	steering := desired.Sub(*s.Velocity)               // steering = desired - current
	return steering
}

// Flee tries to move the entity away from a threat at a given desired speed.
type Flee struct {
	Position     *math.Vector
	Velocity     *math.Vector
	Threat       math.Vector
	DesiredSpeed float32
}

func (f *Flee) Calculate() math.Vector {
	desired := (*f.Position).Sub(f.Threat)
	desired = desired.Normalize().Muls(f.DesiredSpeed)
	steering := desired.Sub(*f.Velocity)
	return steering
}

// Arrival slows down as the entity approaches its target within a certain radius.
type Arrival struct {
	Position      *math.Vector
	Velocity      *math.Vector
	Target        math.Vector
	DesiredSpeed  float32
	SlowingRadius float32
}

func (a *Arrival) Calculate() math.Vector {
	toTarget := a.Target.Sub(*a.Position)
	distance := toTarget.Len()
	// If we’re far away, proceed at full speed.
	// If we’re within the SlowingRadius, scale the speed down proportionally.
	speed := a.DesiredSpeed
	if distance < a.SlowingRadius {
		speed = speed * (distance / a.SlowingRadius)
	}

	desired := toTarget.Normalize().Muls(speed)
	steering := desired.Sub(*a.Velocity)
	return steering
}
