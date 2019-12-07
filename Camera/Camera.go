package Camera

import (
	"goRay/Object"
	"goRay/Ray"
	"goRay/Vector"
	"image/color"
	"math"
	"math/rand"
	"sync"
)

type PixelGrabber interface {
	CastRays() []Pixel
}

type Camera struct {
	height                    int
	width                     int
	origin                    Vector.Vector
	ObjectList                []Object.Object
	pixelList                 []Pixel
	ScreenCellMatrix          [][]*Vector.Vector
	YRotation                 float64
	cameraRotationTransformer func(Vector.Vector) Vector.Vector
	CameraPosition            Vector.Vector
	primaryRays               []Ray.Ray
	antiAliasingFactor        int
}

func New(width int, height int, origin Vector.Vector) *Camera {

	heightMatrix := getScreenMatrix(float64(height))
	widthMatrix := getScreenMatrix(float64(width))

	screenCellMatrix := make([][]*Vector.Vector, height)
	for row, y := range heightMatrix {
		screenCellMatrix[row] = make([]*Vector.Vector, width)
		for col, x := range widthMatrix {
			screenCellMatrix[row][col] = GetPixelHeadingVector(y, x, float64(height))
		}
	}

	return &Camera{
		ScreenCellMatrix:          screenCellMatrix,
		cameraRotationTransformer: Vector.RotateYBuilder(0),
		height:                    height,
		width:                     width,
		origin:                    origin,
		ObjectList:                []Object.Object{},
		YRotation:                 0,
		CameraPosition:            Vector.Vector{},
		antiAliasingFactor:        0,
	}
}

func (c *Camera) TranslateCamera(vector Vector.Vector) {
	c.CameraPosition = c.CameraPosition.Translate(vector)
}

func (c *Camera) RotateCamera(rads float64) {
	c.YRotation = rads
	c.UpdateCamRotationTransformer()
}

func (c *Camera) IncrementForward() {
	camDirectionVector := Vector.New(0, 0, 1).RotateY(c.YRotation)
	c.CameraPosition = c.CameraPosition.Translate(camDirectionVector)
}

func (c *Camera) DecrementForward() {
	camDirectionVector := Vector.New(0, 0, -1).RotateY(c.YRotation)
	c.CameraPosition = c.CameraPosition.Translate(camDirectionVector)
}

func (c *Camera) IncrementYRotation() {
	c.YRotation = c.YRotation + math.Pi/32
	c.UpdateCamRotationTransformer()
}

func (c *Camera) DecrementYRotation() {
	c.YRotation = c.YRotation - math.Pi/32
	c.UpdateCamRotationTransformer()
}

func (c *Camera) UpdateCamRotationTransformer() {
	c.cameraRotationTransformer = Vector.RotateYBuilder(c.YRotation)
}

func (c *Camera) GetRotationLine() (x1 float32, y1 float32, x2 float32, y2 float32) {
	rotationVector := c.cameraRotationTransformer(*Vector.New(0.0, 0.0, 1.0))
	return 0, 0, float32(rotationVector.X()), float32(rotationVector.Z())
}

func (c *Camera) SetObject(object Object.Object) {
	c.ObjectList = append(c.ObjectList, object)
}

func (c *Camera) ClearObjects() {
	c.ObjectList = []Object.Object{}
}

func (c *Camera) CastRays() []Pixel {
	c.pixelList = []Pixel{}
	c.primaryRays = []Ray.Ray{}

	for yIndex := 0; yIndex < c.height; yIndex++ {
		for xIndex := 0; xIndex < c.width; xIndex++ {
			headingVector := c.ScreenCellMatrix[yIndex][xIndex]
			rotatedVector := c.cameraRotationTransformer(*headingVector)

			primaryRay := Ray.New(c.CameraPosition, rotatedVector)

			var pixel Pixel
			if c.antiAliasingFactor > 0 {
				pixel = c.processAntiAliasing(headingVector, xIndex, yIndex, c.antiAliasingFactor)
			} else {
				r, g, b, a := colorVectorToRGB(c.getColor(primaryRay))
				pixel = Pixel{
					color: color.RGBA{R: r, G: g, B: b, A: a},
					x:     xIndex,
					y:     yIndex,
				}
			}
			c.pixelList = append(c.pixelList, pixel)
		}
	}

	return c.pixelList
}

func (c *Camera) GetPrimaryRays() []Ray.Ray {
	return c.primaryRays
}

