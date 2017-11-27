package player

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/moowiz/game/camera"
	"github.com/moowiz/game/input"
	"github.com/moowiz/game/level"
)

type Isometric struct {
	Camera *camera.Isometric
}

func NewIsometric(source camera.PositionSource) *Isometric {
	return &Isometric{
		Camera: camera.NewIsometric(source),
	}
}

func (i *Isometric) HandleInput(w *glfw.Window, l *level.Level, ki *input.KeyInput, mi *input.MouseInput) bool {
	if mi != nil {

	}
	return false
}

func (i *Isometric) Update(elapsed float64) {
}

func (i *Isometric) UpdateWindow(w *glfw.Window, elapsed float64) {
}
