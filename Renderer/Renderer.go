package Renderer

import (
	"github.com/veandco/go-sdl2/sdl"
	"goRay/Camera"
	"goRay/Object"
	"goRay/Vector"
	"time"
)

func Render(w, h int32, camera Camera.Camera) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("GoTracer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	var sphere = Object.NewSphere(*Vector.New(0, 0, 0), 30)
	var vector = Vector.New(0, 0, 2)

	running := true
	for running {

		time.Sleep(1000)
		fillSurface(camera.CastRays(), w, h, surface)
		camera.ClearObjects()
		vector = Vector.New(vector.X(), vector.Y(), vector.Z() + 1)
		sphere = Object.NewSphere(*vector, sphere.Radius())
		camera.SetObject(sphere)

		window.UpdateSurface()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func fillSurface(pixels []Camera.Pixel, w int32, h int32, surface *sdl.Surface) {
	for _, p := range pixels {
		wUnit := 800 / w
		hUnit := 600 / h

		rect := sdl.Rect{X: int32(p.X()) * wUnit, Y: int32(p.Y()) * hUnit, W: wUnit, H: hUnit}
		format := sdl.PixelFormat{
			Format: sdl.PIXELFORMAT_ABGR4444,
		}
		r, g, b, a := p.Color().RGBA()
		color := sdl.MapRGBA(&format, uint8(r), uint8(g), uint8(b), uint8(a))

		err := surface.FillRect(&rect, color)

		if err != nil {
			panic(err)
		}
	}
}