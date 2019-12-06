package Object

import (
	"goRay/Ray"
	"goRay/Vector"
	"math"
	"testing"
)

var white = *Vector.New(1,1,1)

func TestGetIntersectionPoint(t *testing.T) {
	origin := Vector.Vector{}
	tests := []struct {
		spherePos          Vector.Vector
		intersectionVector Vector.Vector
		ray                Ray.Ray
		t                  float64
	}{
		{
			spherePos:          *Vector.New(0, 0, 50),
			ray:                Ray.New(origin, *Vector.New(0, 0, 1)),
			intersectionVector: *Vector.New(0, 0, -1),
			t:                  45,
		},
		{
			spherePos:          *Vector.New(5, 0, 50),
			ray:                Ray.New(origin, *Vector.New(0, 0, 1)),
			intersectionVector: *Vector.New(-1, 0, 0),
			t:                  50,
		},
		{
			spherePos:          *Vector.New(-5, 0, 50),
			ray:                Ray.New(origin, *Vector.New(0, 0, 1)),
			intersectionVector: *Vector.New(1, 0, 0),
			t:                  50,
		},
	}

	for i, tt := range tests {
		sphere := NewSphere(tt.spherePos, white, 5)

		res := sphere.GetHitNormal(tt.ray, tt.t)

		if tt.intersectionVector != res {
			t.Errorf("Test %d: Incorrect interception point, expected %v, got %v", i, tt.intersectionVector, res)
		}
	}
}

func TestSphereIntersectionTests(t *testing.T) {
	headingForwardOnXAxis := *Vector.New(1, 0, 0)
	headingForwardOnXandYAxis := *Vector.New(1, 1, 0)
	headingBackwardsOnXAxis := *Vector.New(-1, 0, 0)
	tests := []struct {
		rayDirection   Vector.Vector
		spherePosition Vector.Vector
		intersects     bool
	}{
		{
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(10, 0, 0),
			intersects:     true,
		}, {
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(0, 10, 0),
			intersects:     false,
		}, {
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(-10, 0, 0),
			intersects:     false,
		}, {
			rayDirection:   headingForwardOnXAxis,
			spherePosition: *Vector.New(10, 3, 0),
			intersects:     true,
		}, {
			rayDirection:   headingBackwardsOnXAxis,
			spherePosition: *Vector.New(-100, 0, 0),
			intersects:     true,
		}, {
			rayDirection:   headingBackwardsOnXAxis.RotateY(math.Pi),
			spherePosition: *Vector.New(100, 0, 0),
			intersects:     true,
		}, {
			rayDirection:   headingForwardOnXandYAxis,
			spherePosition: *Vector.New(100, 100, 0),
			intersects:     true,
		},
	}

	for i, tt := range tests {
		intersects, _ := testRaySphereIntersection(tt.rayDirection, tt.spherePosition)

		if intersects != tt.intersects {
			t.Errorf("Test %d: Expected interesection to be '%t', got '%t'",
				i+1, tt.intersects, intersects, )
		}
	}
}

func TestSphereDistanceTests(t *testing.T) {

	rays := []Ray.Ray{
		Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0)),
		Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0)),
		Ray.New(*Vector.New(0, 0, 0), Vector.New(1, 1, 1).Normalize()),
	}

	spheres := []Sphere{
		*NewSphere(*Vector.New(10, 0, 0), white, 1),
		*NewSphere(*Vector.New(100, 0, 0), white, 10),
		*NewSphere(*Vector.New(10, 10, 10), white, 1),
	}

	results := []float64{
		9,
		90,
		Pythagoras3d(10, 10, 10) - 1,
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

func Pythagoras3d(x, y, z float64) float64 {
	hyp2d := math.Sqrt(x*x + y*y)
	hyp3d := math.Sqrt(z*z + hyp2d*hyp2d)
	return hyp3d
}

func testRaySphereIntersection(rayDirection Vector.Vector, spherePosition Vector.Vector) (bool, float64) {
	origin := Vector.New(0, 0, 0)
	sphere := NewSphere(spherePosition, white, 3)
	ray := Ray.New(*origin, rayDirection)
	return sphere.IntersectDistance(ray)
}
