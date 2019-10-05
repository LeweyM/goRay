package Object

import (
	"fmt"
	"goRay/Ray"
	"goRay/Vector"
)

type Sphere struct {
	position Vector.Vector
	radius   int
}

func (s *Sphere) String() string {
	return fmt.Sprintf("{postion: %s, radius: %b}", s.position, s.radius)
}

func New(position Vector.Vector, radius int) Sphere {
	return Sphere{position: position, radius: radius}
}

func (s *Sphere) intersects(r Ray.Ray) bool {
	// r.dir.x * ? == s.pos.x
	// ? == s.pos.x / r.dir.x
	// if ? > 0, intersects

	offsetDeltaX := float64(s.radius) - r.GetPosition().X()

	intersectsOnXAxis := s.position.X() / r.GetDirection().X() > 0 && offsetDeltaX >= 0
	intersectsOnYAxis := s.position.Y() / r.GetDirection().Y() > 0

	if r.GetDirection().X() == 0 {
		intersectsOnXAxis = false
	}
	if r.GetDirection().Y() == 0 {
		intersectsOnYAxis = false
	}

	return intersectsOnXAxis || intersectsOnYAxis
}
