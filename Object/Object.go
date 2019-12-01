package Object

import (
	"github.com/veandco/go-sdl2/sdl"
	"goRay/Ray"
)

type Object interface {
	IntersectDistance(ray Ray.Ray) (bool, float64)
	Draw(renderer *sdl.Renderer, xOffset, yOffset int32)
}
