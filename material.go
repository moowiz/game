package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	_ "golang.org/x/image/bmp"

	"github.com/moowiz/game/shader"
)

func parseMaterialFromFile(filename string, dataDir string) (*material, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	return parseMaterial(scanner, dataDir)

}

func parseMaterial(scanner *bufio.Scanner, dataDir string) (*material, error) {
	mat := &material{}
	var err error
	for scanner.Scan() && !strings.HasPrefix(scanner.Text(), "newmtl") {
		innerScan := bufio.NewScanner(bytes.NewReader(scanner.Bytes()))
		innerScan.Split(bufio.ScanWords)
		if !innerScan.Scan() {
			return nil, fmt.Errorf("bad material line %s", innerScan.Text())
		}

		switch innerScan.Text() {
		case "Ka":
			mat.ambientColor, err = parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
		case "Kd":
			mat.diffuseColor, err = parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
		case "Ks":
			mat.specularColor, err = parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
		case "map_Ka":
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected filename while reading material")
			}
			mat.ambientTexMap, err = newTexture(filepath.Join(dataDir, innerScan.Text()))
			if err != nil {
				return nil, err
			}
		case "map_Kd":
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected filename while reading material")
			}
			mat.diffuseTexMap, err = newTexture(filepath.Join(dataDir, innerScan.Text()))
			if err != nil {
				return nil, err
			}
		case "illum":
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected number while reading illum of material")
			}
			illum, err := strconv.ParseInt(innerScan.Text(), 32, 10)
			if err != nil {
				return nil, fmt.Errorf("while parsing '%s': %s", innerScan.Text(), err)
			}
			mat.illum = int(illum)
		default:
			//fmt.Printf("ignoring %s\n", scanner.Text())
		}

	}

	//TODO: Figure out which shaders to use
	mat.program, err = shader.ProgramFromIllum(mat.illum)
	if err != nil {
		return nil, err
	}

	return mat, nil
}

type material struct {
	ambientColor  []float32
	diffuseColor  []float32
	specularColor []float32
	ambientTexMap uint32
	diffuseTexMap uint32
	illum         int
	program       uint32
	rgba          *image.RGBA
}

func (m *material) init() {
}

func (m *material) draw() {
	if m.diffuseColor != nil {
		diffuseLoc := gl.GetUniformLocation(m.program, gl.Str("diffuseColor\x00"))
		gl.Uniform3f(diffuseLoc, m.diffuseColor[0], m.diffuseColor[1], m.diffuseColor[2])
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, m.diffuseTexMap)

	if m.ambientColor != nil {
		ambientLoc := gl.GetUniformLocation(m.program, gl.Str("ambientColor\x00"))
		gl.Uniform3f(ambientLoc, m.ambientColor[0], m.ambientColor[1], m.ambientColor[2])
	}
	if m.specularColor != nil {
		specularLoc := gl.GetUniformLocation(m.program, gl.Str("specularColor\x00"))
		gl.Uniform3f(specularLoc, m.specularColor[0], m.specularColor[1], m.specularColor[2])
	}
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("texture %s not found on disk: %s", file, err))
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	newP := make([]uint8, len(rgba.Pix))
	maxX := rgba.Rect.Max.X
	maxY := rgba.Rect.Max.Y
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			for offset := 0; offset < 4; offset++ {
				v := rgba.Pix[y*rgba.Stride+x*4+offset]
				newP[(maxY-y-1)*rgba.Stride+x*4+offset] = v
			}
		}
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
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
		gl.Ptr(newP))

	return texture, nil
}
