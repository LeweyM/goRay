package Ray

import (
	"fmt"
	"goRay/Vector"
)

type Ray struct {
	position  Vector.Vector
	direction Vector.Vector
}

func (r Ray) String() string {
	return fmt.Sprintf("{postion: %s, direction: %s}", r.position, r.direction)
}

func (r Ray) Origin() Vector.Vector {
	return r.position
}

func (r Ray) Direction() Vector.Vector {
	return r.direction
}

func New(position, direction Vector.Vector) Ray {
	return Ray{position: position, direction: direction}
}


