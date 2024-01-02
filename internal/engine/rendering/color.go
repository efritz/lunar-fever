package rendering

type Color struct {
	R float32
	G float32
	B float32
	A float32
}

var (
	White = Color{1, 1, 1, 1}
	Black = Color{0, 0, 0, 1}
)
