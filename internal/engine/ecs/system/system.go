package system

type System interface {
	Init()
	Exit()
	Process(elapsedMs int64)
}
