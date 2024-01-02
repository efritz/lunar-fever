package rendering

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

func compileShader(shaderType uint32, source string) uint32 {
	csources, free := gl.Strs(source)
	defer free()

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, csources, nil)
	gl.CompileShader(shader)

	if log, ok := glTestSuccess(shader, gl.COMPILE_STATUS, glShaderLogFns); !ok {
		panic(fmt.Errorf("failed to compile shader: %s", log))
	}

	return shader
}

type glProgramShaderLogFns struct {
	getiv      func(id uint32, pname uint32, params *int32)
	getInfoLog func(id uint32, bufSize int32, length *int32, infoLog *uint8)
}

var (
	glProgramLogFns = glProgramShaderLogFns{gl.GetProgramiv, gl.GetProgramInfoLog}
	glShaderLogFns  = glProgramShaderLogFns{gl.GetShaderiv, gl.GetShaderInfoLog}
)

func glTestSuccess(id uint32, pname uint32, fns glProgramShaderLogFns) (string, bool) {
	var status int32
	fns.getiv(id, pname, &status)

	if status != gl.FALSE {
		return "", true
	}

	var logLength int32
	fns.getiv(id, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength)+1)
	fns.getInfoLog(id, logLength, nil, gl.Str(log))

	return log, false
}
