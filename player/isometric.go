package player

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/moowiz/game/camera"
	"github.com/moowiz/game/input"
)

type Isometric struct {
	Camera *camera.Isometric
}

func NewIsometric(source camera.PositionSource) *Isometric {
	return &Isometric{
		Camera: camera.NewIsometric(source),
	}
}

func (i *Isometric) HandleInput(ki *input.KeyInput) bool {
	return false
}

func (i *Isometric) Update(elapsed float64) {
}

func (i *Isometric) UpdateWindow(w *glfw.Window, elapsed float64) {
}
