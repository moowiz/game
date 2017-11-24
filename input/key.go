package input

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type KeyInput struct {
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mods     glfw.ModifierKey
}
