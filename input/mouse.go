package input

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type MouseInput struct {
	Button glfw.MouseButton
	Action glfw.Action
	Mods   glfw.ModifierKey
}
