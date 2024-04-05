package viewlib

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawOp struct {
	pass *Pass

	scale    float64
	rotation float32

	posX, posY       float64
	originX, originY float64

	skewX, skewY float64

	image *ebiten.Image

	ops    *ebiten.DrawImageOptions
	filter ebiten.Filter
}

// CenterOrigin sets the origin of the sprite to its center.
func (d *DrawOp) CenterOrigin() *DrawOp {
	bounds := d.image.Bounds().Size()
	d.originX = float64(bounds.X / 2)
	d.originY = float64(bounds.Y / 2)
	return d
}

// Skew sets the skew of the sprite.
func (d *DrawOp) Skew(skewX, skewY float64) *DrawOp {
	d.skewX = skewX
	d.skewY = skewY
	return d
}

// Origin sets the origin of the sprite.
func (d *DrawOp) Origin(originX, originY float64) *DrawOp {
	d.originX = originX
	d.originY = originY
	return d
}

// Scale sets the scale of the sprite.
func (d *DrawOp) Scale(scale float64) *DrawOp {
	d.scale = scale
	return d
}

// Rotation sets the rotation of the sprite.
func (d *DrawOp) Rotation(rotation float32) *DrawOp {
	d.rotation = rotation
	return d
}

// Position sets the position of the sprite.
func (d *DrawOp) Position(posX, posY float64) *DrawOp {
	d.posX = posX
	d.posY = posY
	return d
}

// Filter sets the draw filter mode.
func (d *DrawOp) Filter(filter ebiten.Filter) *DrawOp {
	d.filter = filter
	return d
}

func (d *DrawOp) Render() {
	d.ops.GeoM.Translate(-d.originX, -d.originY)
	d.ops.GeoM.Rotate(float64(d.rotation))
	d.ops.GeoM.Translate(d.originX, d.originY)
	spritePosX, spritePosY := d.posX-d.originX, d.posY-d.originY
	d.ops.GeoM.Translate(spritePosX, spritePosY)

	// Non essential operations are checked first
	if d.scale != 1 {
		d.ops.GeoM.Scale(d.scale, d.scale)
	}

	if d.skewX != 0 && d.skewY != 0 {
		d.ops.GeoM.Skew(d.skewX, d.skewY)
	}

	if d.filter != ebiten.FilterNearest {
		d.ops.Filter = d.filter
	}

	d.pass.Draw(d.image, d.ops)
}

// Draw returns a new DrawOp which can be used to customize how the image is rendered.
// By using Draw instead of manual ebiten drawing, you get automatic handling of rotations and sprite origins.
// DrawOp makes sure the draw operations are performed in the correct order.
// Call Render() to draw onto the render pass.
func Draw(pass *Pass, image *ebiten.Image) *DrawOp {
	return &DrawOp{
		pass:   pass,
		scale:  1.0,
		image:  image,
		filter: ebiten.FilterNearest,
		ops:    &ebiten.DrawImageOptions{},
	}
}

// LayerRenderable must be implemented somewhere and passed in when calling DrawPasses.
type LayerRenderable interface {
	LayerRender(image *ebiten.Image)
}

// DrawPasses clears and draws the screen based on its input arguments.
// The `screen` parameter is the image to draw the passes on.
// The `passes` parameter should be a depth sorted slice of render passes.
// When drawing, draw to the pass's Surface image.
// The `renderable` parameter is an interface which should be implemented on your Game or scene system, as a replacement
// since DrawPasses should be called from the original Draw(*ebiten.Image) function.
func DrawPasses(screen *ebiten.Image, passes []*Pass, renderable LayerRenderable) {
	for _, pass := range passes {
		if pass.Surface == nil {
			return
		}
		pass.Surface.Clear()
	}

	renderable.LayerRender(screen)

	for _, pass := range passes {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(pass.Surface, op)
	}
}
