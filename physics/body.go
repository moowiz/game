package physics

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

var _ = math.Pi
var _ = fmt.Print

type body interface {
	Position() mgl32.Vec3
	SetPosition(mgl32.Vec3)
	Collides(body) bool
	debugDraw(program uint32)
	debugInit()
	fmt.Stringer
}

type Body struct {
	b        body
	forces   map[ForceSource]Force
	position mgl32.Vec3
	velocity mgl32.Vec3
	mass     float32
}

func (b *Body) Position() mgl32.Vec3 {
	return b.b.Position()
}

func (b *Body) SetPosition(new mgl32.Vec3) {
	b.b.SetPosition(new)
}

func (b *Body) Velocity() mgl32.Vec3 {
	return b.velocity
}

func (b *Body) SetVelocity(new mgl32.Vec3) {
	b.velocity = new
}

func (b *Body) EnsureForce(f Force) {
	_, ok := b.forces[f.Source]
	if ok {
		return
	}
	b.forces[f.Source] = f
}

func (b *Body) Collides(other *Body) bool {
	return b.b.Collides(other.b)
}

func (b *Body) Resolve(other *Body) {
	v1 := b.Velocity()
	v2 := other.Velocity()
	//m1 := b.mass
	//m2 := other.mass
	//angle := math.Pi / 2

	// Δvx,2' = 2[ vx,1 - vx,2 + a.(vy,1 - vy,2 ) ] / [(1+a2).(1+m2 /m1 )]   ,
	//delta := (2 / (1 + m1/m2)) * (v2.Sub(v1))
	delta := (v2.Add(v1)).Mul(2)
	fmt.Printf("me %v v %v other %v v %v delta %v\n", b.Position(), b.Velocity(), other.Position(), other.Velocity(), delta)
	b.velocity = b.velocity.Sub(delta)
}

func (b *Body) Tick() {
	b.b.SetPosition(b.b.Position().Add(b.Velocity()))
}

func (b *Body) String() string {
	return b.b.String()
}
