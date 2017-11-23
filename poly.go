package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type Poly struct {
	// HEY YOU. Yes you. Don't remove these. Removing these makes all polys break. I think
	// it has to do with how pointers are handled in open gl. If we don't keep a pointer to this around, then
	// I think the Go GC cleans it up and we're sad. We should be able to do this more efficiently or
	// something. Idk ¯\_(ツ)_/¯
	verts       []float32
	uvs         []float32
	normals     []float32
	numVerts    int32
	haveUVs     bool
	haveNormals bool

	vao      uint32
	material *material
	model    mgl32.Mat4
	pos      mgl32.Mat4
}

func NewPoly(verts, uvs, normals []float32, material *material) (*Poly, error) {
	p := &Poly{
		verts:       verts,
		uvs:         uvs,
		normals:     normals,
		numVerts:    int32(len(verts) / 3),
		haveUVs:     len(uvs) > 0,
		haveNormals: len(normals) > 0,
		material:    material,
	}

	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	var vertID uint32
	gl.GenBuffers(1, &vertID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertID)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	if p.haveUVs {
		var uvID uint32
		gl.GenBuffers(1, &uvID)
		gl.BindBuffer(gl.ARRAY_BUFFER, uvID)
		gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*4, gl.Ptr(uvs), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	}

	if p.haveNormals {
		var normalID uint32
		gl.GenBuffers(1, &normalID)
		gl.BindBuffer(gl.ARRAY_BUFFER, normalID)
		gl.BufferData(gl.ARRAY_BUFFER, len(normals)*4, gl.Ptr(normals), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	}

	//gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return p, nil
}

func (p *Poly) getProgram() uint32 {
	return p.material.program
}

func (p *Poly) draw(model mgl32.Mat4) {
	//fmt.Printf("drawing %v", p.model)

	if p.material != nil {
		p.material.draw()
	}

	modelUniform := gl.GetUniformLocation(p.getProgram(), gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &p.model[0])

	gl.BindVertexArray(p.vao)
	gl.EnableVertexAttribArray(0)
	if p.haveUVs {
		gl.EnableVertexAttribArray(1)
	}
	if p.haveNormals {
		gl.EnableVertexAttribArray(2)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, p.numVerts)

	gl.DisableVertexAttribArray(0)
	if p.haveUVs {
		gl.DisableVertexAttribArray(1)
	}
	if p.haveNormals {
		gl.DisableVertexAttribArray(2)
	}
}
