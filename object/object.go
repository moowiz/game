package object

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/moowiz/game/physics"
)

// Object contains a polygon, along with an associated physics body.
type Object struct {
	Id   string
	Poly *Poly
	Body *physics.Body
}

// draw
func (o *Object) Draw() {
	pos := o.Body.Position()
	model := mgl32.Translate3D(pos[0], pos[1], pos[2])

	o.Poly.draw(model)
}

type objectFile struct {
	Model   string          `json:"model"`
	Physics *physicsSection `json:"physics"`
}

type physicsSection struct {
	AABB *physics.AABB `json:"aabb"`
	Mass float32       `json:"mass"`
}

func (p *physicsSection) Process() (*physics.Body, error) {
	if p.AABB != nil {
		//fmt.Printf("aa %s p %s\n", p.AABB, p)
		return p.AABB.NewBody(p.Mass), nil
	}
	return nil, fmt.Errorf("no valid physics")
}

func ReadObject(filename string) (*Object, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	of := &objectFile{}
	err = json.NewDecoder(f).Decode(&of)
	if err != nil {
		panic(err)
	}

	poly, err := readOBJ(of.Model)
	if err != nil {
		return nil, err
	}

	body, err := of.Physics.Process()
	if err != nil {
		return nil, err
	}

	return &Object{
		Poly: poly,
		Body: body,
	}, nil
}
