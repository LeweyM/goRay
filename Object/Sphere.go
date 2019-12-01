package Object

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"goRay/Ray"
	"goRay/Vector"
	"math"
)

type Sphere struct {
	center Vector.Vector
	radius int
}

func (s *Sphere) Draw(renderer *sdl.Renderer, xOffset, yOffset int32) {
	DrawCircle(renderer, int32(s.center.X()) + xOffset, int32(s.center.Z()) + yOffset, int32(s.radius))
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

func (s *Sphere) solveQuadratic(r Ray.Ray) (bool, float64, float64) {
	loc := r.Origin().Minus(s.center)

	a := r.Direction().Dot(r.Direction())
	b := (r.Direction().Scale(2)).Dot(loc)
	c := loc.Dot(loc) - float64(s.radius * s.radius)

	discriminant := b*b - (4 * a * c)

	if discriminant < 0 {
		return false, 0, 0
	} else if discriminant == 0 {
		x0 := -0.5 * b / a
		return true, x0, x0
	} else {
		var q float64
		if b > 0 {
			q = -0.5 * (b - math.Sqrt(discriminant))
		} else {
			q = -0.5 * (b + math.Sqrt(discriminant))
		}
		x0 := q / a
		x1 := c / q

		if x0 > x1 {
			return true, x1, x0
		}
		return true, x0, x1
	}
}

func (s *Sphere) IntersectDistance(r Ray.Ray) (bool, float64) {
	bidirectionalIntersection, t0, t1 := s.solveQuadratic(r)

	if !bidirectionalIntersection {
		return false, 0
	}

	if t0 < 0 {
		t0 = t1
		if t0 < 0 {
			return false, 0
		}
	}

	return true, t0

}

func (s *Sphere) Radius() int {
	return s.radius
}


func DrawCircle(renderer *sdl.Renderer, centreX, centreY, radius int32) {
	x := radius - 1
	y := int32(0)
	tx := int32(1)
	ty := int32(1)
	error := tx - radius*2

	for x >= y {
		//  Each of the following renders an octant of the circle
		renderer.DrawPoint(centreX+x, centreY-y)
		renderer.DrawPoint(centreX+x, centreY+y)
		renderer.DrawPoint(centreX-x, centreY-y)
		renderer.DrawPoint(centreX-x, centreY+y)
		renderer.DrawPoint(centreX+y, centreY-x)
		renderer.DrawPoint(centreX+y, centreY+x)
		renderer.DrawPoint(centreX-y, centreY-x)
		renderer.DrawPoint(centreX-y, centreY+x)

		if error <= 0 {
			y++
			error += ty
			ty += 2
		}

		if error > 0 {
			x--
			tx += 2
			error += tx - radius*2
		}
	}
}
