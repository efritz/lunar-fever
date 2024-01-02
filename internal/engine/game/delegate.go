package game

type Delegate interface {
	Init()
	Exit()
	Update(elapsedMs int64)
	Render(elapsedMs int64)
	Active() bool
}
