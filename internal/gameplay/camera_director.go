package gameplay

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine"
)

type CameraDirector struct {
	*engine.Context

	startX            float32
	startY            float32
	targetX           float32
	targetY           float32
	currentMagnitude  float32
	startingMagnitude float32
	time              int64
	totalTime         int64
}

func (d *CameraDirector) Init() {}
func (d *CameraDirector) Exit() {}

func (d *CameraDirector) LookAt(x, y float32, timeMs int64) {
	x1, y1, x2, y2 := d.Camera.Bounds()
	d.startX = x1 + (x2-x1)/2
	d.startY = y1 + (y2-y1)/2

	d.targetX = x
	d.targetY = y

	d.time = timeMs
	d.totalTime = timeMs
}

func (d *CameraDirector) AddShake(magnitude float32) {
	if magnitude > d.currentMagnitude {
		d.currentMagnitude = magnitude
		d.startingMagnitude = magnitude
	}
}

func (d *CameraDirector) Process(elapsedMs int64) {
	if d.time > 0 {
		x := ease(d.startX, d.targetX, 1-float32(d.time)/float32(d.totalTime))
		y := ease(d.startY, d.targetY, 1-float32(d.time)/float32(d.totalTime))

		d.time -= elapsedMs
		d.Camera.LookAt(x, y)
	}

	if d.currentMagnitude > 0 {
		min := float32(+3*elapsedMs) * d.currentMagnitude / d.startingMagnitude
		max := float32(-3*elapsedMs) * d.currentMagnitude / d.startingMagnitude

		d.currentMagnitude -= float32(elapsedMs)
		d.Camera.Translate(math.Random(min, max), math.Random(min, max))
	}
}

func ease(min, max, percent float32) float32 {
	return (1-math.Pow32(1-percent, 3))*(max-min) + min
}
