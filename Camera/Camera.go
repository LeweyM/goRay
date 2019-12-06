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
	rotationVector := c.cameraRotationTransformer(*Vector.New(0.0,0.0,1.0))
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
			headingVector    := c.ScreenCellMatrix[yIndex][xIndex]
			rotatedVector    := c.cameraRotationTransformer(*headingVector)

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
	yRotateFn := Vector.RotateYBuilder(c.YRotation)

	rayWorker := func(wg *sync.WaitGroup, mu *sync.Mutex, x, y int) {
		defer wg.Done()

		headingVector := c.ScreenCellMatrix[y][x]
		rotatedHeadingVector := yRotateFn(*headingVector)

		translatedRotatedVector := rotatedHeadingVector.Translate(c.CameraPosition)

		ray := Ray.New(c.origin, translatedRotatedVector)

		mu.Lock()
		c.pixelList = append(c.pixelList, c.getPixel(x, y, ray))
		mu.Unlock()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for yIndex := 0; yIndex < c.height; yIndex++ {
		for xIndex := 0; xIndex < c.width; xIndex++ {
			wg.Add(1)
			go rayWorker(&wg, &mu, xIndex, yIndex)
		}
	}

	wg.Wait()
	return c.pixelList
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

			rr, gg, bb := brighten(colorVector.X(), colorVector.Y(), colorVector.Z(), facingRatio)

			return Pixel{
				color: color.RGBA{R: rr, G: gg, B: bb, A: 255},
				x:     x,
				y:     y,
			}
		}
	}

	return Pixel{
		color: color.RGBA{R: 0, G: 0, B: 160, A: 0},
		x:     x,
		y:     y,
	}
}

func brighten(r, g, b, f float64) (uint8, uint8, uint8) {
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
