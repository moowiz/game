package main

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/golang/freetype"
	"golang.org/x/image/font"

	"github.com/moowiz/game/shader"
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
	rgba    *image.RGBA
}

func (f *Font) init() error {
	var err error
	f.program, err = shader.ProgramForShaders("text.frag", "text.vert")
	if err != nil {
		return err
	}

	f.size = 24.0
	f.spacing = 1.5
	f.setUpGL()
	return err
}

func getFontName() (string, error) {
	if runtime.GOOS == "darwin" {
		return "/Library/Fonts/Arial.ttf", nil
	} else if runtime.GOOS == "windows" {
		return "C:\\Windows\\Fonts\\arial.ttf", nil
	}
	return "", fmt.Errorf("OS %s not supported", runtime.GOOS)
}

func newFont() (*Font, error) {
	fontName, err := getFontName()
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadFile(fontName)
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
	dpi := 72.0

	// Initialize the context.
	fg := image.Black
	f.rgba = image.NewRGBA(image.Rect(0, 0, windowWidth, windowHeight))
	c.SetDPI(dpi)
	c.SetFont(freeFont)
	c.SetFontSize(f.size)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)
	c.SetClip(f.rgba.Bounds())
	c.SetDst(f.rgba)
	bg := image.Transparent
	draw.Draw(f.rgba, f.rgba.Bounds(), bg, image.ZP, draw.Src)

	return f, err
}

func (f *Font) Printf(x, y int, text string, a ...interface{}) error {
	s := fmt.Sprintf(text, a...)
	pt := freetype.Pt(x, y+int(f.c.PointToFixed(f.size)>>6))
	_, err := f.c.DrawString(s, pt)
	if err != nil {
		return err
	}
	return nil
}

func (f *Font) Draw() error {
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
		int32(f.rgba.Rect.Size().X),
		int32(f.rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(f.rgba.Pix))

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

	// Clear for next time
	bg := image.Transparent
	draw.Draw(f.rgba, f.rgba.Bounds(), bg, image.ZP, draw.Src)

	gl.Disable(gl.BLEND)
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
