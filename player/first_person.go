package player

import (
	"math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/moowiz/game/camera"
	"github.com/moowiz/game/input"
	"github.com/moowiz/game/physics"
)

type FPPlayer struct {
	// FIXME: Should not be public
	Body *physics.Body
	// FIXME: Should not be public
	Camera       *camera.FPCamera
	move, strafe int
}

func NewFPPlayer() *FPPlayer {
	box := physics.AABB{
		Center:   mgl32.Vec3{0, 0, 0},
		HalfSize: mgl32.Vec3{0.01, 0.01, 0.01},
	}

	var p *FPPlayer

	p = &FPPlayer{
		Body: box.NewBody(float32(math.Inf(1))),
		Camera: camera.NewFPCamera(func() mgl32.Vec3 {
			return p.Body.Position()
		}),
	}
	return p
}

func (f *FPPlayer) HandleInput(ki *input.KeyInput) bool {
	swallow := false
	if ki.Key == glfw.KeyW {
		swallow = true
		if ki.Action == glfw.Release {
			f.move = 0
		} else if ki.Action == glfw.Press {
			f.move = 1
		}
	}
	if ki.Key == glfw.KeyS {
		swallow = true
		if ki.Action == glfw.Release {
			f.move = 0
		} else if ki.Action == glfw.Press {
			f.move = -1
		}
	}
	if ki.Key == glfw.KeyD {
		swallow = true
		if ki.Action == glfw.Release {
			f.strafe = 0
		} else if ki.Action == glfw.Press {
			f.strafe = 1
		}
	}
	if ki.Key == glfw.KeyA {
		swallow = true
		if ki.Action == glfw.Release {
			f.strafe = 0
		} else if ki.Action == glfw.Press {
			f.strafe = -1
		}
	}
	return swallow
}

func (p *FPPlayer) Update(elapsed float64) {
	dV := p.Camera.GetDirection().Mul(float32(p.move))
	dV = dV.Add(p.Camera.GetRight().Mul(float32(p.strafe)))

	// No vertical movement for now
	dV[1] = 0
	dV = dV.Normalize().Mul(float32(elapsed))

	if !math.IsNaN(float64(dV.Len())) {
		p.Body.SetVelocity(dV)
	} else {
		p.Body.SetVelocity(mgl32.Vec3{})
	}
}

func (p *FPPlayer) UpdateWindow(w *glfw.Window, elapsed float64) {
	p.Camera.Update(w, elapsed)
}