func (c *Camera) CastRaysConcurrent() []Pixel {
	c.pixelList = []Pixel{}

	rayWorker := func(wg *sync.WaitGroup, list []Pixel, xStart, xEnd, yStart, yEnd int) {
		defer wg.Done()

		for y := yStart; y < yEnd; y++ {
			for x := xStart; x < xEnd; x++ {
				headingVector := c.ScreenCellMatrix[y][x]
				rotatedVector := c.cameraRotationTransformer(*headingVector)

				primaryRay := Ray.New(c.CameraPosition, rotatedVector)

				var pixel Pixel
				if c.antiAliasingFactor > 0 {
					pixel = c.processAntiAliasing(headingVector, x, y, c.antiAliasingFactor)
				} else {
					r, g, b, a := colorVectorToRGB(c.getColor(primaryRay))
					pixel = Pixel{
						color: color.RGBA{R: r, G: g, B: b, A: a},
						x:     x,
						y:     y,
					}
				}
				list[y*c.width+x] = pixel
			}
		}
	}

	var wg sync.WaitGroup
	list := make([]Pixel, c.width*c.height)

	cpuSplitFactor := 4

	for c.height%cpuSplitFactor != 0 && cpuSplitFactor != 1 {
		cpuSplitFactor--
	}

	yPixelGroupSize := c.height / cpuSplitFactor
	xPixelGroupSize := c.width / cpuSplitFactor
	for yStart := 0; yStart < c.height; yStart += yPixelGroupSize {
		for xStart := 0; xStart < c.width; xStart += xPixelGroupSize {
			wg.Add(1)
			go rayWorker(&wg, list, xStart, xStart+xPixelGroupSize, yStart, yStart+yPixelGroupSize)
		}
	}

	wg.Wait()
	return list
}

// get heading vector for pixel from
// worldOrigin by pixelHeight and pixelWidth
func GetPixelHeadingVector(pYOffset, pXOffset, zoomUnit float64) *Vector.Vector {
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

func (c *Camera) getColor(ray Ray.Ray) Vector.Vector {
	type objectWithDistance struct {
		object Object.Object
		t      float64
	}
	intersectionObject := objectWithDistance{
		object: nil,
		t:      math.MaxFloat64,
	}

	for _, object := range c.ObjectList {
		intersects, t := object.IntersectDistance(ray)

		if intersects && t < intersectionObject.t {
			intersectionObject = objectWithDistance{
				object: object,
				t:      t,
			}
		}
	}

	if intersectionObject.object != nil {
		return getColorFromObject(ray, intersectionObject.t, intersectionObject.object)
	} else {
		vector := getBackgroundColor(ray)
		return vector
	}
}

func getColorFromObject(ray Ray.Ray, t float64, object Object.Object) Vector.Vector {
	hitNormal := object.GetHitNormal(ray, t)
	colorVector := object.GetSurfaceColor()
	facingRatio := hitNormal.Dot(ray.Direction().Reverse())
	facingRatio = math.Max(0, facingRatio)

	return colorVector.Scale(facingRatio * 255.99)
}

func getBackgroundColor(ray Ray.Ray) Vector.Vector {
	dir := ray.Direction()
	t := dir.Y()*0.5 + 1
	blue := Vector.New(0.5, 0.7, 1.0)
	white := Vector.New(1, 1, 1)
	lerp := white.Scale(1 - t).Translate(blue.Scale(t))

	return lerp.Scale(255.99)
}

func colorVectorToRGB(colorVector Vector.Vector) (uint8, uint8, uint8, uint8) {
	r := uint8(colorVector.X())
	g := uint8(colorVector.Y())
	b := uint8(colorVector.Z())
	return r, g, b, 255
}

func getScreenMatrix(scale float64) []float64 {
	var cells []float64
	for i := 0.0; i < scale; i++ {
		cells = append(cells, -(scale/2)+(i+0.5))
	}
	return cells
}

func (c *Camera) processAntiAliasing(headingVector *Vector.Vector, xIndex, yIndex, aaFactor int) Pixel {
	var pixel Pixel
	var colorVector Vector.Vector
	for aa := 0; aa < aaFactor; aa++ {
		randX := rand.Float64() / float64(aaFactor)
		randY := rand.Float64() / float64(aaFactor)
		randomOffset := Vector.New(randX, randY, 0)
		aaHeadingVector := headingVector.Translate(*randomOffset).Normalize()
		aaRotatedVector := c.cameraRotationTransformer(aaHeadingVector)

		aaRay := Ray.New(c.CameraPosition, aaRotatedVector)

		colorVector = colorVector.Translate(c.getColor(aaRay))
	}
	r, g, b, a := colorVectorToRGB(colorVector.Scale(1 / float64(aaFactor)))

	pixel = Pixel{
		color: color.RGBA{R: r, G: g, B: b, A: a},
		x:     xIndex,
		y:     yIndex,
	}

	return pixel
}

func (c *Camera) SetAntiAliasing(aaFactor int) {
	c.antiAliasingFactor = aaFactor
}
