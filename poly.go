package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type Poly struct {
	verts   []float32
	uvs     []float32
	normals []float32

	vao      uint32
	material *material
	model    mgl32.Mat4
	pos      mgl32.Mat4
}

func (p *Poly) init() {
	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	var vertId uint32
	gl.GenBuffers(1, &vertId)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertId)
	gl.BufferData(gl.ARRAY_BUFFER, len(p.verts)*4, gl.Ptr(p.verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(0)

	if len(p.uvs) > 0 {
		var uvId uint32
		gl.GenBuffers(1, &uvId)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvId)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.uvs)*4, gl.Ptr(p.uvs), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
		gl.DisableVertexAttribArray(1)
	}

	if len(p.normals) > 0 {

		var normalId uint32
		gl.GenBuffers(1, &normalId)
		gl.BindBuffer(gl.ARRAY_BUFFER, normalId)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.normals)*4, gl.Ptr(p.normals), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
		gl.DisableVertexAttribArray(2)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (p *Poly) getProgram() uint32 {
	return p.material.program
}

func (p *Poly) draw(model mgl32.Mat4) {
	fmt.Printf("drawing %v", p.model)

	if p.material != nil {
		p.material.draw()
	}

	modelUniform := gl.GetUniformLocation(p.getProgram(), gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	gl.BindVertexArray(p.vao)
	gl.EnableVertexAttribArray(0)
	if len(p.uvs) > 0 {
		gl.EnableVertexAttribArray(1)
	}
	if len(p.normals) > 0 {
		gl.EnableVertexAttribArray(2)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.verts)/3))

	gl.DisableVertexAttribArray(0)
	if len(p.uvs) > 0 {
		gl.DisableVertexAttribArray(1)
	}
	if len(p.normals) > 0 {
		gl.DisableVertexAttribArray(2)
	}
}
