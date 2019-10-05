package Camera

import (
	"goRay/Vector"
	"math"
)

type Camera struct {
}

// get heading vector for pixel from
// worldOrigin by pixelHeight and pixelWidth
func (c *Camera) GetPixelHeadingVector(pHeight, pWidth float64) *Vector.Vector {
	zoomUnit := 1.0

	xAngle := math.Atan(pWidth / zoomUnit)
	yAngle := math.Atan(pHeight / zoomUnit)

	// pixel heading vector
	z := math.Cos(xAngle) * math.Cos(yAngle)
	x := math.Sin(xAngle) * math.Cos(yAngle)
	y := math.Sin(yAngle)

	return Vector.New(x, y, z)
}
