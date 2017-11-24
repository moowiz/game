package main

import (
	"fmt"
	"math"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moowiz/game/camera"
	"github.com/moowiz/game/physics"
)

var _ = fmt.Print

type player struct {
	camera                         *camera.FPCamera
	body                           *physics.Body
	mat                            *mgl32.Mat4
	moveAmount                     int
	strafeAmount                   int
	horizontalAngle, verticalAngle float64
}

const mouseSpeed = 0.5
const PlayerSpeed = 1

func newPlayer() *player {
	box := physics.AABB{
		Center:   mgl32.Vec3{0, 0, 0},
		HalfSize: mgl32.Vec3{0.01, 0.01, 0.01},
	}

	var p *player

	p = &player{
		body: box.NewBody(float32(math.Inf(1))),
		camera: camera.NewFPCamera(func() mgl32.Vec3 {
			return p.body.Position()
		}),
	}
	return p
}

func (p *player) update(elapsed float64) {
	dV := p.camera.GetDirection().Mul(float32(p.moveAmount))
	dV = dV.Add(p.camera.GetRight().Mul(float32(p.strafeAmount)))

	// No vertical movement for now
	dV[1] = 0
	dV = dV.Normalize().Mul(float32(elapsed))

	if !math.IsNaN(float64(dV.Len())) {
		p.body.SetVelocity(dV)
	} else {
		p.body.SetVelocity(mgl32.Vec3{})
	}
}

func (p *player) updateCamera(program uint32) {
}

var lastX, lastY float64 = -10000, -10000

func (p *player) updateWindow(w *glfw.Window, elapsed float64) {
	p.camera.Update(w, elapsed)
}

func (p *player) move(amount int) {
	p.moveAmount = amount
}

func (p *player) strafe(amount int) {
	p.strafeAmount = amount
}
