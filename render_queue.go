package viewlib

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

// queuedRender represents a new render operation to be done in the future.
type queuedRender struct {
	op         *DrawOp
	renderFunc func(screen *ebiten.Image, camera *Camera)
	mode       RenderMode
	ordering   int
}

// RenderMode represents the type of rendering that will be done.
// RenderModeCanvas is used for screen space rendering.
// RenderModeWorld is used for world space rendering.
type RenderMode uint

const (
	RenderModeCanvas RenderMode = iota + 1
	RenderModeWorld
)

var renderQueue []*queuedRender

// QueueRender adds a draw operation to the render queue.
// The `mode` argument specifies whether the draw operation is executed on the world space layer or the canvas layer.
func QueueRender(op *DrawOp, mode RenderMode, layer int) {
	renderQueue = append(renderQueue, &queuedRender{
		op:         op,
		renderFunc: nil,
		mode:       mode,
		ordering:   layer,
	})

	sort.Slice(renderQueue, func(i, j int) bool {
		return renderQueue[i].ordering < renderQueue[j].ordering
	})
}

// QueueFunc is used to queue an arbitrary screen aware function to run.
func QueueFunc(renderFunc func(screen *ebiten.Image, camera *Camera), layer int) {
	renderQueue = append(renderQueue, &queuedRender{
		op:         nil,
		renderFunc: renderFunc,
		ordering:   layer,
	})

	sort.Slice(renderQueue, func(a, b int) bool {
		return renderQueue[a].ordering < renderQueue[b].ordering
	})
}

// RenderTo should be called by your render loop, and will draw all queued operations to the screen.
// It will clear the screen, and clear the queue.
func RenderTo(screen *ebiten.Image, camera *Camera) {
	screen.Clear()

	// Commit each sorted entry to the screen
	for _, entry := range renderQueue {
		if entry.op != nil {
			entry.op.commit(screen, camera)
		}

		if entry.renderFunc != nil {
			entry.renderFunc(screen, camera)
		}
	}

	// Clear the queue
	renderQueue = nil
}
