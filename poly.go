package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type Poly struct {
	verts    []float32
	vertId   uint32
	uvs      []float32
	uvId     uint32
	normals  []float32
	normalId uint32
	material *material
	model    mgl32.Mat4
	pos      mgl32.Mat4
	angle    float64
}

var bufID uint32 = 0

func (p *Poly) init() {
	gl.GenBuffers(1, &p.vertId)
	gl.GenBuffers(1, &p.uvId)
	gl.GenBuffers(1, &p.normalId)

	p.model = mgl32.Ident4()
	p.pos = mgl32.Translate3D(0, 0, 0)
}

func (p *Poly) getProgram() uint32 {
	return p.material.program
}

func (p *Poly) draw() {
	errar := gl.GetError()
	fmt.Printf("Got errar %s\n", errar)
	if p.material != nil {
		gl.ActiveTexture(gl.TEXTURE0)
		//fmt.Printf("tex %s\n", p.material.diffuseTexMap)
		gl.BindTexture(gl.TEXTURE_2D, p.material.diffuseTexMap)
		textureUniform := gl.GetUniformLocation(p.getProgram(), gl.Str("tex\x00"))
		gl.Uniform1i(textureUniform, 0)

		p.material.draw()
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, p.vertId)
	gl.BufferData(gl.ARRAY_BUFFER, len(p.verts)*4, gl.Ptr(p.verts), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(p.getProgram(), gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	if p.uvs != nil {
		texCoordAttrib := uint32(gl.GetAttribLocation(p.getProgram(), gl.Str("UV\x00")))
		gl.BindBuffer(gl.ARRAY_BUFFER, p.uvId)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.uvs)*4, gl.Ptr(p.uvs), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(texCoordAttrib)
		gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	}

	if len(p.normals) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, p.normalId)
		gl.BufferData(gl.ARRAY_BUFFER, len(p.normals)*4, gl.Ptr(p.normals), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(p.normalId)
		gl.VertexAttribPointer(p.normalId, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	}

	model := mgl32.HomogRotate3D(float32(p.angle), mgl32.Vec3{0, 1, 0})
	model = p.pos.Mul4(model)
	modelUniform := gl.GetUniformLocation(p.getProgram(), gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	errar = gl.GetError()
	fmt.Printf("Got second errar %s\n", errar)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(p.verts)/3))
	errar = gl.GetError()
	fmt.Printf("Got errar end %s\n", errar)
}
