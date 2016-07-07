package physics

import (
	"github.com/go-gl/mathgl/mgl32"
)

// Either a special source ("gravity"), or another object
type ForceSource struct {
	ID   string
	Body Body
}

type Force struct {
	Source    ForceSource
	Magnitude mgl32.Vec3
}
