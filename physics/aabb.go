package physics

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type AABB struct {
	Center   mgl32.Vec3 `json:"center"`
	HalfSize mgl32.Vec3 `json:"half_size"`
	vao      uint32
	program  uint32
}

func (a *AABB) String() string {
	return fmt.Sprintf("AABB center %v half size %v", a.Center, a.HalfSize)
}

func (a *AABB) Position() mgl32.Vec3 {
	return a.Center
}

func (a *AABB) SetPosition(new mgl32.Vec3) {
	a.Center = new
}

func (a *AABB) Collides(b body) bool {
	other, ok := b.(*AABB)
	if !ok {
		panic("Only AABB for now")
	}

	for ind := 0; ind < 3; ind += 1 {
		if mgl32.Abs(a.Center[ind]-other.Center[ind]) > (a.HalfSize[ind] + other.HalfSize[ind]) {
			return false
		}
	}

	return true
}

func (a *AABB) debugDraw(program uint32) {
	if a.vao == 0 {
		panic("debug draw with no vao")
	}

	model := mgl32.Translate3D(a.Center[0], a.Center[1], a.Center[2])
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(a.vao)
	gl.EnableVertexAttribArray(0)

	for i := 0; i < 3*6; i += 3 {
		gl.DrawArrays(gl.LINE_LOOP, 0, int32(3*6))
	}

	gl.DisableVertexAttribArray(0)
}

func (a *AABB) debugInit() {
	verts := []float32{}
	for i := -1.0; i < 2; i += 2 {
		for j := -1.0; j < 2; j += 2 {
			for k := -1.0; k < 2; k += 2 {
				v := []float32{
					a.HalfSize[0] * float32(i),
					a.HalfSize[1] * float32(j),
					a.HalfSize[2] * float32(k),
				}
				verts = append(verts, v...)
			}
		}
	}
	fmt.Println("verts", verts)
	gl.GenVertexArrays(1, &a.vao)
	gl.BindVertexArray(a.vao)

	var vertId uint32
	gl.GenBuffers(1, &vertId)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertId)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (a *AABB) NewBody(mass float32) *Body {
	return &Body{
		b:    a,
		mass: mass,
	}
}
