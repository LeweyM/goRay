package Object

import "goRay/Ray"

type Object interface {
	intersects(ray Ray.Ray) bool
}
