package Object

import (
	"goRay/Ray"
	"goRay/Vector"
	"testing"
)

func TestSphereIntersectionTests(t *testing.T) {
	offsetXRay := Ray.New(*Vector.New(1, 0, 0), *Vector.New(1, 0, 0))
	xRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(1, 0, 0))
	yRay := Ray.New(*Vector.New(0, 0, 0), *Vector.New(0, 1, 0))

	rays := []Ray.Ray{
		xRay,
		xRay,
		xRay,
		yRay,
		yRay,
		offsetXRay,
	}

	spheres := []Sphere{
		New(*Vector.New(1,0,0),1),
		New(*Vector.New(10,0,0),1),
		New(*Vector.New(0,10,0),1),
		New(*Vector.New(0,10,0),1),
		New(*Vector.New(0,-10,0),1),
		New(*Vector.New(10,0,0),1),
	}

	results := []bool {
		true,
		true,
		false,
		true,
		false,
		true,
	}
	
	for i, ray := range rays {
		if spheres[i].intersects(ray) != results[i] {
			t.Errorf("Test %d: Expected ray %b intersection with sphere %b to be: %t",
				i, rays[i], spheres[i], results[i])
		}
	}
}
