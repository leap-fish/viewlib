package viewlib

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Camera can look at positions, zoom and rotate.
// The Camera implementation is a modified https://github.com/MelonFunction/ebiten-camera, in order to support how
// viewlib performs draw layering.
type Camera struct {
	X, Y, Scale   float64
	Width, Height int
	Surface       *ebiten.Image

	Screen *ebiten.Image
}

// NewCamera returns a new Camera.
func NewCamera(width, height int, x, y, zoom float64) *Camera {
	return &Camera{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		Scale:   zoom,
		Surface: ebiten.NewImage(width, height),
	}
}

// SetPosition looks at a position
func (c *Camera) SetPosition(x, y float64) *Camera {
	c.X = x
	c.Y = y
	return c
}

// MovePosition moves the Camera by x and y.
// Use SetPosition if you want to set the position
func (c *Camera) MovePosition(x, y float64) *Camera {
	c.X += x
	c.Y += y
	return c
}

// Zoom *= the current zoom
func (c *Camera) Zoom(mul float64) *Camera {
	c.Scale *= mul
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(c.Width, c.Height)
	return c
}

// SetZoom sets the zoom
func (c *Camera) SetZoom(zoom float64) *Camera {
	c.Scale = zoom
	if c.Scale <= 0.01 {
		c.Scale = 0.01
	}
	c.Resize(c.Width, c.Height)
	return c
}

// Resize resizes the camera Surface
func (c *Camera) Resize(w, h int) *Camera {
	c.Width = w
	c.Height = h
	newW := int(float64(w) * 1.0 / c.Scale)
	newH := int(float64(h) * 1.0 / c.Scale)
	if newW <= 16384 && newH <= 16384 {
		c.Surface.Dispose()
		c.Surface = ebiten.NewImage(newW, newH)
	}
	return c
}

func (c *Camera) GetTranslation(ops *ebiten.DrawImageOptions, x, y float64) *ebiten.DrawImageOptions {
	w, h := c.Surface.Size()
	ops.GeoM.Translate(float64(w)/2, float64(h)/2)
	ops.GeoM.Translate(-c.X+x, -c.Y+y)
	return ops
}

func (c *Camera) GetRotation(ops *ebiten.DrawImageOptions, rot, originX, originY float64) *ebiten.DrawImageOptions {
	ops.GeoM.Translate(originX, originY)
	ops.GeoM.Rotate(rot)
	ops.GeoM.Translate(-originX, -originY)
	return ops
}

func (c *Camera) GetScale(ops *ebiten.DrawImageOptions, scaleX, scaleY float64) *ebiten.DrawImageOptions {
	ops.GeoM.Scale(scaleX, scaleY)
	return ops
}

func (c *Camera) GetSkew(ops *ebiten.DrawImageOptions, skewX, skewY float64) *ebiten.DrawImageOptions {
	ops.GeoM.Skew(skewX, skewY)
	return ops
}

// GetScreenCoords converts world coords into screen coords
func (c *Camera) GetScreenCoords(x, y float64) (float64, float64) {
	w, h := c.Width, c.Height

	x, y = x-c.X, y-c.Y
	x, y = x*c.Scale, y*c.Scale

	// Translate to screen center
	return x + float64(w)/2, y + float64(h)/2
}

// GetWorldCoords converts screen coords into world coords
func (c *Camera) GetWorldCoords(x, y float64) (float64, float64) {
	w, h := c.Width, c.Height

	x, y = x-float64(w)/2, y-float64(h)/2

	// Translate the coordinates
	x += float64(w) / 2
	y += float64(h) / 2

	return x, y
}

// Center returns the center point of the camera, based on its Width and Height.
func (c *Camera) Center() (float64, float64) {
	return float64(c.Width) * 0.5, float64(c.Height) * 0.5
}

// GetCursorCoords converts cursor/screen coords into world coords
func (c *Camera) GetCursorCoords() (float64, float64) {
	cx, cy := ebiten.CursorPosition()
	return c.GetWorldCoords(float64(cx), float64(cy))
}

// WorldMatrix modifies the `ops` parameter to be world relative.
func (c *Camera) WorldMatrix(ops *ebiten.DrawImageOptions) {
	centerX, centerY := c.Center()
	ops.GeoM.Scale(c.Scale, c.Scale)
	ops.GeoM.Translate(centerX, centerY)
	ops.GeoM.Translate(-c.X*c.Scale, -c.Y*c.Scale)
}
