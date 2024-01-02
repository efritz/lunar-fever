package math

import (
	"math"
	"math/rand"
)

func Random(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func Abs32(v float32) float32 {
	return float32(math.Abs(float64(v)))
}

func Sqrt32(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

func Pow32(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Sin32(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func Cos32(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func Atan232(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

const epsilon = float32(0.0001)

func Equal(a, b float32) bool {
	return Abs32(a-b) <= epsilon
}

func Sign(v float32) int {
	if v == 0 {
		return 0
	}

	if v < 0 {
		return -1
	}

	return +1
}

func Min[E int | int32 | int64 | float32 | float64](a, b E) E {
	if a < b {
		return a
	}

	return b
}

func Max[E int | int32 | int64 | float32 | float64](a, b E) E {
	if a < b {
		return b
	}

	return a
}

func Clamp[E int32 | int64 | float32 | float64](v, min, max E) (_ E, clamped bool) {
	if v <= min {
		return min, true
	}
	if v >= max {
		return max, true
	}

	return v, false
}

func NextMultiple(v int, scale int) int {
	if v < 0 {
		scale = -scale
	}

	if remainder := v % scale; remainder != 0 {
		return v + scale - remainder
	}

	return v
}

func PrevMultiple(v int, scale int) int {
	if v >= scale {
		return v - (v % scale)
	}

	if v < 0 {
		return ((-v + scale) / scale) * scale * -1
	}

	return 0
}
