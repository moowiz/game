package shader

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var illumToShader = map[int]string{
	0: "basic.frag",
	1: "ambient.frag",
	2: "specular.frag",
}

func ProgramFromIllum(illum int) (uint32, error) {
	shaders, err := shadersForIllum(illum)
	if err != nil {
		return 0, err
	}

	prog, err := newProgram(shaders...)
	if err != nil {
		return 0, err
	}
	return prog, err
}

func ProgramForShaders(filenames ...string) (uint32, error) {
	shaders, err := readShaders(filenames...)
	if err != nil {
		return 0, err
	}

	prog, err := newProgram(shaders...)
	if err != nil {
		return 0, err
	}
	return prog, err
}

func shadersForIllum(illum int) ([]uint32, error) {
	shaderName, ok := illumToShader[illum]
	if !ok {
		return nil, fmt.Errorf("no shader found for illumination %v", illum)
	}
	return readShaders("default.vert", shaderName)
}

func readShaders(filenames ...string) ([]uint32, error) {
	shaders := make([]uint32, len(filenames))
	var err error
	for i, fname := range filenames {
		shaders[i], err = readShader(fname)
		if err != nil {
			return nil, err
		}
	}

	return shaders, nil
}

func getCurrentFilepath() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

var shaderForFilename map[string]uint32 = make(map[string]uint32)

func readShader(filename string) (uint32, error) {
	//fmt.Printf("map %v\n", len(shaderForFilename))
	if shader, ok := shaderForFilename[filename]; ok {
		return shader, nil
	}

	shaderType := typeFromExtension(filepath.Ext(filename))
	currentFile := getCurrentFilepath()
	shaderDir := filepath.Join(filepath.Dir(currentFile), "shaders")
	fullFilepath := filepath.Join(shaderDir, filename)
	contents, err := ioutil.ReadFile(fullFilepath)
	if err != nil {
		return 0, err
	}

	shader, err := compileShader(string(contents), shaderType)
	if err != nil {
		return 0, err
	}
	shaderForFilename[filename] = shader
	return shader, err
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

type shaderKey struct {
	A, B uint32
}

var shadersToProgram map[shaderKey]uint32 = make(map[shaderKey]uint32)

func newProgram(shaders ...uint32) (uint32, error) {
	if len(shaders) == 2 {
		if prog, ok := shadersToProgram[shaderKey{shaders[0], shaders[1]}]; ok {
			return prog, nil
		}
	}
	program := gl.CreateProgram()
	fmt.Printf("making new program %v for shaders %v\n", program, shaders)

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
	if len(shaders) == 2 {
		shadersToProgram[shaderKey{shaders[0], shaders[1]}] = program
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
