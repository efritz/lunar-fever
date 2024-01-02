package rendering

type Texture struct {
	ID             uint32
	Width          float32
	Height         float32
	U1, V1, U2, V2 float32
}

func NewTexture(id uint32, width, height float32) Texture {
	return Texture{
		ID:     id,
		Width:  width,
		Height: height,
		U1:     0,
		V1:     0,
		U2:     1,
		V2:     1,
	}
}

func (t Texture) Region(x, y, width, height float32) Texture {
	return Texture{
		ID:     t.ID,
		Width:  t.Width,
		Height: t.Height,
		U1:     x / t.Width,
		V1:     y / t.Height,
		U2:     (x + width) / t.Width,
		V2:     (y + height) / t.Height,
	}
}
