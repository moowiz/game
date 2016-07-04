package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var illumToShader = map[int]string{
	0: "basic.frag",
	1: "ambient.frag",
	2: "ambient.frag",
}

func compileShaders(illum int) ([]uint32, error) {
	shaderName, ok := illumToShader[illum]
	fmt.Printf("doing %v for %v\n", shaderName, illum)
	if !ok {
		return nil, fmt.Errorf("no shader found for illumination %v", illum)
	}
	v, err := readShader(filepath.Join(rootDir, "shaders/default.vert"))
	if err != nil {
		return nil, err
	}
	f, err := readShader(filepath.Join(rootDir, "shaders/"+shaderName))
	if err != nil {
		return nil, err
	}

	return []uint32{v, f}, nil
}

func readShader(filename string) (uint32, error) {
	shaderType := typeFromExtension(filepath.Ext(filename))
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	return compileShader(string(contents), shaderType)
}

func typeFromExtension(ext string) uint32 {
	switch ext {
	case ".vert":
		return gl.VERTEX_SHADER
	case ".frag":
		return gl.FRAGMENT_SHADER
	default:
		panic(fmt.Sprintf("invalid shader extension: %s", ext))
	}
}

func newProgram(shaders ...uint32) (uint32, error) {
	program := gl.CreateProgram()
	fmt.Printf("making %v for %v\n", program, shaders)

	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
