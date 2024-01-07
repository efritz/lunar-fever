package rendering

import (
	"github.com/efritz/lunar-fever/internal/common/math"
	"github.com/go-gl/gl/all-core/gl"
)

const (
	SpriteBatchSize    = 1024
	numVerticesPerRect = int32(6) // drawn via two triangles
	sizeOfFloat        = int32(4) // int32 is 4 bytes
)

var (
	numValuesPerVertex int32 = func() (total int32) {
		for _, attr := range vertexAttributes {
			total += attr.numValues
		}

		return total
	}()

	renderStride = int32(numValuesPerVertex) * sizeOfFloat
)

type SpriteBatch struct {
	program          *ShaderProgram
	viewMatrix       math.Matrix4f32
	projectionMatrix math.Matrix4f32
	drawing          bool
	vertices         []float32
	textureID        uint32
}

func NewSpriteBatch(program *ShaderProgram) *SpriteBatch {
	n := SpriteBatchSize * numVerticesPerRect * numValuesPerVertex
	projectionMatrix := math.OrthoMatrix(0, DisplayWidth, 0, DisplayHeight, 1, -1)

	return &SpriteBatch{
		program:          program,
		viewMatrix:       math.IdentityMatrix,
		projectionMatrix: projectionMatrix,
		vertices:         make([]float32, 0, n),
	}
}

func (sb *SpriteBatch) combinedMatrix() math.Matrix4f32 {
	return sb.projectionMatrix.Transpose().Multiply(sb.viewMatrix)
}

func (sb *SpriteBatch) SetViewMatrix(m math.Matrix4f32) {
	if !sb.drawing {
		sb.viewMatrix = m
		return
	}

	sb.flush()
	sb.viewMatrix = m
	sb.program.SetProjectionMatrix(sb.combinedMatrix())
}

func (sb *SpriteBatch) Begin() {
	if sb.drawing {
		panic("already drawing")
	}

	sb.program.Use()
	sb.program.SetTexture()
	sb.program.SetProjectionMatrix(sb.combinedMatrix())
	sb.drawing = true
}

func (sb *SpriteBatch) End() {
	if !sb.drawing {
		panic("not drawing")
	}

	sb.drawing = false
	sb.flush()
}

type DrawOptions struct {
	Color         Color
	Rotation      float32
	OriginX       float32
	OriginY       float32
	ScaleX        float32
	ScaleY        float32
	SpriteEffects SpriteEffects
}

type SpriteEffects uint8

const (
	SpriteEffectFlipHorizontal SpriteEffects = 1 << iota
	SpriteEffectFlipVertical
)

type DrawOptionFunc func(o *DrawOptions)

func WithColor(color Color) DrawOptionFunc {
	return func(o *DrawOptions) { o.Color = color }
}

func WithRotation(rotation float32) DrawOptionFunc {
	return func(o *DrawOptions) { o.Rotation = rotation }
}

func WithOrigin(originX, originY float32) DrawOptionFunc {
	return func(o *DrawOptions) { o.OriginX = originX; o.OriginY = originY }
}

func WithScale(scaleX, scaleY float32) DrawOptionFunc {
	return func(o *DrawOptions) { o.ScaleX = scaleX; o.ScaleY = scaleY }
}

func WithSpriteEffects(spriteEffects SpriteEffects) DrawOptionFunc {
	return func(o *DrawOptions) { o.SpriteEffects = spriteEffects }
}

func (sb *SpriteBatch) Draw(texture Texture, x, y, w, h float32, optionFns ...DrawOptionFunc) {
	options := DrawOptions{
		Color:  White,
		ScaleX: 1,
		ScaleY: 1,
	}

	for _, fn := range optionFns {
		fn(&options)
	}

	sb.DrawWithOptions(texture, x, y, w, h, options)
}

func (sb *SpriteBatch) DrawWithOptions(texture Texture, x, y, w, h float32, options DrawOptions) {
	if !sb.drawing {
		panic("not drawing")
	}

	// Flush if we're changing textures or buffer is full
	if texture.ID != sb.textureID || len(sb.vertices) == cap(sb.vertices) {
		sb.flush()
	}

	var (
		x1, y1 float32
		x2, y2 float32
		x3, y3 float32
		x4, y4 float32
		u1, v1 = texture.U1, texture.V1
		u2, v2 = texture.U2, texture.V2
	)

	fx1, fy1 := 0-options.OriginX, 0-options.OriginY
	fx2, fy2 := w-options.OriginX, h-options.OriginY

	if options.ScaleX != 1 || options.ScaleY != 1 {
		fx1 *= options.ScaleX
		fx2 *= options.ScaleX
		fy1 *= options.ScaleY
		fy2 *= options.ScaleY
	}

	if options.Rotation != 0 {
		sin := math.Sin32(options.Rotation)
		cos := math.Cos32(options.Rotation)

		applyRotation := func(x, y float32) (float32, float32) {
			return (cos*x - sin*y), (sin*x + cos*y)
		}

		x1, y1 = applyRotation(fx1, fy1)
		x2, y2 = applyRotation(fx2, fy1)
		x3, y3 = applyRotation(fx2, fy2)
		x4, y4 = applyRotation(fx1, fy2)
	} else {
		x1, y1 = fx1, fy1
		x2, y2 = fx2, fy1
		x3, y3 = fx2, fy2
		x4, y4 = fx1, fy2
	}

	x1, y1 = x1+x+options.OriginX, y1+y+options.OriginY
	x2, y2 = x2+x+options.OriginX, y2+y+options.OriginY
	x3, y3 = x3+x+options.OriginX, y3+y+options.OriginY
	x4, y4 = x4+x+options.OriginX, y4+y+options.OriginY

	if options.SpriteEffects&SpriteEffectFlipHorizontal != 0 {
		u1, u2 = u2, u1
	}
	if options.SpriteEffects&SpriteEffectFlipVertical != 0 {
		v1, v2 = v2, v1
	}

	sb.vertices = append(sb.vertices,
		x1, y1, u1, v1, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
		x2, y2, u2, v1, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
		x4, y4, u1, v2, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
		x2, y2, u2, v1, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
		x3, y3, u2, v2, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
		x4, y4, u1, v2, options.Color.R, options.Color.G, options.Color.B, options.Color.A,
	)

	// Keep reference to texture for future flush
	sb.textureID = texture.ID
}

func (sb *SpriteBatch) flush() {
	if len(sb.vertices) == 0 {
		return
	}

	sb.render()
	sb.vertices = sb.vertices[:0]
}

func (sb *SpriteBatch) render() {
	if sb.textureID == 0 {
		panic("texture is not set")
	}

	var (
		numVertices  = len(sb.vertices)
		bufferSize   = numVertices * int(sizeOfFloat)
		numTriangles = int32(numVertices) / numValuesPerVertex
	)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, bufferSize, gl.Ptr(sb.vertices), gl.STATIC_DRAW)

	offset := uintptr(0)
	for _, attr := range vertexAttributes {
		gl.EnableVertexAttribArray(attr.location)
		gl.VertexAttribPointerWithOffset(attr.location, int32(attr.numValues), gl.FLOAT, false, renderStride, offset)
		offset += uintptr(attr.numValues * sizeOfFloat)
	}

	gl.BindTexture(gl.TEXTURE_2D, sb.textureID)
	gl.DrawArrays(gl.TRIANGLES, 0, numTriangles)
}
