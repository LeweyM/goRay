package Object

import (
	"goRay/Ray"
)

type Object interface {
	IntersectDistance(ray Ray.Ray) (bool, float64)
}
