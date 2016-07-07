package physics

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var _ = fmt.Println

type World struct {
	bodies    []*Body
	debug     bool
	debugProg uint32
}

func (w *World) AddBody(b *Body) {
	w.bodies = append(w.bodies, b)
}

func (w *World) DebugInit(prog uint32) {
	w.debugProg = prog
	w.debug = true
	for _, body := range w.bodies {
		body.b.debugInit()
	}
}

func (w *World) Update(callback func(uint32)) {
	compared := map[[2]int]bool{}
	for i, body := range w.bodies {
		for j, other := range w.bodies {
			tup := [2]int{i, j}
			_, seen := compared[tup]
			if !seen && i != j && (body.Collides(other) || other.Collides(body)) {
				body.Resolve(other)
				other.Resolve(body)
				compared[tup] = true
				compared[[2]int{j, i}] = true
			}
		}
	}

	if w.debug {
		gl.UseProgram(w.debugProg)
		callback(w.debugProg)
	}

	for _, body := range w.bodies {
		body.Tick()
		if w.debug {
			body.b.debugDraw(w.debugProg)
		}
	}
}
