package main

import (
	"goRay/Camera"
	"goRay/Object"
	"goRay/Renderer"
	"goRay/Vector"
)

func main() {

	h := 100
	w := 100

	camOrigin := Vector.New(0, 0, 100)
	camera := Camera.New(w, h, *camOrigin)

	red := *Vector.New(1.0, 0.0, 0.0)
	green := *Vector.New(0.0, 1.0, 0.0)
	purple := *Vector.New(1.0, 0.0, 1.0)

	sphere := Object.NewSphere(*Vector.New(0, 0, 50), red, 10)
	sphere1 := Object.NewSphere(*Vector.New(20, 10, 50), green, 10)
	sphere2 := Object.NewSphere(*Vector.New(40, 5, 50), purple, 10)

	camera.SetObject(sphere)
	camera.SetObject(sphere1)
	camera.SetObject(sphere2)

	Renderer.Render(int32(w), int32(h), *camera)

}
