package Vector

import (
	"fmt"
	"math"
)

type Vector struct {
	x float64
	y float64
	z float64
}

func (v Vector) String() string {
	return fmt.Sprintf("x:%f, y:%f, z%f", v.x, v.y, v.z)
}

func New(x float64, y float64, z float64) *Vector {

	return &Vector{x: x, y: y, z: z}
}

func (v *Vector) X() float64 {
	return v.x
}

func (v *Vector) Y() float64 {
	return v.y
}

func (v *Vector) Z() float64 {
	return v.z
}

func (v Vector) Minus(v2 Vector) Vector {
	return Vector{
		x: v.x - v2.x,
		y: v.y - v2.y,
		z: v.z - v2.z,
	}
}

func (v Vector) Dot(v2 Vector) float64 {
	return (v.x * v2.x) + (v.y * v2.y) + (v.z * v2.z)
}

// Deprecated due to performance
func (v Vector) RotateY(radAngle float64) Vector {
	x := (v.X() * math.Cos(radAngle)) + (v.Z() * math.Sin(radAngle))
	z := (-v.X() * math.Sin(radAngle)) + (v.Z() * math.Cos(radAngle))
	return *New(x, v.Y(), z)
}

func RotateYBuilder(radAngle float64) func(v Vector) Vector {
	cos := math.Cos(radAngle)
	sin := math.Sin(radAngle)
	return func(v Vector) Vector {
		x := (v.X() * cos) + (v.Z() * sin)
		z := (-v.X() * sin) + (v.Z() * cos)
		return *New(x, v.Y(), z)
	}
}

func (v Vector) Normalize() Vector {
	magnitude := v.magnitude()
	if magnitude != 0 {
		v.x = v.x / magnitude
		v.y = v.y / magnitude
		v.z = v.z / magnitude
	}
	return v
}

func (v Vector) magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vector) TranslateZ(f float64) Vector {
	return *New(v.x, v.y, v.z + f)
}

func (v Vector) TranslateX(f float64) Vector {
	return *New(v.x + f, v.y, v.z)
}

func (v Vector) Translate(vector Vector) Vector {
	return *New(v.x + vector.x, v.y + vector.y, v.z + vector.z)
}

func (v Vector) Scale(f float64) Vector {
	return *New(v.x * f, v.y * f, v.z * f)
}

func (v Vector) DistanceBetween(vector Vector) float64 {
	return v.Minus(vector).magnitude()
}


