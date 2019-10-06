package Camera

import (
	"goRay/Object"
	"goRay/Ray"
	"goRay/Vector"
	"image/color"
	"math"
)

type PixelGrabber interface {
	CastRays() []Pixel
}

type Camera struct {
	widthMatrix  []float64
	heightMatrix []float64
	height       int
	width        int
	origin       Vector.Vector
	objectList   []Object.Object
	pixelList    []Pixel
}

func (c *Camera) SetObject(object Object.Object) {
	c.objectList = append(c.objectList, object)
}

func (c *Camera) ClearObjects() {
	c.objectList = []Object.Object{}
}

func (c *Camera) New(width int, height int, origin Vector.Vector) *Camera {
	return &Camera{
		widthMatrix:  getScreenMatrix(float64(width)),
		heightMatrix: getScreenMatrix(float64(height)),
		height:       height,
		width:        width,
		origin:       Vector.Vector{},
		objectList:   []Object.Object{},
	}
}

func (c *Camera) CastRays() []Pixel {
	c.pixelList = []Pixel{}

	for yIndex, y := range c.heightMatrix {
		for xIndex, x := range c.widthMatrix {
			ray := Ray.New(c.origin, *c.GetPixelHeadingVector(y, x, float64(c.height))) //TODO
			c.pixelList = append(c.pixelList, c.GetPixel(xIndex, yIndex, ray))
		}
	}
	return c.pixelList
}

// get heading vector for pixel from
// worldOrigin by pixelHeight and pixelWidth
func (c *Camera) GetPixelHeadingVector(pYOffset, pXOffset, zoomUnit float64) *Vector.Vector {
	zAngle := math.Atan(pXOffset / zoomUnit)
	yAngle := math.Atan(pYOffset / zoomUnit)

	// pixel heading vector
	z := math.Cos(zAngle) * math.Cos(yAngle)
	x := math.Sin(zAngle) * math.Cos(yAngle)
	y := math.Sin(yAngle)

	return Vector.New(x, y, z)
}

type Pixel struct {
	color color.Color
	x     int
	y     int
}

func (p Pixel) Y() int {
	return p.y
}

func (p Pixel) X() int {
	return p.x
}

func (p Pixel) Color() color.Color {
	return p.color
}

func (c *Camera) GetPixel(x, y int, ray Ray.Ray) Pixel {
	for _, object := range c.objectList {
		intersects, _ := object.IntersectDistance(ray)
		if intersects {
			return Pixel{
				color: color.White,
				x:     x,
				y:     y,
			}
		}
	}
	return Pixel{
		color: color.Black,
		x:     x,
		y:     y,
	}
}

func getScreenMatrix(scale float64) []float64 {
	var cells []float64
	for i := 0.0; i < scale; i++ {
		cells = append(cells, -(scale/2)+(i+0.5))
	}
	return cells
}
