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

type TextOptions struct {
	Color Color
	Scale float32
}

type TextOptionFunc func(o *TextOptions)

func WithTextColor(color Color) TextOptionFunc   { return func(o *TextOptions) { o.Color = color } }
func WithTextScale(scale float32) TextOptionFunc { return func(o *TextOptions) { o.Scale = scale } }

func (f Font) Printf(x, y float32, text string, optionFns ...TextOptionFunc) {
	options := TextOptions{
		Color: Black,
		Scale: 1,
	}

	for _, fn := range optionFns {
		fn(&options)
	}

	f.font.SetColor(options.Color.R, options.Color.G, options.Color.B, options.Color.A)
	f.font.Printf(x, y, options.Scale, text)
	gl.Enable(gl.BLEND)
}
