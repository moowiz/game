package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	Update(w *glfw.Window, elapsed float64)
	Setup(program uint32)
}

type PositionSource func() mgl32.Vec3
