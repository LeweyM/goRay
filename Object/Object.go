package Object

import "goRay/Ray"

type Object interface {
	intersectDistance(ray Ray.Ray) bool
}
