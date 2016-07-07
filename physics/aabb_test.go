package physics

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestAABB(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test AABB", t, func() {
		a := &AABB{
			Center:   mgl32.Vec3{0, 0, 0},
			HalfSize: mgl32.Vec3{0.5, 0.5, 0.5},
		}

		Convey("Should hit", func() {
			b := &AABB{
				Center:   mgl32.Vec3{0.1, 0.1, 0.1},
				HalfSize: mgl32.Vec3{0.5, 0.5, 0.5},
			}

			Convey("close box", func() {
				So(a.Collides(b), ShouldBeTrue)
			})

			Convey("corner box", func() {
				b.Center = mgl32.Vec3{0.49, 0.49, 0.49}
				So(a.Collides(b), ShouldBeTrue)
			})
		})

		Convey("Shouldn't hit far box", func() {
			b := &AABB{
				Center:   mgl32.Vec3{2, 2, 2},
				HalfSize: mgl32.Vec3{0.5, 0.5, 0.5},
			}
			So(a.Collides(b), ShouldBeFalse)
			Convey("but does when first box is wide", func() {
				a.HalfSize = mgl32.Vec3{4.5, 4.5, 4.5}
				So(a.Collides(b), ShouldBeTrue)
			})
		})

	})
}
