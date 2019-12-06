package Object

import (
	"github.com/veandco/go-sdl2/sdl"
	"goRay/Ray"
	"goRay/Vector"
)

type Object interface {
	IntersectDistance(ray Ray.Ray) (bool, float64)
	Draw(renderer *sdl.Renderer, xOffset, yOffset int32)
	GetHitNormal(ray Ray.Ray, t float64) Vector.Vector
}
