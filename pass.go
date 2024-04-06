package viewlib

import "github.com/hajimehoshi/ebiten/v2"

type Pass struct {
	Surface *ebiten.Image
	camera  *Camera
	canvas  bool
}

func (p *Pass) Draw(image *ebiten.Image, ops *ebiten.DrawImageOptions) {
	// If we're in canvas mode, we draw the image as it should be at the correct position already.
	if p.canvas {
		p.Surface.DrawImage(image, ops)
	} else {
		// If we are in non-canvas mode we have to modify the image with data from our camera.
		p.camera.WorldMatrix(ops)
		p.Surface.DrawImage(image, ops)
	}
}

// UpdateCamera can be used if the camera's width/height changes
func (p *Pass) UpdateCamera(camera *Camera) {
	p.camera = camera
	p.Surface = ebiten.NewImage(camera.Width, camera.Height)
}

func NewRenderPass(camera *Camera) *Pass {
	return &Pass{
		Surface: ebiten.NewImage(camera.Width, camera.Height),
		camera:  camera,
	}
}

func NewCanvasRenderPass(camera *Camera) *Pass {
	return &Pass{
		Surface: ebiten.NewImage(camera.Width, camera.Height),
		camera:  camera,
		canvas:  true,
	}
}
