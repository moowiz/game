// Renders a textured spinning cube using GLFW 3.1 and OpenGL 4.1 core forward-compatible profile.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	gtime "time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600
const secondsPerFrame = 1 / 60

var rootDir string

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	var err error
	rootDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Root is %v\n", rootDir)

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	p := newPlayer()
	world := &world{
		player: p,
		projection: mgl32.Perspective(
			mgl32.DegToRad(60.0), float32(windowWidth)/windowHeight, 0.1, 10.0),
		lightPos: []float32{4, 4, 4},
	}

	// Configure the vertex data
	square := makeSquare()

	chicken, err := readOBJ("data/r2d2.obj")
	if err != nil {
		panic(err)
	}
	chicken.init()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyW {
			if action == glfw.Release {
				p.move(0)
			} else if action == glfw.Press {
				p.move(1)
			}
		}
		if key == glfw.KeyS {
			if action == glfw.Release {
				p.move(0)
			} else if action == glfw.Press {
				p.move(-1)
			}
		}
		if key == glfw.KeyD {
			if action == glfw.Release {
				p.strafe(0)
			} else if action == glfw.Press {
				p.strafe(1)
			}
		}
		if key == glfw.KeyA {
			if action == glfw.Release {
				p.strafe(0)
			} else if action == glfw.Press {
				p.strafe(-1)
			}
		}
	})

	previousTime := glfw.GetTime()
	sofar := 0
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		// Time is in seconds
		time := glfw.GetTime()
		elapsed := time - previousTime
		sofar += elapsed
		fmt.Printf("elapsed %v\n", elapsed)
		previousTime = time

		p.updateWindow(window, elapsed)

		world.apply(square.getProgram(), elapsed)
		square.draw()

		world.apply(chicken.getProgram(), elapsed)
		chicken.draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

		// Constant FPS
		diff := secondsPerFrame - elapsed
		fmt.Println("sofar", sofar)
		//gtime.Sleep(gtime.Duration(diff) * gtime.Second)
	}
}

func makeSquare() *Poly {
	square, err := readOBJ("data/plane.obj")
	if err != nil {
		panic(err)
	}
	square.init()
	square.pos = mgl32.Translate3D(0, 5, 0)
	return square
}

/*
func makeCube() *Poly {
	return &Poly{
		verts: []float32{
			//  X, Y, Z, U, V
			// Bottom
			-1.0, -1.0, -1.0, 0.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0,
			-1.0, -1.0, 1.0, 0.0, 1.0,
			1.0, -1.0, -1.0, 1.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 1.0,
			-1.0, -1.0, 1.0, 0.0, 1.0,

			// Top
			-1.0, 1.0, -1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 0.0, 1.0,
			1.0, 1.0, -1.0, 1.0, 0.0,
			1.0, 1.0, -1.0, 1.0, 0.0,
			-1.0, 1.0, 1.0, 0.0, 1.0,
			1.0, 1.0, 1.0, 1.0, 1.0,

			// Front
			-1.0, -1.0, 1.0, 1.0, 0.0,
			1.0, -1.0, 1.0, 0.0, 0.0,
			-1.0, 1.0, 1.0, 1.0, 1.0,
			1.0, -1.0, 1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 0.0, 1.0,
			-1.0, 1.0, 1.0, 1.0, 1.0,

			// Back
			-1.0, -1.0, -1.0, 0.0, 0.0,
			-1.0, 1.0, -1.0, 0.0, 1.0,
			1.0, -1.0, -1.0, 1.0, 0.0,
			1.0, -1.0, -1.0, 1.0, 0.0,
			-1.0, 1.0, -1.0, 0.0, 1.0,
			1.0, 1.0, -1.0, 1.0, 1.0,

			// Left
			-1.0, -1.0, 1.0, 0.0, 1.0,
			-1.0, 1.0, -1.0, 1.0, 0.0,
			-1.0, -1.0, -1.0, 0.0, 0.0,
			-1.0, -1.0, 1.0, 0.0, 1.0,
			-1.0, 1.0, 1.0, 1.0, 1.0,
			-1.0, 1.0, -1.0, 1.0, 0.0,

			// Right
			1.0, -1.0, 1.0, 1.0, 1.0,
			1.0, -1.0, -1.0, 1.0, 0.0,
			1.0, 1.0, -1.0, 0.0, 0.0,
			1.0, -1.0, 1.0, 1.0, 1.0,
			1.0, 1.0, -1.0, 0.0, 0.0,
			1.0, 1.0, 1.0, 0.0, 1.0,
		},
	}
}
*/

type world struct {
	player     *player
	projection mgl32.Mat4
	lightPos   []float32
}

func (w *world) apply(program uint32, elapsed float64) {
	gl.UseProgram(program)

	w.player.update(program, elapsed)
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &w.projection[0])

	lightUniform := gl.GetUniformLocation(program, gl.Str("light\x00"))
	gl.Uniform3f(lightUniform, w.lightPos[0], w.lightPos[1], w.lightPos[2])
}
