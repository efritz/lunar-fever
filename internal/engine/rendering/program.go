package rendering

import (
	"fmt"

	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/go-gl/gl/all-core/gl"
)

type ShaderProgram struct {
	id               uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
}

func NewShaderProgram(vertexShaderSource, fragmentShaderSource string) *ShaderProgram {
	program := gl.CreateProgram()
	vertexShader := compileShader(gl.VERTEX_SHADER, vertexShaderSource)
	fragmentShader := compileShader(gl.FRAGMENT_SHADER, fragmentShaderSource)

	for _, attrs := range vertexAttributes {
		gl.BindAttribLocation(program, attrs.location, gl.Str(attrs.name))
	}

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	if log, ok := glTestSuccess(program, gl.LINK_STATUS, glProgramLogFns); !ok {
		panic(fmt.Errorf("failed to link program: %s", log))
	}

	return &ShaderProgram{
		id:               program,
		vertexShaderID:   vertexShader,
		fragmentShaderID: fragmentShader,
	}
}

func (p *ShaderProgram) Use() {
	gl.UseProgram(p.id)
}

const uTexture = "u_texture" + null

func (p *ShaderProgram) SetTexture() {
	uTextureLocation := gl.GetUniformLocation(p.id, gl.Str(uTexture))
	if uTextureLocation == -1 {
		panic(fmt.Errorf("unknown uniform variable %q", uTexture))
	}

	gl.Uniform1i(uTextureLocation, 0)
}

const uProjView = "u_projView" + null

func (p *ShaderProgram) SetProjectionMatrix(m math.Matrix4f32) {
	uProjViewLocation := gl.GetUniformLocation(p.id, gl.Str(uProjView))
	if uProjViewLocation == -1 {
		panic(fmt.Errorf("unknown uniform variable %q", uProjView))
	}

	buffer := m.IntoBuffer(make([]float32, 0, 16))
	gl.UniformMatrix4fv(uProjViewLocation, 1, false, &buffer[0])
}
