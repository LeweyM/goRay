package Object

import (
	"goRay/Ray"
	"goRay/Vector"
)

type Object interface {
	intersectDistance(ray Ray.Ray) (bool, Vector.Vector)
}
