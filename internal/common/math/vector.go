package math

type Vector struct {
	X float32
	Y float32
}

func (v Vector) Len() float32          { return Sqrt32(v.X*v.X + v.Y*v.Y) }
func (v Vector) Neg() Vector           { return Vector{-v.X, -v.Y} }
func (v Vector) Orthogonalize() Vector { return Vector{+v.Y, -v.X} }

func (v Vector) Normalize() Vector {
	if len := v.Len(); len != 0 {
		return v.Divs(len)
	}

	return Vector{}
}

func (v Vector) Dot(other Vector) float32   { return v.X*other.X + v.Y*other.Y }
func (v Vector) Cross(other Vector) float32 { return v.X*other.Y - v.Y*other.X }
func (v Vector) Add(other Vector) Vector    { return Vector{v.X + other.X, v.Y + other.Y} }
func (v Vector) Sub(other Vector) Vector    { return Vector{v.X - other.X, v.Y - other.Y} }
func (v Vector) Mul(other Vector) Vector    { return Vector{v.X * other.X, v.Y * other.Y} }

func (v Vector) Crosss(n float32) Vector { return Vector{v.Y * -n, v.X * +n} }
func (v Vector) Adds(n float32) Vector   { return Vector{v.X + n, v.Y + n} }
func (v Vector) Subs(n float32) Vector   { return Vector{v.X - n, v.Y - n} }
func (v Vector) Muls(n float32) Vector   { return Vector{v.X * n, v.Y * n} }
func (v Vector) Divs(n float32) Vector   { return Vector{v.X / n, v.Y / n} }
