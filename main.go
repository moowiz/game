package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"unsafe"
	//gtime "time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const windowWidth = 800
const windowHeight = 600
const secondsPerFrame = 1 / 60

var rootDir string

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func callbackThing(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {

	fmt.Printf("DEBUG MAN\n")
}
func checkErr() {
	res := gl.GetError()
	if res != 0 {
		fmt.Printf("err %v\n", res)
	}
}

func main() {
	var err error
	rootDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}

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

	level, err := loadLevel("data/levels/basic.json")
	if err != nil {
		panic(err)
	}
	p := level.player

	font, err := newFont()
	if err != nil {
		panic(err)
	}

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyQ {
			w.SetShouldClose(true)
		}
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
	sinceFPS := 0.0
	fps := -1
	for !window.ShouldClose() {
		checkErr()
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		// Time is in seconds
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		p.updateWindow(window, elapsed)

		level.draw(elapsed)

		if sinceFPS > 1 || fps == -1 {
			fps = int(1 / elapsed)
			sinceFPS = 0
		} else {
			sinceFPS += elapsed
		}

		font.Printf(10, 30, "%v %v", p.body.Position()[0], p.body.Position()[2])
		if fps != -1 {
			font.Printf(10, 10, "%v", fps)
		}
		font.Draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()

		// Constant FPS
		//diff := secondsPerFrame - elapsed
		//fmt.Println("sofar", sofar, "diff", diff)
		//gtime.Sleep(gtime.Duration(diff) * gtime.Second)
	}
}

func makeSquare() *Poly {
	square, err := readOBJ("data/plane.obj")
	if err != nil {
		panic(err)
	}
	return square
}
