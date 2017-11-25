package camera

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Isometric struct {
	Source   PositionSource
	angle    float64
	distance float32
	mat      mgl32.Mat4
}

func NewIsometric(source PositionSource) *Isometric {
	return &Isometric{
		Source:   source,
		distance: 10,
	}

}

func (i *Isometric) getMat4() mgl32.Mat4 {
	pos := i.Source()
	return mgl32.LookAt(
		i.distance, i.distance, i.distance,
		pos[0], pos[1], pos[2],
		0, 1, 0,
	)
}

func (i *Isometric) Setup(program uint32) {
	i.mat = i.getMat4()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &i.mat[0])
}

func (i *Isometric) Update(w *glfw.Window, elapsed float64) {
}
