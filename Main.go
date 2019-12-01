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

	sphere := Object.NewSphere(*Vector.New(0, 0, 50), 10)
	sphere1 := Object.NewSphere(*Vector.New(20, 0, 50), 10)
	sphere2 := Object.NewSphere(*Vector.New(40, 0, 50), 10)

	camera.SetObject(sphere)
	camera.SetObject(sphere1)
	camera.SetObject(sphere2)

	Renderer.Render(int32(w), int32(h), *camera)

}
