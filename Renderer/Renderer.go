package Renderer

import (
	"github.com/veandco/go-sdl2/sdl"
	"goRay/Camera"
	"goRay/Object"
	"goRay/Vector"
)

func Render(w, h int32, camera Camera.Camera) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("GoTracer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		600, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateSoftwareRenderer(surface)
	if err != nil {
		panic(err)
	}

	running := true
	for running {

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		drawPixels(camera.CastRays(), w, h, renderer)

		drawPrimaryRays(camera, renderer)
		drawVerticalPrimaryRays(camera, renderer)

		//drawRays(camera, renderer)

		//drawRotationLine(camera, renderer)

		//drawMiniMap(camera.CameraPosition, renderer)

		drawObjects(camera.ObjectList, renderer)

		renderer.Present()

		window.UpdateSurface()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_LEFT {
					camera.DecrementYRotation()
				}
				if e.Keysym.Sym == sdl.K_RIGHT {
					camera.IncrementYRotation()
				}
				if e.Keysym.Sym == sdl.K_UP {
					camera.IncrementForward()
				}
				if e.Keysym.Sym == sdl.K_DOWN {
					camera.DecrementForward()
				}
				break
			}
		}
	}
}

func drawObjects(objectList []Object.Object, renderer *sdl.Renderer) {
	for _, object := range objectList {
		object.Draw(renderer, 100, 100)
	}
}

func drawRays(camera Camera.Camera, renderer *sdl.Renderer) {
	for _, vectorList := range camera.ScreenCellMatrix {
		for _, vector := range vectorList {
			vector1 := vector.RotateY(camera.YRotation)
			vector2 := vector1.Translate(camera.CameraPosition)
			camx, camz := camera.CameraPosition.X(), camera.CameraPosition.Z()
			drawLine(float32(camx), float32(camz), float32(vector2.X()), float32(vector2.Z()), 100, renderer)
		}
	}
}

func drawPrimaryRays(camera Camera.Camera, renderer *sdl.Renderer) {
	camera.CastRays()
	for _, ray := range camera.GetPrimaryRays() {
		//hits := false
		for _, obj := range camera.ObjectList {
			intersects, _ := obj.IntersectDistance(ray)
			if intersects {
				//hits = true
			}
		}

		startPoint := ray.Origin()
		direction := ray.Direction()
		endPoint := startPoint.Translate(direction.Scale(30))

		x1 := startPoint.X()
		y1 := startPoint.Z()
		x2 := endPoint.X()
		y2 := endPoint.Z()

		//if hits {
		//	drawHitLine(float32(x1), float32(y1), float32(x2), float32(y2), 300, renderer)
		//} else {
			drawLine(float32(x1), float32(y1), float32(x2), float32(y2), 100, renderer)
		//}
	}
}

func drawVerticalPrimaryRays(camera Camera.Camera, renderer *sdl.Renderer) {
	camera.CastRays()
	for _, ray := range camera.GetPrimaryRays() {
		hits := false
		for _, obj := range camera.ObjectList {
			intersects, _ := obj.IntersectDistance(ray)
			if intersects {
				hits = true
			}
		}

		startPoint := ray.Origin()
		direction := ray.Direction()
		endPoint := startPoint.Translate(direction.Scale(40))


		x1 := startPoint.X()
		y1 := startPoint.Y()
		x2 := endPoint.X()
		y2 := endPoint.Y()

		if hits {
			drawHitLine(float32(x1), float32(y1), float32(x2), float32(y2), 300, renderer)
		} else {
			drawLine(float32(x1), float32(y1), float32(x2), float32(y2), 300, renderer)
		}
	}
}

func drawHitLine(x1, y1, x2, y2, offset float32, renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(0, 255, 0, 50)
	if err != nil {
		panic(err)
	}
	err = renderer.DrawLineF(x1+offset, y1+offset, x2+offset, y2+offset)
	if err != nil {
		panic(err)
	}
}

func drawLine(x1, y1, x2, y2, offset float32, renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(200, 100, 200, 50)
	if err != nil {
		panic(err)
	}
	err = renderer.DrawLineF(x1+offset, y1+offset, x2+offset, y2+offset)
	if err != nil {
		panic(err)
	}
}

func drawRotationLine(camera Camera.Camera, renderer *sdl.Renderer) {
	x1, y1, x2, y2 := camera.GetRotationLine()
	drawLine(x1, y1, x2, y2, 50, renderer)
}

func drawMiniMap(cameraPosition Vector.Vector, renderer *sdl.Renderer) {
	err := renderer.SetDrawColor(200, 255, 100, 50)
	if err != nil {
		panic(err)
	}
	xOffset := 100.0
	yOffset := 100.0

	x1, y1 := xOffset, yOffset
	x2, y2 := cameraPosition.X()*-30+xOffset, cameraPosition.Z()*-30+yOffset

	err = renderer.DrawLineF(float32(x1), float32(y1), float32(x2), float32(y2))
	if err != nil {
		panic(err)
	}
}

func drawPixels(pixels []Camera.Pixel, w int32, h int32, renderer *sdl.Renderer) {
	for _, p := range pixels {
		wUnit := 600 / w
		hUnit := 600 / h

		r, g, b, a := p.Color().RGBA()
		err := renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), uint8(a))
		if err != nil {
			panic(err)
		}

		rect := sdl.Rect{X: int32(p.X()) * wUnit, Y: int32(p.Y()) * hUnit, W: wUnit, H: hUnit}
		err = renderer.FillRect(&rect)
		if err != nil {
			panic(err)
		}
	}
}

