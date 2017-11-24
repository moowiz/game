package player

import (
	"github.com/go-gl/glfw/v3.1/glfw"

	"github.com/moowiz/game/input"
)

type Player interface {
	HandleInput(ki *input.KeyInput) bool
	UpdateWindow(w *glfw.Window, elapsed float64)
	Update(elapsed float64)
}
