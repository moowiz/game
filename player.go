package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type player struct {
	camera                         mgl32.Mat4
	position                       mgl32.Vec3
	moveAmount                     int
	strafeAmount                   int
	horizontalAngle, verticalAngle float64
}

const mouseSpeed = 0.5
const PlayerSpeed = 1

func newPlayer() *player {
	return &player{
		position:        mgl32.Vec3{0, 0, 0},
		horizontalAngle: 0,
		verticalAngle:   math.Pi / 2,
	}
}

func (p *player) getMat4() mgl32.Mat4 {
	//fmt.Printf("pos %s dir %s up %s\n", p.position, p.getDirection(), p.getUp())
	return mgl32.LookAtV(p.position, p.position.Add(p.getDirection()), p.getUp())
}

func (p *player) getDirection() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(p.verticalAngle) * math.Sin(p.horizontalAngle)),
		float32(math.Sin(p.verticalAngle)),
		float32(math.Cos(p.verticalAngle) * math.Cos(p.horizontalAngle)),
	}
}

func (p *player) getRight() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Sin(p.horizontalAngle - math.Pi/2.0)),
		0,
		float32(math.Cos(p.horizontalAngle - math.Pi/2.0)),
	}
}

func (p *player) getUp() mgl32.Vec3 {
	// Keep order to make +y be up
	return p.getRight().Cross(p.getDirection())
}

func (p *player) update(program uint32, elapsed float64) {
	dV := p.getDirection().Mul(float32(p.moveAmount))
	dV = dV.Add(p.getRight().Mul(float32(p.strafeAmount)))
	dV[1] = 0
	dV = dV.Normalize().Mul(float32(elapsed))

	if !math.IsNaN(float64(dV.Len())) {
		p.position = p.position.Add(dV)
	}

	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	p.camera = p.getMat4()
	gl.UniformMatrix4fv(cameraUniform, 1, false, &p.camera[0])
}

var lastX, lastY float64 = -10000, -10000

func (p *player) updateWindow(w *glfw.Window, elapsed float64) {
	x, y := w.GetCursorPos()
	if lastX == -10000 {
		lastX = x
		lastY = y
		return
	}
	dx := lastX - x
	dy := lastY - y

	p.horizontalAngle += mouseSpeed * elapsed * dx
	p.verticalAngle = float64(mgl32.Clamp(float32(
		p.verticalAngle+mouseSpeed*elapsed*dy), -math.Pi/2, math.Pi/2))
	lastX = x
	lastY = y
}

func (p *player) move(amount int) {
	p.moveAmount = amount
}

func (p *player) strafe(amount int) {
	p.strafeAmount = amount
}
