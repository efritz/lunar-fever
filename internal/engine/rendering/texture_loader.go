package rendering

import (
	"github.com/efritz/lunar-fever/assets"
	"github.com/go-gl/gl/all-core/gl"
)

type TextureLoader struct {
	cache map[string]Texture
}

func NewTextureLoader() *TextureLoader {
	return &TextureLoader{
		cache: map[string]Texture{},
	}
}

func (l *TextureLoader) Load(name string) Texture {
	if texture, ok := l.cache[name]; ok {
		return texture
	}

	texture := loadTexture(name)
	l.cache[name] = texture
	return texture
}

var texParameters = map[uint32]int32{
	gl.TEXTURE_MIN_FILTER: gl.NEAREST,
	gl.TEXTURE_MAG_FILTER: gl.NEAREST,
}

func loadTexture(name string) Texture {
	rgba, err := assets.LoadTexture(name)
	if err != nil {
		panic(err)
	}

	var (
		bounds = rgba.Bounds()
		width  = bounds.Dx()
		height = bounds.Dy()
	)

	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)

	gl.TexImage2D(
		gl.TEXTURE_2D,    // target
		0,                // level
		gl.RGBA,          // internal format
		int32(width),     // width
		int32(height),    // height
		0,                // border
		gl.RGBA,          //format
		gl.UNSIGNED_BYTE, // xtype; uint8
		gl.Ptr(rgba.Pix), // pixels
	)

	for k, v := range texParameters {
		gl.TexParameteri(gl.TEXTURE_2D, k, v)
	}

	return NewTexture(id, float32(width), float32(height))
}
