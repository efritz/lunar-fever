package math

type Matrix32 struct {
	M00, M01, M10, M11 float32
}

func (m Matrix32) Mul(v Vector) Vector {
	return Vector{
		X: m.M00*v.X + m.M01*v.Y,
		Y: m.M10*v.X + m.M11*v.Y,
	}
}

//
//

type Matrix4f32 struct {
	M00, M01, M02, M03 float32
	M10, M11, M12, M13 float32
	M20, M21, M22, M23 float32
	M30, M31, M32, M33 float32
}

var IdentityMatrix = Matrix4f32{
	M00: 1,
	M11: 1,
	M22: 1,
	M33: 1,
}

func OrthoMatrix(x1, x2, y1, y2, z2, z1 float32) Matrix4f32 {
	return Matrix4f32{
		M00: 2 / (x2 - x1),
		M11: 2 / (y1 - y2),
		M22: 2 / (z1 - z2) * -1,
		M03: -(x2 + x1) / (x2 - x1),
		M13: -(y1 + y2) / (y1 - y2),
		M23: -(z1 + z2) / (z1 - z2),
		M33: 1,
	}
}

func (m Matrix4f32) Transpose() Matrix4f32 {
	return Matrix4f32{
		M00: m.M00, M01: m.M10, M02: m.M20, M03: m.M30,
		M10: m.M01, M11: m.M11, M12: m.M21, M13: m.M31,
		M20: m.M02, M21: m.M12, M22: m.M22, M23: m.M32,
		M30: m.M03, M31: m.M13, M32: m.M23, M33: m.M33,
	}
}

func (m Matrix4f32) IntoBuffer(buf []float32) []float32 {
	return append(buf,
		m.M00, m.M01, m.M02, m.M03,
		m.M10, m.M11, m.M12, m.M13,
		m.M20, m.M21, m.M22, m.M23,
		m.M30, m.M31, m.M32, m.M33,
	)
}

func (left Matrix4f32) Multiply(right Matrix4f32) Matrix4f32 {
	return Matrix4f32{
		M00: left.M00*right.M00 + left.M10*right.M01 + left.M20*right.M02 + left.M30*right.M03,
		M01: left.M01*right.M00 + left.M11*right.M01 + left.M21*right.M02 + left.M31*right.M03,
		M02: left.M02*right.M00 + left.M12*right.M01 + left.M22*right.M02 + left.M32*right.M03,
		M03: left.M03*right.M00 + left.M13*right.M01 + left.M23*right.M02 + left.M33*right.M03,
		M10: left.M00*right.M10 + left.M10*right.M11 + left.M20*right.M12 + left.M30*right.M13,
		M11: left.M01*right.M10 + left.M11*right.M11 + left.M21*right.M12 + left.M31*right.M13,
		M12: left.M02*right.M10 + left.M12*right.M11 + left.M22*right.M12 + left.M32*right.M13,
		M13: left.M03*right.M10 + left.M13*right.M11 + left.M23*right.M12 + left.M33*right.M13,
		M20: left.M00*right.M20 + left.M10*right.M21 + left.M20*right.M22 + left.M30*right.M23,
		M21: left.M01*right.M20 + left.M11*right.M21 + left.M21*right.M22 + left.M31*right.M23,
		M22: left.M02*right.M20 + left.M12*right.M21 + left.M22*right.M22 + left.M32*right.M23,
		M23: left.M03*right.M20 + left.M13*right.M21 + left.M23*right.M22 + left.M33*right.M23,
		M30: left.M00*right.M30 + left.M10*right.M31 + left.M20*right.M32 + left.M30*right.M33,
		M31: left.M01*right.M30 + left.M11*right.M31 + left.M21*right.M32 + left.M31*right.M33,
		M32: left.M02*right.M30 + left.M12*right.M31 + left.M22*right.M32 + left.M32*right.M33,
		M33: left.M03*right.M30 + left.M13*right.M31 + left.M23*right.M32 + left.M33*right.M33,
	}
}

func (m Matrix4f32) Translate(x, y float32) Matrix4f32 {
	m.M30 += m.M00*x + m.M10*y
	m.M31 += m.M01*x + m.M11*y
	m.M32 += m.M02*x + m.M12*y
	m.M33 += m.M03*x + m.M13*y
	return m
}

func (m Matrix4f32) Scale(x, y, z float32) Matrix4f32 {
	m.M00 = m.M00 * x
	m.M01 = m.M01 * x
	m.M02 = m.M02 * x
	m.M03 = m.M03 * x
	m.M10 = m.M10 * y
	m.M11 = m.M11 * y
	m.M12 = m.M12 * y
	m.M13 = m.M13 * y
	m.M20 = m.M20 * z
	m.M21 = m.M21 * z
	m.M22 = m.M22 * z
	m.M23 = m.M23 * z
	return m
}
