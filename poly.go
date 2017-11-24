package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

// Poly contains a set of vertices and associated vertex data, along with
// a material.
type Poly struct {
	verts    []float32
	numVerts int
	uvs      []float32
	normals  []float32

	vao      uint32
	material *material
}

func NewPoly(material *material, verts, uvs, normals []float32) (*Poly, error) {
	p := &Poly{
		material: material,
		verts:    verts,
		numVerts: len(verts),
		uvs:      uvs,
		normals:  normals,
	}
	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	var vertID uint32
	gl.GenBuffers(1, &vertID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertID)
	gl.BufferData(gl.ARRAY_BUFFER, p.numVerts*4, gl.Ptr(p.verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(0)

	if len(p.uvs) > 0 {
		var uvID uint32
		gl.GenBuffers(1, &uvID)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvID)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.uvs)*4, gl.Ptr(p.uvs), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
		gl.DisableVertexAttribArray(1)
	}

	if len(p.normals) > 0 {
		var normalID uint32
		gl.GenBuffers(1, &normalID)
		gl.BindBuffer(gl.ARRAY_BUFFER, normalID)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.normals)*4, gl.Ptr(p.normals), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
		gl.DisableVertexAttribArray(2)
	}

	return p, nil
}

func (p *Poly) getProgram() uint32 {
	return p.material.program
}

// draw draws this poly in the given location in world space.
func (p *Poly) draw(model mgl32.Mat4) {
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

	gl.DrawArrays(gl.TRIANGLES, 0, int32(p.numVerts/3))

	gl.DisableVertexAttribArray(0)
	if len(p.uvs) > 0 {
		gl.DisableVertexAttribArray(1)
	}
	if len(p.normals) > 0 {
		gl.DisableVertexAttribArray(2)
	}
}
