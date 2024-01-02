package rendering

import (
	"github.com/efritz/lunar-fever/assets"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/nullboundary/glfont"
)

type Font struct {
	font *glfont.Font
}

func LoadFont(name string) (*Font, error) {
	fontBytes, err := assets.LoadRawFont(name)
	if err != nil {
		return nil, err
	}

	font, err := glfont.LoadFontBytes(fontBytes, 50, DisplayWidth, DisplayHeight)
	if err != nil {
		return nil, err
	}

	return &Font{
		font: font,
	}, nil
}

func (f Font) Printf(x, y, scale float32, format string, args ...interface{}) {
	f.font.SetColor(1.0, 1.0, 0.0, 1.0) // TODO - with options
	f.font.Printf(x, y, scale, format, args...)
	gl.Enable(gl.BLEND)
}
