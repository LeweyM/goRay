package Camera

import (
	"goRay/Object"
	"goRay/Ray"
	"goRay/Vector"
	"image/color"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

func TestCamera_GetPixelHeadingVector(t *testing.T) {
	type args struct {
		pHeight  float64
		pWidth   float64
		unitZoom float64
	}
	tests := []struct {
		name string
		args args
		want *Vector.Vector
	}{
		{
			name: "looking straight along z axis",
			args: args{
				pHeight:  0.0,
				pWidth:   0.0,
				unitZoom: 1.0,
			},
			want: Vector.New(0, 0, 1),
		},
		{
			name: "looking 90degrees right from z axis",
			args: args{
				pHeight:  0.0,
				pWidth:   1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees left from z axis",
			args: args{
				pHeight:  0.0,
				pWidth:   -1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(-0.7071067811, 0, 0.7071067811),
		},
		{
			name: "looking 90degrees up from z axis",
			args: args{
				pHeight:  1.0,
				pWidth:   0.0,
				unitZoom: 1.0,
			},
			want: Vector.New(0, 0.7071067811, 0.7071067811),
		},
		{
			name: "looking 90degrees up and left from z axis",
			args: args{
				pHeight:  1.0,
				pWidth:   -1.0,
				unitZoom: 1.0,
			},
			want: Vector.New(-0.5, 0.7071067811, 0.5), // BUT WHY THO???
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPixelHeadingVector(tt.args.pHeight, tt.args.pWidth, tt.args.unitZoom)

			if vectorIsEqual(got, tt.want) {
				t.Errorf("GetPixelHeadingVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getScreenMatrix(t *testing.T) {
	type args struct {
		height float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "scale 3",
			args: args{
				height: 3,
			},
			want: []float64{-1, 0, 1},
		},
		{
			name: "scale 5",
			args: args{
				height: 5,
			},
			want: []float64{-2, -1, 0, 1, 2},
		},
		{
			name: "scale 4",
			args: args{
				height: 4,
			},
			want: []float64{-1.5, -0.5, 0.5, 1.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getScreenMatrix(tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScreenMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCamera_GetPixel(t *testing.T) {
	zRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(0, 0, 1))
	type fields struct {
		height       int
		width        int
		origin       Vector.Vector
		objectList   []Object.Object
		pixelList    []Pixel
	}
	type args struct {
		x   int
		y   int
		ray Ray.Ray
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Pixel
	}{
		{
			name: "without objects",
			fields: fields{
				origin: Vector.Vector{},
			},
			args: args{
				x:   0,
				y:   0,
				ray: zRay,
			},
			want: Pixel{
				color: color.Black,
				x:     0,
				y:     0,
			},
		},
		{
			name: "with an object in front",
			fields: fields{
				origin:     Vector.Vector{},
				objectList: []Object.Object{Object.NewSphere(*Vector.New(0, 0, 5), 2, )},
			},
			args: args{
				x:   0,
				y:   0,
				ray: zRay,
			},
			want: Pixel{
				color: color.White,
				x:     0,
				y:     0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Camera{
				height:       tt.fields.height,
				width:        tt.fields.width,
				origin:       tt.fields.origin,
				ObjectList:   tt.fields.objectList,
				pixelList:    tt.fields.pixelList,
			}
			if got := c.GetPixel(tt.args.x, tt.args.y, tt.args.ray); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPixel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCamera_GetPixelHeadingVector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetPixelHeadingVector(10, 10, 1)
	}
}

func BenchmarkCamera_CastRays10(b *testing.B) {benchmarkCastRays(b, 10)}
func BenchmarkCamera_CastRays1000(b *testing.B) {benchmarkCastRays(b, 1000)}
func BenchmarkCamera_CastRays1000000(b *testing.B) {benchmarkCastRays(b, 1000000)}

func benchmarkCastRays(b *testing.B, n int) {
	c := New(3, 3, *Vector.New(0, 0, 0))
	setNSpheres(c, n)
	for i := 0; i < b.N; i++ {
		c.CastRays()
	}
}
//func benchmarkCastRaysConcurrent(b *testing.B, n int) {
//	camera := Camera{}
//	c := camera.New(3, 3, *Vector.New(0, 0, 0))
//	setNSpheres(c, n)
//	for i := 0; i < b.N; i++ {
//		c.CastRaysConcurrent()
//	}
//}

func TestCastRays(t *testing.T) {
	quarterTurnRight := math.Pi / 2
	tests := []struct {
		translation Vector.Vector
		rotationRads float64
		intersects bool
	}{
		{
			translation: *Vector.New(0,0,0),
			rotationRads:0,
			intersects: true,
		},
		{
			translation: *Vector.New(0,0,100),
			rotationRads:0,
			intersects: false,
		},
		{
			translation:  *Vector.New(0,0,100),
			rotationRads: quarterTurnRight * 2,
			intersects:   true,
		},
		{
			translation:  *Vector.New(50,0,50),
			rotationRads: quarterTurnRight * 3,
			intersects:   true,
		},
		{
			translation:  *Vector.New(50,0,0),
			rotationRads: quarterTurnRight * 3.5,
			intersects:   true,
		},
	}

	for i, tt := range tests {
		centerRayPixel := testCenterRayOfCamera(&tt.translation, tt.rotationRads)
		if tt.intersects != (centerRayPixel.Color() != color.Black) {
			t.Errorf("Test %d: Center of camera: Expected intersection to be '%t'", i + 1, tt.intersects)
		}
	}
}

func TestWalking(t *testing.T) {
	spherePosition := Vector.New(0,0, 50)
	origin := Vector.New(0,0, 0)

	camera := New(1, 1, *origin)
	camera.SetObject(Object.NewSphere(*spherePosition, 1))

	turnQuarterLeft := func() {turnLeft(camera, 16)}
	turnEighthLeft := func() {turnLeft(camera, 8)}
	objectVisible := func(i int) {checkIntersection(camera, t,true, i)}
	objectInvisible := func(i int) {checkIntersection(camera, t,false, i)}


	objectVisible(1)
	walkForward(camera, 50) //standing on sphere
	objectVisible(2)
	walkForward(camera, 2) //standing on edge of sphere
	objectInvisible(3)
	walkForward(camera, 8)
	turnQuarterLeft()
	turnQuarterLeft()
	objectVisible(4)
	walkForward(camera, 10)
	//back to origin
	turnQuarterLeft()
	walkForward(camera, 200)
	turnQuarterLeft()
	turnQuarterLeft()
	objectVisible(5)
	turnQuarterLeft()
	walkForward(camera, 200)
	turnQuarterLeft()
	turnQuarterLeft()
	turnEighthLeft()
	objectVisible(6)
	walkForward(camera, 400)
	turnQuarterLeft()
	turnQuarterLeft()
	objectVisible(7)
	walkForward(camera, 200)
	//back to origin
}

func walkForward(c *Camera, steps int) {
	for i := 0;  i<steps; i++ {
		c.IncrementForward()
	}
}

func turnLeft(c *Camera, steps int) {
	for i := 0;  i<steps; i++ {
		c.IncrementYRotation()
	}
}

func checkIntersection(c *Camera, t *testing.T, shouldIntersect bool, testNum int) {
	centerPixel := c.CastRays()[0]
	if shouldIntersect && centerPixel.color == color.Black {
		t.Errorf("Test %d: Should hit object but doesn't", testNum)
	}
	if !shouldIntersect && centerPixel.color != color.Black {
		t.Errorf("Test %d: Shouldn't hit object but does", testNum)
	}
}

func testCenterRayOfCamera(translation *Vector.Vector, rotationRads float64) Pixel {
	spherePosition := Vector.New(0,0, 50)
	origin := Vector.New(0,0, 0)

	camera := New(1, 1, *origin)
	camera.SetObject(Object.NewSphere(*spherePosition, 3))
	camera.RotateCamera(rotationRads)
	camera.TranslateCamera(*translation)

	rays := camera.CastRays()
	return rays[0]
}

func setNSpheres(camera *Camera, n int) {
	for i := 0; i < n; i++ {
		camera.SetObject(getRandomSphere())
	}
}

func getRandomSphere() *Object.Sphere {
	vector := Vector.New(rand.Float64(), rand.Float64(), rand.Float64())
	return Object.NewSphere(*vector, rand.Int())
}

func vectorIsEqual(v1 *Vector.Vector, v2 *Vector.Vector) bool {
	return withinErrorMargin(v1.X(), v2.X()) || withinErrorMargin(v1.Y(), v2.Y()) || withinErrorMargin(v1.Z(), v2.Z())
}

func withinErrorMargin(f1 float64, f2 float64) bool {
	return math.Abs(f1-f2) > 0.00000001
}
