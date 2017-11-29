package player

import (
	"fmt"

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
		x, y := w.GetCursorPos()
		obj, vec, err := l.Raytrace(float32(x), float32(y))
		fmt.Printf("got %v %v %v FROM calling raytrace\n", obj, vec, err)
	}
	return false
}

func (i *Isometric) Update(elapsed float64) {
}

func (i *Isometric) UpdateWindow(w *glfw.Window, elapsed float64) {
}
