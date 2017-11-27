package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	//gtime "time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/moowiz/game/input"
	"github.com/moowiz/game/level"
	"github.com/moowiz/game/player"
)

const windowWidth = 1024
const windowHeight = 768
const secondsPerFrame = 1 / 60

var rootDir string

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
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
	//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	player := player.NewIsometric(func() mgl32.Vec3 { return mgl32.Vec3{0, 0, 0} })
	lvl, err := level.LoadLevel("data/levels/basic.json", windowWidth, windowHeight, player.Camera)
	if err != nil {
		panic(err)
	}

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
		if key == glfw.KeyK {
			lvl.LightPower += 1
		}
		if key == glfw.KeyL {
			lvl.LightPower -= 1
		}
		ki := &input.KeyInput{
			Key:      key,
			Scancode: scancode,
			Action:   action,
			Mods:     mods,
		}
		if player.HandleInput(w, lvl, ki, nil) {
			return
		}
	})
	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		mi := &input.MouseInput{
			Button: button,
			Action: action,
			Mods:   mod,
		}
		if player.HandleInput(w, lvl, nil, mi) {
			return
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

		player.UpdateWindow(window, elapsed)
		player.Update(elapsed)

		lvl.Draw(elapsed)

		if sinceFPS > 1 || fps == -1 {
			fps = int(1 / elapsed)
			sinceFPS = 0
		} else {
			sinceFPS += elapsed
		}

		font.Printf(10, 50, "%v", lvl.LightPower)
		//font.Printf(10, 30, "%v %v", p.Body.Position()[0], p.Body.Position()[2])
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
