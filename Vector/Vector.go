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

func (v Vector) Minus(vector Vector) Vector {
	return Vector{
		x: v.x - vector.x,
		y: v.y - vector.y,
		z: v.z - vector.z,
	}
}

func New(x float64, y float64, z float64) *Vector {
	return &Vector{x: x, y: y, z: z}
}


