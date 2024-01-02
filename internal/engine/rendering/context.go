package rendering

import (
	"github.com/efritz/lunar-fever/assets"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Context struct {
	Window        *glfw.Window
	ShaderProgram *ShaderProgram
	TextureLoader *TextureLoader
	SpriteBatch   *SpriteBatch
}

func InitContext() (Context, error) {
	window, err := initGlfw()
	if err != nil {
		return Context{}, err
	}

	if err := initOpenGL(); err != nil {
		return Context{}, err
	}

	// TODO - bundle more closely to this package
	vertexShaderSource, err := assets.LoadShader("default_vert")
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := assets.LoadShader("default_frag")
	if err != nil {
		panic(err)
	}
	shaderProgram := NewShaderProgram(vertexShaderSource, fragmentShaderSource)

	return Context{
		Window:        window,
		ShaderProgram: shaderProgram,
		TextureLoader: NewTextureLoader(),
		SpriteBatch:   NewSpriteBatch(shaderProgram),
	}, nil
}

func initGlfw() (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(DisplayWidth, DisplayHeight, WindowTitle, nil, nil)
	if err != nil {
		return nil, err
	}

	window.MakeContextCurrent()
	return window, nil
}

func initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.TEXTURE_2D)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, DisplayWidth, 0, DisplayHeight, 1, -1)
	gl.MatrixMode(gl.MODELVIEW)
	return nil
}
