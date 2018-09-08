package i3

import (
	log "github.com/sirupsen/logrus"
	"github.com/mikkeloscar/flis/layout"
)

// Layout defines an i3 inspired layout handler.
type Layout struct {
	root *layout.Root
}

// New initiliazes a new empty layout.
func New() *Layout {
	return &Layout{
		root: layout.NewRoot(),
	}
}

// Arrange arranges the layout starting from container specified by start.
func (l *Layout) Arrange(start layout.Container) {
	// TODO: implement with more than one view
	switch c := start.(type) {
	case *layout.View:
		pG := c.Parent().Geometry()
		g := c.Geometry()
		g.Origin = pG.Origin
		g.Size = pG.Size
		log.Debugf("Arranging view %s %dx%d (%d,%d)", c.Title(),
			g.Size.W, g.Size.H, g.Origin.X, g.Origin.Y)
	}
}

// ArrangeRoot arranges the layout starting from the root container and all the
// way down through the layout tree.
func (l *Layout) ArrangeRoot() {
	l.Arrange(l.root)
}
