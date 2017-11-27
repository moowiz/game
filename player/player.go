package player

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/moowiz/game/camera"
	"github.com/moowiz/game/input"
)

type Player interface {
	Camera() camera.Camera
	HandleInput(ki *input.KeyInput) bool
	UpdateWindow(w *glfw.Window, elapsed float64)
	Update(elapsed float64)
}
