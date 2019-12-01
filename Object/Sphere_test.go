package Object

import (
	"goRay/Ray"
	"goRay/Vector"
	"math"
	"testing"
)

func TestSphereIntersectionTests(t *testing.T) {
	headingForwardOnXAxis := *Vector.New(1, 0, 0)
	headingForwardOnXandYAxis := *Vector.New(1, 1, 0)
	headingBackwardsOnXAxis := *Vector.New(-1, 0, 0)
	tests := []struct{
		rayDirection Vector.Vector
		spherePosition Vector.Vector
		intersects bool
	}{
		{
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(10,0,0),
			intersects:     true,
		},{
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(0,10,0),
			intersects:     false,
		},{
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(-10,0,0),
			intersects:     false,
		},{
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(10,3,0),
			intersects:     true,
		},{
			rayDirection:   headingBackwardsOnXAxis,
			spherePosition: *Vector.New(-100,0,0),
			intersects:     true,
		},{
			rayDirection:   headingBackwardsOnXAxis.RotateY(math.Pi),
			spherePosition: *Vector.New(100,0,0),
			intersects:     true,
		},{
			rayDirection:   headingForwardOnXandYAxis,
			spherePosition: *Vector.New(100,100,0),
			intersects:     true,
		},
	}

	for i, tt := range tests {
		intersects, _ := testRaySphereIntersection(tt.rayDirection, tt.spherePosition)

		if intersects != tt.intersects {
			t.Errorf("Test %d: Expected interesection to be '%t', got '%t'",
				i + 1, tt.intersects, intersects, )
		}
	}
}

func testRaySphereIntersection(rayDirection Vector.Vector, spherePosition Vector.Vector) (bool, float64) {
	origin := Vector.New(0, 0, 0)
	sphere := NewSphere(spherePosition, 3)
	ray := Ray.New(*origin, rayDirection)
	return sphere.IntersectDistance(ray)
}

func TestSphereDistanceTests(t *testing.T) {

	rays := []Ray.Ray{
		Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0)),
		Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0)),
		Ray.New(*Vector.New(0, 0, 0), Vector.New(1, 1, 1).Normalize()),
	}

	spheres := []Sphere{
		*NewSphere(*Vector.New(10, 0, 0), 1),
		*NewSphere(*Vector.New(100, 0, 0), 10),
		*NewSphere(*Vector.New(10, 10, 10), 1),
	}

	results := []float64{
		9,
		90,
		Pythagorus3d(10, 10, 10) - 1,
	}

	for i, ray := range rays {
		ok, distance := spheres[i].IntersectDistance(ray)

		if !ok {
			t.Fatalf("Test %d: Doesn't intersect", i)
		}

		delta := math.Abs(distance - results[i])
		errorMargin := 0.0000000001
		if delta > errorMargin {
			t.Errorf("Test %d: Expected: %g got: %g",
				i, results[i], distance)
		}
	}
}

func Pythagorus3d(x, y, z float64) float64 {
	hyp2d := math.Sqrt(x*x + y*y)
	hyp3d := math.Sqrt(z*z + hyp2d*hyp2d)
	return hyp3d
}
