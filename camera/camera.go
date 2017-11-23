package camera

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Camera interface {
	Update(w *glfw.Window, elapsed float64)
}
