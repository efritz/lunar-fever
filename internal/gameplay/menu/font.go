package menu

import "github.com/efritz/lunar-fever/internal/engine/rendering"

var font *rendering.Font

func initFonts() {
	var err error
	if font, err = rendering.LoadFont("Roboto-Light"); err != nil {
		panic(err)
	}
}
