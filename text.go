package main

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

var _ = fmt.Print

type Font struct {
	c       *freetype.Context
	bufName uint32
	texture uint32
	program uint32
	vao     uint32
	size    float64
	spacing float64
}

func (f *Font) init() error {
	shaders, err := readShaders("shaders/text.frag", "shaders/text.vert")
	if err != nil {
		return err
	}

	f.program, err = newProgram(shaders...)
	f.size = 24.0
	f.spacing = 1.5
	return err
}

func newFont() (*Font, error) {
	contents, err := ioutil.ReadFile("/Library/Fonts/Arial.ttf")
	if err != nil {
		return nil, err
	}

	freeFont, err := freetype.ParseFont(contents)
	if err != nil {
		return nil, err
	}

	c := freetype.NewContext()
	f := &Font{
		c: c,
	}
	if err := f.init(); err != nil {
		return nil, err
	}
	f.setUpGL()
	dpi := 72.0

	// Initialize the context.
	fg := image.Black
	c.SetDPI(dpi)
	c.SetFont(freeFont)
	c.SetFontSize(f.size)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	return f, err
}

func (f *Font) DrawString(text string) error {
	bg := image.Transparent
	rgba := image.NewRGBA(image.Rect(0, 0, windowWidth, windowHeight))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	f.c.SetClip(rgba.Bounds())
	f.c.SetDst(rgba)

	// Draw the text.
	pt := freetype.Pt(10, 10+int(f.c.PointToFixed(f.size)>>6))
	_, err := f.c.DrawString(text, pt)
	if err != nil {
		return err
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		//gl.Ptr(newArr))
		gl.Ptr(rgba.Pix))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.UseProgram(f.program)
	projection := mgl32.Ortho2D(0, windowWidth, windowHeight, 0)
	projectionUniform := gl.GetUniformLocation(f.program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
	gl.BindVertexArray(f.vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.TRIANGLES, 0, int32(6))

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	return nil
}

func (f *Font) setUpGL() {
	z := float32(1.0)
	x := float32(800.0)
	y := float32(600.0)
	verts := []float32{
		0, 0, z,
		0, y, z,
		x, 0, z,
		x, 0, z,
		0, y, z,
		x, y, z,
	}
	uvs := []float32{
		0, 0,
		0, 1,
		1, 0,
		1, 0,
		0, 1,
		1, 1,
	}

	gl.GenTextures(1, &f.texture)
	gl.GenVertexArrays(1, &f.vao)
	gl.BindVertexArray(f.vao)

	var vertId uint32
	gl.GenBuffers(1, &vertId)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertId)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	var uvId uint32
	gl.GenBuffers(1, &uvId)
	gl.BindBuffer(gl.ARRAY_BUFFER, uvId)
	gl.BufferData(gl.ARRAY_BUFFER, len(uvs)*4, gl.Ptr(uvs), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.DisableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
