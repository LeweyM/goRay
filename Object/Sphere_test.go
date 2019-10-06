package Object

import (
	"goRay/Ray"
	"goRay/Vector"
	"math"
	"testing"
)

func TestSphereIntersectionTests(t *testing.T) {
	xRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0))
	yRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(0, 1, 0))
	xyzRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 1, 1))

	rays := []Ray.Ray{
		xRay,
		xRay,
		xRay,
		yRay,
		yRay,
		Ray.New(*Vector.New(0, 10, 0), *Vector.New(1, 0, 0)),
		xyzRay,
	}

	spheres := []Sphere{
		*NewSphere(*Vector.New(1, 0, 0), 1),
		*NewSphere(*Vector.New(10, 0, 0), 1),
		*NewSphere(*Vector.New(0, 10, 0), 1),
		*NewSphere(*Vector.New(0, 10, 0), 1),
		*NewSphere(*Vector.New(0, -10, 0), 1),
		*NewSphere(*Vector.New(100, 0, 0), 10),
		*NewSphere(*Vector.New(100, 100, 100), 1),
	}

	results := []bool{
		true,
		true,
		false,
		true,
		false,
		true,
		true,
	}

	for i, ray := range rays {
		if intersects, _ := spheres[i].IntersectDistance(ray); intersects != results[i] {
			t.Errorf("Test %d: Expected ray %b intersection with sphere %b to be: %t",
				i, rays[i], spheres[i], results[i])
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
