package camera

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/efritz/lunar-fever/internal/engine/rendering"
)

type Camera struct {
	x float32
	y float32
	z float32
	w float32
	h float32
}

func NewCamera() *Camera {
	return &Camera{
		z: 1,
		w: rendering.DisplayWidth,
		h: rendering.DisplayHeight,
	}
}

func (c *Camera) Translate(x, y float32) {
	c.x += x
	c.y += y
}

func (c *Camera) LookAt(x, y float32) {
	c.x = -(x - c.w/2)
	c.y = -(y - c.h/2)
}

func (c *Camera) Zoom(z float32) {
	c.z, _ = math.Clamp(c.z+z, 0.25, 1.25)
}

func (c *Camera) Bounds() (x1, y1, x2, y2 float32) {
	x := -c.x
	y := -c.y

	if c.z != 1 {
		x += c.w/2 - c.w/2/c.z
		y += c.h/2 - c.h/2/c.z
	}

	return x, y, x + c.w/c.z, y + c.h/c.z
}

func (c *Camera) ViewMatrix() math.Matrix4f32 {
	return math.IdentityMatrix.Translate(c.w/2, c.h/2).Scale(c.z, c.z, 1).Translate(c.x-(c.w/2), c.y-(c.h/2))
}

func (c *Camera) Unprojectx(value float32) float32 { return c.unproject(value, c.x, c.w) }
func (c *Camera) UnprojectY(value float32) float32 { return c.unproject(value, c.y, c.h) }

func (c *Camera) unproject(value, pos, size float32) float32 {
	return value/c.z - ((size + 2*pos*c.z - c.z*size) / (2 * c.z))
}
