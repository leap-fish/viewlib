# viewlib
A small library for [Ebitengine](https://ebitengine.org/) for layered rendering and includes a basic camera implementation
modified from [ebiten-camera](https://github.com/MelonFunction/ebiten-camera).

# Layered rendering

```go
var (
	cam *viewlib.Camera

	renderPassContent = viewlib.NewRenderPass(cam)
	renderPassUI      = viewlib.NewCanvasRenderPass(cam)

	passes = []*viewlib.Pass{
		renderPassContent,
		renderPassUI,
	}
)

func main() {
	ebiten.SetWindowSize(1366, 768)
	ebiten.SetWindowTitle("viewlib")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) LayerRender(screen *ebiten.Image) {
	// Draw your game!
	viewlib.Draw(renderPassContent /* image */).Render()
}

func (g *Game) Draw(screen *ebiten.Image) {
	viewlib.DrawPasses(screen, passes, g)
}

```