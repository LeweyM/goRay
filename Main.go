package main

import (
	"goRay/Camera"
	"goRay/Object"
	"goRay/Renderer"
	"goRay/Vector"
)

func main() {

	cam := Camera.Camera{}

	camOrigin := Vector.New(0, 0, 0)
	camera := cam.New(800, 600, *camOrigin)

	vector := Vector.New(0, 0, 10)
	sphere := Object.NewSphere(*vector, 2)

	camera.SetObject(sphere)

	Renderer.Render(800, 600, *camera)

}
