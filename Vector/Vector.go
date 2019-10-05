package Vector

import "fmt"

type Vector struct {
	x float64
	y float64
	z float64
}

func (v Vector) Z() float64 {
	return v.z
}

func (v Vector) X() float64 {
	return v.x
}

func (v Vector) Y() float64 {
	return v.y
}

func (v Vector) String() string {
	return fmt.Sprintf("x:%f, y:%f, z%f", v.x, v.y, v.z)
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

func New(x float64, y float64, z float64) *Vector {
	return &Vector{x: x, y: y, z: z}
}


