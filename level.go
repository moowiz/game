package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moowiz/game/physics"
)

var _ = fmt.Print

type Level struct {
	objects    []*Object
	phys       *physics.World
	player     *player
	projection mgl32.Mat4
	lightPos   []float32
}

func loadLevel(filename string) (*Level, error) {
	l := &Level{
		phys:   &physics.World{},
		player: newPlayer(),
		projection: mgl32.Perspective(
			mgl32.DegToRad(60.0), float32(windowWidth)/windowHeight, 0.1, 100.0),
		lightPos: []float32{4, 4, 4},
	}
	lf := &levelFile{}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(lf)

	if err != nil {
		return nil, err
	}

	for _, objLoad := range lf.Objects {
		obj, err := readObject("data/objects/" + objLoad.Object + ".json")
		if err != nil {
			return nil, err
		}

		obj.id = objLoad.ID
		if obj.body != nil && objLoad.Location != nil {
			obj.body.SetPosition(floatToVec3(objLoad.Location))
		}
		l.objects = append(l.objects, obj)
		l.phys.AddBody(obj.body)
	}
	l.phys.AddBody(l.player.body)
	shaders, err := readShaders("shaders/wireframe.vert")
	if err != nil {
		return nil, err
	}
	prog, err := newProgram(shaders...)
	if err != nil {
		return nil, err
	}
	l.phys.DebugInit(prog)

	return l, nil
}

func (l *Level) applyBasics(program uint32) {
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &l.projection[0])

	lightUniform := gl.GetUniformLocation(program, gl.Str("light\x00"))
	gl.Uniform3f(lightUniform, l.lightPos[0], l.lightPos[1], l.lightPos[2])
}

func (l *Level) draw(elapsed float64) {
	l.player.update(elapsed)

	l.phys.Update(l.applyBasics)
	for _, obj := range l.objects {
		program := obj.poly.getProgram()
		gl.UseProgram(program)
		l.player.updateCamera(program)
		l.applyBasics(program)

		obj.draw()
	}
}
func floatToVec3(arr []float32) mgl32.Vec3 {
	if len(arr) != 3 {
		panic("bad conversion")
	}

	return mgl32.Vec3{arr[0], arr[1], arr[2]}
}

type levelFile struct {
	Objects []*objLoad
}

type objLoad struct {
	ID       string    `json:"id"`
	Object   string    `json:"object"`
	Location []float32 `json:"location"`
}
