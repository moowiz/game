package camera

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var _ = fmt.Print

type FPCamera struct {
	lastX, lastY                   float64
	horizontalAngle, verticalAngle float64
	// mat is a matrix used to store camera matrix. It needs to be a member of
	// this struct so that it doesn't get garbage collected.
	mat mgl32.Mat4
	gen PositionSource
}

const mouseSpeed = 0.5

func NewFPCamera(gen PositionSource) *FPCamera {
	return &FPCamera{
		horizontalAngle: 0,
		verticalAngle:   0,
		lastX:           -10000,
		lastY:           -10000,
		gen:             gen,
	}
}

func (c *FPCamera) getMat4() mgl32.Mat4 {
	//fmt.Printf("pos %s dir %s up %s\n", p.position, p.getDirection(), p.getUp())
	pos := c.gen()
	return mgl32.LookAtV(pos, pos.Add(c.GetDirection()), c.getUp())
}

func (c *FPCamera) GetDirection() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Cos(c.verticalAngle) * math.Sin(c.horizontalAngle)),
		float32(math.Sin(c.verticalAngle)),
		float32(math.Cos(c.verticalAngle) * math.Cos(c.horizontalAngle)),
	}
}

func (c *FPCamera) GetRight() mgl32.Vec3 {
	return mgl32.Vec3{
		float32(math.Sin(c.horizontalAngle - math.Pi/2.0)),
		0,
		float32(math.Cos(c.horizontalAngle - math.Pi/2.0)),
	}
}

func (c *FPCamera) getUp() mgl32.Vec3 {
	// Keep cross order to make +y be up
	return c.GetRight().Cross(c.GetDirection())
}

func (c *FPCamera) Setup(program uint32) {
	c.mat = c.getMat4()
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &c.mat[0])
}

func (c *FPCamera) Update(w *glfw.Window, elapsed float64) {
	x, y := w.GetCursorPos()
	if c.lastX == -10000 {
		c.lastX = x
		c.lastY = y
		return
	}
	dx := c.lastX - x
	dy := c.lastY - y

	c.horizontalAngle += mouseSpeed * elapsed * dx
	c.verticalAngle = float64(mgl32.Clamp(float32(
		c.verticalAngle+mouseSpeed*elapsed*dy), -math.Pi/2, math.Pi/2))
	c.lastX = x
	c.lastY = y
}
