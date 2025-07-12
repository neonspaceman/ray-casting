package vec

import (
	"math"
)

type Vec2 struct {
	X, Y float32
}

func NewVec2(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

func NewRotated(angel float32) Vec2 {
	return Vec2{X: 1, Y: 0}.Rot(angel)
}

func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vec2) Mul(other Vec2) Vec2 {
	return Vec2{X: v.X * other.X, Y: v.Y * other.Y}
}

func (v Vec2) MulValue(value float32) Vec2 {
	return Vec2{X: v.X * value, Y: v.Y * value}
}

func (v Vec2) Div(other Vec2) Vec2 {
	return Vec2{X: v.X / other.X, Y: v.Y / other.Y}
}

func (v Vec2) DivValue(value float32) Vec2 {
	return Vec2{X: v.X / value, Y: v.Y / value}
}

func (v Vec2) Rad() float32 {
	return float32(math.Atan2(float64(v.Y), float64(v.X)))
}

func (v Vec2) Len() float32 {
	x, y := float64(v.X), float64(v.Y)
	return float32(math.Sqrt(x*x + y*y))
}

func (v Vec2) Len2() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

func (v Vec2) Norm() Vec2 {
	length := v.Len()
	return Vec2{X: v.X / length, Y: v.Y / length}
}

func (v Vec2) Rot(r float32) Vec2 {
	sin, cos := float32(math.Sin(float64(r))), float32(math.Cos(float64(r)))

	return Vec2{X: cos*v.X - sin*v.Y, Y: sin*v.X + cos*v.Y}
}
