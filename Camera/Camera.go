package Camera

import (
	"goRay/Object"
	"goRay/Ray"
	"goRay/Vector"
	"image/color"
	"math"
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

			c.primaryRays = append(c.primaryRays, primaryRay)
			c.pixelList = append(c.pixelList, c.getPixel(xIndex, yIndex, primaryRay))
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

				list[y * c.width + x] = c.getPixel(x, y, primaryRay)
			}
		}
	}

	var wg sync.WaitGroup
	list := make([]Pixel, c.width * c.height)

	cpuSplitFactor := 6

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

func (c *Camera) getPixel(x, y int, ray Ray.Ray) Pixel {
	for _, object := range c.ObjectList {
		intersects, t := object.IntersectDistance(ray)

		if intersects {
			hitNormal := object.GetHitNormal(ray, t)
			colorVector := object.GetSurfaceColor()
			facingRatio := hitNormal.Dot(ray.Direction().Reverse())
			facingRatio = math.Max(0, facingRatio)

			rr, gg, bb := scaleRGB(colorVector.X(), colorVector.Y(), colorVector.Z(), facingRatio)

			return Pixel{
				color: color.RGBA{R: rr, G: gg, B: bb, A: 255},
				x:     x,
				y:     y,
			}
		}
	}

	return getBackgroundPixel(ray, x, y)
}

func getBackgroundPixel(ray Ray.Ray, x, y int) Pixel {
	dir := ray.Direction()
	t := dir.Y()*0.5 + 1
	blue := Vector.New(0.5, 0.7, 1.0)
	white := Vector.New(1, 1, 1)
	lerp := white.Scale(1 - t).Translate(blue.Scale(t))

	rr, gg, bb := scaleRGB(lerp.X(), lerp.Y(), lerp.Z(), 1.0)

	return Pixel{
		color: color.RGBA{R: rr, G: gg, B: bb, A: 255},
		x:     x,
		y:     y,
	}
}

func scaleRGB(r, g, b, f float64) (uint8, uint8, uint8) {
	scale := f * 255.99
	rr := uint8(scale * r)
	gg := uint8(scale * g)
	bb := uint8(scale * b)
	return rr, gg, bb
}

func getScreenMatrix(scale float64) []float64 {
	var cells []float64
	for i := 0.0; i < scale; i++ {
		cells = append(cells, -(scale/2)+(i+0.5))
	}
	return cells
}
