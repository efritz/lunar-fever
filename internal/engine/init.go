package engine

import "runtime"

func init() {
	runtime.LockOSThread()
}
