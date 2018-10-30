// Scalable Vector Graphics (SVG) visualizations of backgammon using svgo
package svg

import (
	"io"

	svgo "github.com/ajstarks/svgo"

	"github.com/chandler37/gobackgammon/svg"
)

func New(w io.Writer) svg.Drawer {
	return &svgoDrawer{svgo.New(w)}
}

type svgoDrawer struct {
	canvas *svgo.SVG
}

func (d *svgoDrawer) Start(w, h int) {
	d.canvas.Start(w, h)
}

func (d *svgoDrawer) End() {
	d.canvas.End()
}

func (d *svgoDrawer) Rect(x int, y int, w int, h int, s ...string) {
	d.canvas.Rect(x, y, w, h, s...)
}

func (d *svgoDrawer) CenterRect(x int, y int, w int, h int, s ...string) {
	d.canvas.CenterRect(x, y, w, h, s...)
}

func (d *svgoDrawer) Circle(x int, y int, r int, s ...string) {
	d.canvas.Circle(x, y, r, s...)
}

func (d *svgoDrawer) Line(x1 int, y1 int, x2 int, y2 int, s ...string) {
	d.canvas.Line(x1, y1, x2, y2, s...)
}

func (d *svgoDrawer) Polyline(x []int, y []int, s ...string) {
	d.canvas.Polyline(x, y, s...)
}

func (d *svgoDrawer) Text(x int, y int, t string, s ...string) {
	d.canvas.Text(x, y, t, s...)
}
