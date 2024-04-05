package viewlib

import "github.com/hajimehoshi/ebiten/v2"

type Pass struct {
	Surface *ebiten.Image
	camera  *Camera
	canvas  bool
}

func (p *Pass) PreparePass(screen *ebiten.Image) {
	p.Surface = screen
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

func NewRenderPass(camera *Camera) *Pass {
	return &Pass{
		camera: camera,
	}
}

func NewCanvasRenderPass(camera *Camera) *Pass {
	return &Pass{
		camera: camera,
		canvas: true,
	}
}
