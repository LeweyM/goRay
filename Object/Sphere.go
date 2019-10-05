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

	// sphere: (x−x0)^2+(y−y0)^2+(z−z0)^2=r^2.  where x-x0 is pos-center
	//sphere: dot(pos^2) == r^2
	// ray: pos + t*dir

	//    vec3 oc = r.origin() - center;
	//    float a = dot(r.direction(), r.direction());
	//    float b = 2.0 * dot(oc, r.direction());
	//    float c = dot(oc,oc) - radius*radius;
	//    float discriminant = b*b - 4*a*c;
	//    return (discriminant>0);

	oc := r.GetPosition().Minus(s.position)
	a := r.GetDirection().dot(r.GetPosition())

	// intersection:

	return intersectsOnXAxis || intersectsOnYAxis
}
