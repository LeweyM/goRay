package Object

import (
	"fmt"
	"goRay/Ray"
	"goRay/Vector"
	"math"
)

type Sphere struct {
	center Vector.Vector
	radius int
}

func NewSphere(center Vector.Vector, radius int) *Sphere {
	return &Sphere{
		center: center,
		radius: radius,
	}
}

func (s *Sphere) String() string {
	return fmt.Sprintf("{postion: %s, radius: %b}", s.center, s.radius)
}

func (s *Sphere) IntersectDistance(r Ray.Ray) (bool, float64) {
	oc := r.Origin().Minus(s.center)

	a := r.Direction().Dot(r.Direction())
	b := 2 * oc.Dot(r.Direction())
	c := oc.Dot(oc) - float64(s.radius*s.radius)
	discriminant := b*b - (4 * a * c)

	if discriminant < 0 {
		return false, float64(0) //todo
	} else {
		t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
		//t2 := 0 - b + math.Sqrt(discriminant)/(2*a)
		return t1 >= 0, t1 //todo
	}
}

func (s *Sphere) Radius() int {
	return s.radius
}
